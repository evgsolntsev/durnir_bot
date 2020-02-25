package fight

import (
	"context"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/evgsolntsev/durnir_bot/player"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

var (
	playerDAO      player.DAO
	playerManager  player.Manager
	fighterDAO     fighter.DAO
	fighterManager fighter.Manager
	fightDAO       DAO
	notificator    *NotificatorMock
	manager        Manager
	ctx            context.Context
)

func TestMain(m *testing.M) {
	rand.Seed(time.Now().Unix())
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		os.Exit(1)
	}

	ctx = context.Background()
	playerDAO = player.NewDAO(ctx, session)
	playerManager = player.NewManager(ctx, playerDAO)
	fighterDAO = fighter.NewDAO(ctx, session)
	fighterManager = fighter.NewManager(ctx, fighterDAO)
	fightDAO = NewDAO(ctx, session)
	notificator = NewNotificator(ctx)
	manager = NewManager(ctx, playerManager, fighterManager, fightDAO, notificator)

	defer func() {
		playerDAO.RemoveAll(ctx)
		fighterDAO.RemoveAll(ctx)
		fightDAO.RemoveAll(ctx)
	}()

	code := m.Run()

	os.Exit(code)
}

func TestJoinFighters(t *testing.T) {
	hID := idtype.NewHex()
	newF := &fighter.Fighter{
		ID:        idtype.NewFighter(),
		Hex:       hID,
		JoinFight: true,
	}

	var fs []*fighter.Fighter
	for i := 0; i < 10; i++ {
		fs = append(fs, &fighter.Fighter{
			ID:  idtype.NewFighter(),
			Hex: hID,
		})
	}

	for _, f := range append(fs, newF) {
		otherF, err := fighterDAO.Insert(ctx, f)
		require.NoError(t, err)
		f.ID = otherF.ID
	}

	var fss []FighterState

	for _, f := range fs {
		fss = append(fss, NewFighterState(f))
	}

	fight := &Fight{
		Fighters: fss,
		Hex:      hID,
	}

	fight, err := fightDAO.Insert(ctx, fight)
	require.NoError(t, err)

	err = manager.JoinFighters(ctx, fight)
	require.NoError(t, err)

	dbFight, err := fightDAO.FindOne(ctx, fight.ID)
	require.NoError(t, err)
	require.Len(t, dbFight.Fighters, 11)

	places := make(map[idtype.Fighter]int)
	for i, f := range dbFight.Fighters {
		places[f.ID] = i
	}

	for i, _ := range fs {
		if i == 0 {
			continue
		}

		place, ok := places[fs[i].ID]
		require.True(t, ok)
		prev, ok := places[fs[i-1].ID]
		require.True(t, ok)

		require.Greater(t, place, prev)
	}

}

func TestFightStep(t *testing.T) {
	type testcase struct {
		name                string
		fightersStateBefore []FighterState
		fightersStateAfter  []FighterState
		fighters            []*fighter.Fighter
		players             []*player.Player
		message             string
	}

	for _, test := range []testcase{{
		name: "Two skip",
		fightersStateBefore: []FighterState{{
			Health: 100,
			Mana:   100,
		}, {
			Health: 100,
			Mana:   100,
		}},
		fightersStateAfter: []FighterState{{
			Health: 100,
			Mana:   100,
		}, {
			Health: 100,
			Mana:   100,
		}},
		fighters: []*fighter.Fighter{{
			Name:   "Монстр Ундо",
			Deck:   []fighter.Card{fighter.CardSkip},
			Health: 100,
			Mana:   100,
		}, {
			Name:   "Лебедь-отступник",
			Deck:   []fighter.Card{fighter.CardSkip},
			Health: 100,
			Mana:   100,
		}},
		players: []*player.Player{{
			Name: "Ундо",
		}, nil},
		message: "Монстр Ундо использует карту \"Пропуск\" и пропускает ход.\nЛебедь-отступник использует карту \"Пропуск\" и пропускает ход.",
	}, {
		name: "Hit and heal",
		fightersStateBefore: []FighterState{{
			Health: 100,
			Mana:   100,
		}, {
			Health: 100,
			Mana:   100,
		}},
		fightersStateAfter: []FighterState{{
			Health: 100,
			Mana:   100,
		}, {
			Health: 100,
			Mana:   100,
		}},
		fighters: []*fighter.Fighter{{
			Name:     "Монстр Ундо",
			Deck:     []fighter.Card{fighter.CardHit},
			Health:   100,
			Mana:     100,
			Fraction: 0,
		}, {
			Name:     "Лебедь-отступник",
			Deck:     []fighter.Card{fighter.CardHeal},
			Health:   100,
			Mana:     100,
			Fraction: 1,
		}},
		players: []*player.Player{{
			Name: "Ундо",
		}, nil},
		message: "Монстр Ундо использует карту \"Удар\". Здоровье Лебедь-отступник теперь 91.\nЛебедь-отступник использует карту \"Лечение\". Лебедь-отступник восстанавливает здоровье до 100.",
	}, {
		name: "Hit without an answer",
		fightersStateBefore: []FighterState{{
			Health: 10,
			Mana:   100,
		}, {
			Health: 8,
			Mana:   100,
		}},
		fightersStateAfter: []FighterState{{
			Health: 10,
			Mana:   100,
		}, {
			Health: 0,
			Mana:   100,
		}},
		fighters: []*fighter.Fighter{{
			Name:     "Монстр Ундо",
			Deck:     []fighter.Card{fighter.CardHit},
			Health:   100,
			Mana:     100,
			Fraction: 0,
		}, {
			Name:     "Лебедь-отступник",
			Deck:     []fighter.Card{fighter.CardHeal},
			Health:   100,
			Mana:     100,
			Fraction: 1,
		}},
		players: []*player.Player{{
			Name: "Ундо",
		}, nil},
		message: "Монстр Ундо использует карту \"Удар\". Здоровье Лебедь-отступник теперь 0.",
	}} {
		t.Run(test.name, func(t *testing.T) {
			defer func() {
				playerDAO.RemoveAll(ctx)
				fighterDAO.RemoveAll(ctx)
				fightDAO.RemoveAll(ctx)
				notificator.Reset(ctx)
			}()

			hexID := idtype.NewHex()
			for i, f := range test.fighters {
				f.Hex = hexID
				dbFighter, err := fighterDAO.Insert(ctx, f)
				require.NoError(t, err)
				test.fighters[i] = dbFighter
				test.fightersStateBefore[i].Fraction = dbFighter.Fraction
				test.fightersStateAfter[i].Fraction = dbFighter.Fraction
				test.fightersStateBefore[i].MaxHealth = dbFighter.Health
				test.fightersStateAfter[i].MaxHealth = dbFighter.Health
				test.fightersStateBefore[i].MaxMana = dbFighter.Mana
				test.fightersStateAfter[i].MaxMana = dbFighter.Mana
			}

			for i, p := range test.players {
				if p != nil {
					p.FighterID = &test.fighters[i].ID
					dbPlayer, err := playerDAO.Insert(ctx, p)
					require.NoError(t, err)
					test.players[i] = dbPlayer
				}
			}

			for i, _ := range test.fightersStateBefore {
				test.fightersStateBefore[i].ID = test.fighters[i].ID
				test.fightersStateAfter[i].ID = test.fighters[i].ID
			}

			fight := &Fight{
				Fighters:    test.fightersStateBefore,
				Started:     true,
				Hex:         hexID,
				UpdatedTime: time.Now().Add(-2 * time.Minute),
			}

			fight, err := fightDAO.Insert(ctx, fight)
			require.NoError(t, err)

			err = manager.Step(ctx, hexID)
			require.NoError(t, err)

			expectedMessagesMap := make(map[idtype.Player]string)
			for _, p := range test.players {
				if p != nil {
					expectedMessagesMap[p.ID] = test.message
				}
			}

			realMessagesMap := make(map[idtype.Player]string)
			for _, m := range notificator.Messages {
				realMessagesMap[m.PlayerID] = m.Text
			}

			require.Equal(t, expectedMessagesMap, realMessagesMap)

			dbFight, err := fightDAO.FindOne(ctx, fight.ID)
			require.NoError(t, err)
			for i, _ := range dbFight.Fighters {
				require.Equal(t, test.fightersStateAfter[i], dbFight.Fighters[i])
			}
		})
	}
}

type Message struct {
	Text     string
	PlayerID idtype.Player
}

type NotificatorMock struct {
	Messages []Message
}

func (nm *NotificatorMock) Send(ctx context.Context, playerID idtype.Player, message string) error {
	nm.Messages = append(nm.Messages, Message{
		Text:     message,
		PlayerID: playerID,
	})
	return nil
}

func (nm *NotificatorMock) Reset(ctx context.Context) {
	nm.Messages = []Message{}
}

func NewNotificator(ctx context.Context) *NotificatorMock {
	return &NotificatorMock{}
}

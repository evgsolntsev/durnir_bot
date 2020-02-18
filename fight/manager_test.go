package fight

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/evgsolntsev/durnir_bot/player"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestJoinFighters(t *testing.T) {
	rand.Seed(time.Now().Unix())
	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	playerDAO := player.NewDAO(ctx, session)
	playerManager := player.NewManager(ctx, playerDAO)
	fighterDAO := fighter.NewDAO(ctx, session)
	fighterManager := fighter.NewManager(ctx, fighterDAO)
	fightDAO := NewDAO(ctx, session)
	notificator := NewNotificator(ctx)
	manager := NewManager(ctx, playerManager, fighterManager, fightDAO, notificator)

	defer func() {
		playerDAO.RemoveAll(ctx)
		fighterDAO.RemoveAll(ctx)
		fightDAO.RemoveAll(ctx)
	}()

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

	fight, err = fightDAO.Insert(ctx, fight)
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

type NotificatorMock struct {
}

func (*NotificatorMock) Send(ctx context.Context, playerID idtype.Player, message string) error {
	return nil
}

func NewNotificator(ctx context.Context) Notificator {
	return &NotificatorMock{}
}

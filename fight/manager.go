package fight

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/evgsolntsev/durnir_bot/player"
	"github.com/globalsign/mgo"
)

type Notificator interface {
	Send(context.Context, idtype.Player, string) error
}

var (
	TimeToStart  = time.Minute * 10
	TimeToUpdate = time.Second * -30
)

type Manager interface {
	Step(context.Context, idtype.Hex) error
	JoinFighters(context.Context, *Fight) error
}

type defaultManager struct {
	PlayerManager  player.Manager
	FighterManager fighter.Manager
	FightDAO       DAO
	Notificator    Notificator
}

func (m *defaultManager) Step(ctx context.Context, hexID idtype.Hex) error {
	fight, err := m.FightDAO.FindOneByHex(ctx, hexID)
	if err != nil {
		if err != mgo.ErrNotFound {
			return err
		} else {
			fighters, err := m.FighterManager.FindJoining(ctx, hexID)
			if err != nil {
				return err
			}
			if len(fighters) > 0 {
				fight, err = m.InitFight(ctx, hexID, fighters)
				if err != nil {
					return err
				}
			}
		}
	}

	now := time.Now()
	if !checkPeriod(fight.UpdatedTime, now, TimeToUpdate) {
		return nil
	}

	stopped, err := m.StopFightIfNeededAndLoot(ctx, fight)
	if err != nil {
		return err
	}
	if stopped {
		return nil
	}

	err = m.JoinFighters(ctx, fight)
	if err != nil {
		return err
	}

	if !fight.Started {
		return m.StartFightIfNeeded(ctx, fight)
	}

	var turns []string
	for i, fs := range fight.Fighters {
		if fight.Fighters[i].Health == 0 {
			continue
		}

		f, err := m.FighterManager.GetOne(ctx, fs.ID)
		if err != nil {
			return err
		}

		var message string
		card := f.GetCard(ctx)
		switch card {
		case fighter.CardHeal:
			target, err := m.GetRandomFromSameFraction(ctx, fight, fs)
			if err != nil {
				message = fmt.Sprintf(
					"%s использует карту \"%s\" и получает ошибку: %s.",
					f.Name, card.Name(), err.Error())
			} else {
				targetFighter, err := m.FighterManager.GetOne(ctx, target.ID)
				if err != nil {
					return err
				}
				healedFull := target.Health + 10
				if healedFull < target.MaxHealth {
					target.Health = healedFull
				} else {
					target.Health = target.MaxHealth
				}
				message = fmt.Sprintf(
					"%s использует карту \"%s\". %s восстанавливает здоровье до %d.",
					f.Name, card.Name(), targetFighter.Name, target.Health)
			}
		case fighter.CardHit:
			target, err := m.GetRandomFromAnotherFraction(ctx, fight, fs)
			if err != nil {
				message = fmt.Sprintf(
					"%s использует карту \"%s\" и получает ошибку: %s.",
					f.Name, card.Name(), err.Error())
			} else {
				targetFighter, err := m.FighterManager.GetOne(ctx, target.ID)
				if err != nil {
					return err
				}
				heatedFull := target.Health - 9
				if heatedFull < 0 {
					target.Health = 0
				} else {
					target.Health = heatedFull
				}
				message = fmt.Sprintf(
					"%s использует карту \"%s\". Здоровье %s теперь %d.",
					f.Name, card.Name(), targetFighter.Name, target.Health)
			}
		case fighter.CardSkip:
			message = fmt.Sprintf(
				"%s использует карту \"%s\" и пропускает ход.", f.Name, card.Name())
		default:
			return fmt.Errorf("Unknown card type!")
		}

		turns = append(turns, message)
	}

	fight.UpdatedTime = time.Now()
	err = m.FightDAO.Update(ctx, fight)
	if err != nil {
		return err
	}

	return m.NotificateFighters(ctx, fight, strings.Join(turns, "\n"))
}

func (m *defaultManager) GetRandomFromSameFraction(ctx context.Context, fight *Fight, state FighterState) (*FighterState, error) {
	var states []*FighterState
	for i, _ := range fight.Fighters {
		if fight.Fighters[i].Fraction == state.Fraction {
			states = append(states, &fight.Fighters[i])
		}
	}
	return getRandom(ctx, states)
}

func (m *defaultManager) GetRandomFromAnotherFraction(ctx context.Context, fight *Fight, state FighterState) (*FighterState, error) {
	var states []*FighterState
	for i, _ := range fight.Fighters {
		if fight.Fighters[i].Fraction != state.Fraction {
			states = append(states, &fight.Fighters[i])
		}
	}
	return getRandom(ctx, states)
}

func (m *defaultManager) StartFightIfNeeded(ctx context.Context, fight *Fight) error {
	now := time.Now()
	if !checkPeriod(fight.UpdatedTime, now, TimeToStart) {
		return nil
	}

	fight.UpdatedTime = now
	fight.Started = true
	err := m.FightDAO.Update(ctx, fight)
	if err != nil {
		return err
	}

	message := fmt.Sprintf("Битва на гексе %v началась!", fight.Hex)
	return m.NotificateFighters(ctx, fight, message)
}

func (m *defaultManager) StopFightIfNeededAndLoot(ctx context.Context, fight *Fight) (bool, error) {
	return false, nil
}

func (m *defaultManager) NotificateFighters(ctx context.Context, fight *Fight, message string) error {
	var fighterIDs []idtype.Fighter
	for _, fighterState := range fight.Fighters {
		fighterIDs = append(fighterIDs, fighterState.ID)
	}

	players, err := m.PlayerManager.FindPlayersByFighters(ctx, fighterIDs)
	if err != nil {
		return err
	}

	var errs []error
	for _, p := range players {
		err := m.Notificator.Send(ctx, p.ID, message)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs[0]
	}
	return nil
}

func (m *defaultManager) JoinFighters(ctx context.Context, fight *Fight) error {
	fighters, err := m.FighterManager.FindJoining(ctx, fight.Hex)
	if err != nil {
		return err
	}

	for _, f := range fighters {
		state := NewFighterState(f)
		i := rand.Intn(len(fight.Fighters) + 1)
		if i == 0 {
			fight.Fighters = append([]FighterState{state}, fight.Fighters...)
		} else if i == len(fight.Fighters) {
			fight.Fighters = append(fight.Fighters, state)
		} else {
			left := make([]FighterState, i)
			copy(left, fight.Fighters[0:i])
			right := make([]FighterState, len(fight.Fighters)-i)
			copy(right, fight.Fighters[i:])
			fight.Fighters = append(left, state)
			fight.Fighters = append(fight.Fighters, right...)
		}
	}
	if err := m.FightDAO.Update(ctx, fight); err != nil {
		return err
	}

	for _, f := range fighters {
		if err := m.FighterManager.SetJoining(ctx, f.ID, false); err != nil {
			return err
		}
	}

	return nil
}

func (m *defaultManager) InitFight(ctx context.Context, hexID idtype.Hex, fighters []*fighter.Fighter) (*Fight, error) {
	var fighterStates []FighterState
	for _, f := range fighters {
		fighterStates = append(fighterStates, NewFighterState(f))
	}

	now := time.Now()
	fight := &Fight{
		Fighters:    fighterStates,
		UpdatedTime: now,
		Started:     false,
		Hex:         hexID,
	}
	return m.FightDAO.Insert(ctx, fight)
}

func checkPeriod(a, b time.Time, d time.Duration) bool {
	return b.After(a.Add(d))
}

func NewManager(
	ctx context.Context, playerManager player.Manager,
	fighterManager fighter.Manager, dao DAO, notificator Notificator,
) Manager {
	return &defaultManager{
		PlayerManager:  playerManager,
		FighterManager: fighterManager,
		FightDAO:       dao,
		Notificator:    notificator,
	}
}

func getRandom(ctx context.Context, states []*FighterState) (*FighterState, error) {
	if len(states) == 0 {
		return nil, fmt.Errorf("цель не найдена")
	}
	return states[rand.Intn(len(states))], nil
}

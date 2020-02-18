package fight

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/evgsolntsev/durnir_bot/player"
)

type Notificator interface {
	Send(context.Context, idtype.Player, string) error
}

var (
	TimeToStart  = time.Minute * 10
	TimeToUpdate = time.Second * -30
)

type Manager interface {
	Step(context.Context, *Fight) error
	JoinFighters(context.Context, *Fight) error
}

type defaultManager struct {
	PlayerManager  player.Manager
	FighterManager fighter.Manager
	FightDAO       DAO
	Notificator    Notificator
}

func (m *defaultManager) Step(ctx context.Context, fight *Fight) error {
	err := m.JoinFighters(ctx, fight)
	if err != nil {
		return err
	}

	if !fight.Started {
		return m.StartFightIfNeeded(ctx, fight)
	}

	now := time.Now()
	if !checkPeriod(fight.UpdatedTime, now, TimeToUpdate) {
		return nil
	}

	// TODO: actual step
	return nil
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
	return m.FightDAO.Update(ctx, fight)
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

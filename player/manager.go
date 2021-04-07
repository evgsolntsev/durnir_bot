package player

import (
	"context"
	"fmt"
	"log"

	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Manager interface {
	Update(context.Context, *Player) error
	FindPlayersByFighters(context.Context, []idtype.Fighter) ([]Player, error)
	FindPlayersByFightersMap(context.Context, []idtype.Fighter) (map[idtype.Fighter]Player, error)
	GetOne(context.Context, idtype.Player) (*Player, error)
	GetOneByTelegramId(context.Context, int64) (*Player, error)
	GenerateFighter(context.Context, *Player) error
}

type defaultManager struct {
	PlayerDAO      DAO
	FighterManager fighter.Manager
}

var _ Manager = (*defaultManager)(nil)

func NewManager(
	ctx context.Context,
	dao DAO,
	fighterManager fighter.Manager,
) *defaultManager {
	return &defaultManager{
		PlayerDAO:      dao,
		FighterManager: fighterManager,
	}
}

func (m *defaultManager) FindPlayersByFighters(ctx context.Context, fighterIDs []idtype.Fighter) ([]Player, error) {
	return m.PlayerDAO.FindByFighters(ctx, fighterIDs)
}

func (m *defaultManager) FindPlayersByFightersMap(ctx context.Context, fighterIDs []idtype.Fighter) (map[idtype.Fighter]Player, error) {
	players, err := m.FindPlayersByFighters(ctx, fighterIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[idtype.Fighter]Player)
	for _, p := range players {
		result[*p.FighterID] = p
	}
	return result, nil
}

func (m *defaultManager) Update(ctx context.Context, p *Player) error {
	return m.PlayerDAO.Update(ctx, p)
}

func (d *defaultManager) GetOne(ctx context.Context, pID idtype.Player) (*Player, error) {
	return d.PlayerDAO.FindOne(ctx, pID)
}

func (d *defaultManager) GetOneByTelegramId(ctx context.Context, telegramId int64) (*Player, error) {
	return d.PlayerDAO.FindOneByTelegramId(ctx, telegramId)
}

func (d *defaultManager) GenerateFighter(ctx context.Context, p *Player) error {
	newFighter, err := d.FighterManager.Create(
		ctx, fmt.Sprintf("Монстр %s", p.Name), fighter.FractionPlayers)
	if err != nil {
		return err
	}

	err = d.PlayerDAO.SetFighterID(ctx, p.ID, newFighter.ID)
	if err != nil {
		removingErr := d.FighterManager.RemoveOne(ctx, newFighter.ID)
		if removingErr != nil {
			log.Printf(
				"Failed to remove fighter on creating rollback:\nFirst err: %s\nRollback err: %s",
				err.Error(), removingErr.Error())
		}
		return fmt.Errorf("Something went wrong: %s", err.Error())
	}

	return nil
}

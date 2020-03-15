package player

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Manager interface {
	Update(context.Context, *Player) error
	FindPlayersByFighters(context.Context, []idtype.Fighter) ([]Player, error)
	FindPlayersByFightersMap(context.Context, []idtype.Fighter) (map[idtype.Fighter]Player, error)
	GetOne(context.Context, idtype.Player) (*Player, error)
	GetOneByTelegramId(context.Context, int64) (*Player, error)
}

type defaultManager struct {
	PlayerDAO DAO
}

var _ Manager = (*defaultManager)(nil)

func NewManager(ctx context.Context, dao DAO) *defaultManager {
	return &defaultManager{
		PlayerDAO: dao,
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

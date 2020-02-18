package player

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Manager interface {
	FindPlayersByFighters(context.Context, []idtype.Fighter) ([]Player, error)
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

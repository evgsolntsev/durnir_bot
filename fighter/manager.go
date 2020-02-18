package fighter

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Manager interface {
	FindJoining(context.Context, idtype.Hex) ([]*Fighter, error)
}

type defaultManager struct {
	FighterDAO DAO
}

var _ Manager = (*defaultManager)(nil)

func NewManager(ctx context.Context, dao DAO) *defaultManager {
	return &defaultManager{
		FighterDAO: dao,
	}
}

func (d *defaultManager) FindJoining(ctx context.Context, hexID idtype.Hex) ([]*Fighter, error) {
	return d.FighterDAO.FindJoining(ctx, hexID)
}

package fighter

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Manager interface {
	SetJoining(context.Context, idtype.Fighter, bool) error
	FindJoining(context.Context, idtype.Hex) ([]*Fighter, error)
	GetOne(context.Context, idtype.Fighter) (*Fighter, error)
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

func (d *defaultManager) SetJoining(ctx context.Context, fID idtype.Fighter, join bool) error {
	fighter, err := d.FighterDAO.FindOne(ctx, fID)
	if err != nil {
		return err
	}

	fighter.JoinFight = join
	return d.FighterDAO.Update(ctx, fighter)
}

func (d *defaultManager) GetOne(ctx context.Context, fID idtype.Fighter) (*Fighter, error) {
	return d.FighterDAO.FindOne(ctx, fID)
}

package fighter

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Manager interface {
	SetJoining(context.Context, idtype.Fighter, bool) error
	FindJoining(context.Context, idtype.Hex) ([]*Fighter, error)
	GetOne(context.Context, idtype.Fighter) (*Fighter, error)
	GetMapByIDs(context.Context, []idtype.Fighter) (map[idtype.Fighter]Fighter, error)
	Update(context.Context, *Fighter) error
	RemoveOne(context.Context, idtype.Fighter) error
	Create(context.Context, string) (*Fighter, error)
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

func (d *defaultManager) GetMapByIDs(ctx context.Context, fIDs []idtype.Fighter) (map[idtype.Fighter]Fighter, error) {
	fighters, err := d.FighterDAO.FindByIDs(ctx, fIDs)
	if err != nil {
		return nil, err
	}
	result := make(map[idtype.Fighter]Fighter)
	for _, f := range fighters {
		result[f.ID] = *f
	}
	return result, nil

}

func (d *defaultManager) RemoveOne(ctx context.Context, fID idtype.Fighter) error {
	return d.FighterDAO.RemoveOne(ctx, fID)
}

func (m *defaultManager) Update(ctx context.Context, f *Fighter) error {
	return m.FighterDAO.Update(ctx, f)
}

func (m *defaultManager) Create(ctx context.Context, name string) (*Fighter, error) {
	result := &Fighter{
		Name:      name,
		Health:    100,
		Mana:      100,
		Will:      1,
		Power:     1,
		FearPower: 1,
		Hex:       idtype.StartHex,
	}

	real, err := m.FighterDAO.Insert(ctx, result)
	if err != nil {
		return nil, err
	}

	return real, nil
}

package fighter

import "context"

type Manager interface {
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

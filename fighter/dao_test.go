package fighter

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestFighterDAO(t *testing.T) {
	f := &Fighter{
		Health:    100,
		Mana:      100,
		Shield:    2,
		Will:      3,
		Power:     4,
		FearPower: 3,
		Hex:       5,
		Deck:      []Card{},
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDefaultDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	f.Mana = 300
	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)	

	dbF, err := dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)

	f.Health = 200
	f.Will = 2
	err = dao.Update(ctx, f)
	require.Nil(t, err)

	dbF, err = dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)
}

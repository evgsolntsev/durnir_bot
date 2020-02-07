package fight

import (
	"context"
	"testing"

	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestFighterDAO(t *testing.T) {
	f := &Fight{
		FighterIDs: []idtype.Fighter{idtype.NewFighter(), idtype.NewFighter()},
		Started:    false,
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDefaultDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	f.FighterIDs = append(f.FighterIDs, idtype.NewFighter())
	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	dbF, err := dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)

	f.Started = true
	err = dao.Update(ctx, f)
	require.Nil(t, err)

	dbF, err = dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)
}

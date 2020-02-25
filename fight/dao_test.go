package fight

import (
	"context"
	"testing"

	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestFightDAO(t *testing.T) {
	h1 := idtype.NewHex()
	h2 := idtype.NewHex()
	fs1 := FighterState{
		ID: idtype.NewFighter(),
	}
	fs2 := FighterState{
		ID: idtype.NewFighter(),
	}
	fs3 := FighterState{
		ID: idtype.NewFighter(),
	}
	f := &Fight{
		Fighters: []FighterState{fs1, fs2},
		Started:  false,
		Hex:      h2,
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	f.Fighters = append(f.Fighters, fs3)
	f.Hex = h1
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

	dbF, err = dao.FindOneByHex(ctx, h1)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)

	dbF, err = dao.FindOneByHex(ctx, idtype.NewHex())
	require.Nil(t, dbF)
	require.Error(t, err)
}

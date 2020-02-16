package fight

import (
	"context"
	"testing"

	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestFighterDAO(t *testing.T) {
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
		Fighters:  []FighterState{fs1, fs2},
		Started:    false,
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	f.Fighters = append(f.Fighters, fs3)
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

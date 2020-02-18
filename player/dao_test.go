package player

import (
	"context"
	"fmt"
	"testing"

	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestPlayerDAO(t *testing.T) {
	f := &Player{
		Name: "Сарасти",
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	f.Name = "Ундо"
	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	dbF, err := dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)

	dbF, err = dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)
}

func TestFindByFighters(t *testing.T) {
	f1 := idtype.NewFighter()
	f2 := idtype.NewFighter()
	f3 := idtype.NewFighter()
	p1 := &Player{
		Name:      "Areatangent",
		FighterID: &f1,
	}
	p2 := &Player{
		Name:      "UnheiligZ",
		FighterID: &f2,
	}
	p3 := &Player{
		Name:      "Нечто",
		FighterID: &f3,
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	p1, err = dao.Insert(ctx, p1)
	require.Nil(t, err)
	p2, err = dao.Insert(ctx, p2)
	require.Nil(t, err)
	p3, err = dao.Insert(ctx, p3)
	require.Nil(t, err)

	result, err := dao.FindByFighters(ctx, []idtype.Fighter{f1, f2})
	require.Nil(t, err)

	var realIDs []idtype.Player
	for _, r := range result {
		realIDs = append(realIDs, r.ID)
	}

	expectedIDs := []idtype.Player{p1.ID, p2.ID}
	require.Equal(t, expectedIDs, realIDs)
}

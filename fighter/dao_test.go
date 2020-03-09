package fighter

import (
	"context"
	"testing"

	"github.com/evgsolntsev/durnir_bot/idtype"
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
		Parts:     []Part{},
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	f.ID = idtype.NewFighter()
	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	f.Mana = 300
	f.ID = idtype.NewFighter()
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

func TestFindJoining(t *testing.T) {
	h1 := idtype.NewHex()
	h2 := idtype.NewHex()
	f1 := &Fighter{
		JoinFight: true,
		Hex:       h1,
	}
	f2 := &Fighter{
		JoinFight: false,
		Hex:       h1,
	}
	f3 := &Fighter{
		JoinFight: true,
		Hex:       h1,
	}
	f4 := &Fighter{
		JoinFight: true,
		Hex:       h2,
	}

	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	f1, err = dao.Insert(ctx, f1)
	require.Nil(t, err)
	f2, err = dao.Insert(ctx, f2)
	require.Nil(t, err)
	f3, err = dao.Insert(ctx, f3)
	require.Nil(t, err)
	f4, err = dao.Insert(ctx, f4)
	require.Nil(t, err)
	expectedIDs := []idtype.Fighter{f1.ID, f3.ID}

	fs, err := dao.FindJoining(ctx, h1)
	require.Nil(t, err)

	var realIDs []idtype.Fighter
	for _, f := range fs {
		realIDs = append(realIDs, f.ID)
	}

	require.Equal(t, expectedIDs, realIDs)
}

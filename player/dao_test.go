package player

import (
	"context"
	"os"
	"testing"

	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

var (
	dao     DAO
	manager Manager
)

func TestMain(m *testing.M) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	if err != nil {
		os.Exit(1)
	}

	ctx := context.Background()
	dao = NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)

	code := m.Run()

	os.Exit(code)
}

func TestPlayerDAO(t *testing.T) {
	ctx := context.Background()
	dao.RemoveAll(ctx)

	f := &Player{
		Name:       "Сарасти",
		TelegramId: 123,
		Parts:      []fighter.Part{},
	}

	f, err := dao.Insert(ctx, f)
	require.Nil(t, err)

	f.Name = "Ундо"
	f.TelegramId = 124
	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	dbF, err := dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)

	dbF, err = dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)

	dbF, err = dao.FindOneByTelegramId(ctx, 124)
	require.Nil(t, err)
	require.Equal(t, *f, *dbF)
}

func TestFindByFighters(t *testing.T) {
	ctx := context.Background()

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

	p1, err := dao.Insert(ctx, p1)
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

func TestSetFighterID(t *testing.T) {
	ctx := context.Background()

	fID := idtype.NewFighter()
	p, err := dao.Insert(ctx, &Player{})
	require.Nil(t, err)

	err = dao.SetFighterID(ctx, p.ID, fID)
	require.Nil(t, err)

	dbP, err := dao.FindOne(ctx, p.ID)
	require.Nil(t, err)
	require.Equal(t, fID, *dbP.FighterID)

	f2ID := idtype.NewFighter()
	err = dao.SetFighterID(ctx, p.ID, f2ID)
	require.NotNil(t, err)

	dbP, err = dao.FindOne(ctx, p.ID)
	require.Nil(t, err)
	require.Equal(t, fID, *dbP.FighterID)
}

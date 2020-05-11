package fighter

import (
	"context"
	"os"
	"testing"

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
	manager = NewManager(ctx, dao)

	code := m.Run()

	os.Exit(code)
}

func TestSetJoining(t *testing.T) {
	ctx := context.Background()

	var err error
	f := &Fighter{
		JoinFight: false,
	}
	f, err = dao.Insert(ctx, f)
	require.Nil(t, err)

	dbF, err := dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.Equal(t, false, dbF.JoinFight)

	for _, value := range []bool{true, false} {
		err = manager.SetJoining(ctx, f.ID, value)
		require.Nil(t, err)

		dbF, err = dao.FindOne(ctx, f.ID)
		require.Nil(t, err)
		require.Equal(t, value, dbF.JoinFight)
	}
}

func TestGetMapByIDs(t *testing.T) {
	ctx := context.Background()

	f1 := &Fighter{
		Health: 15,
	}
	f1, err := dao.Insert(ctx, f1)
	require.Nil(t, err)

	f2 := &Fighter{
		Health: 25,
	}
	f2, err = dao.Insert(ctx, f2)
	require.Nil(t, err)

	m, err := manager.GetMapByIDs(ctx, []idtype.Fighter{f1.ID, f2.ID})
	require.Nil(t, err)

	result, ok := m[f1.ID]
	require.True(t, ok)
	require.Equal(t, 15, result.Health)

	result, ok = m[f2.ID]
	require.True(t, ok)
	require.Equal(t, 25, result.Health)
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	f, err := manager.Create(ctx, "KEK")
	require.Nil(t, err)

	dbF, err := dao.FindOne(ctx, f.ID)
	require.Nil(t, err)
	require.NotNil(t, dbF)
	require.Equal(t, "KEK", dbF.Name)
}

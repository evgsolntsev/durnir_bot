package fighter

import (
	"context"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/require"
)

func TestSetJoining(t *testing.T) {
	session, err := mgo.Dial("mongodb://localhost:27017")
	require.Nil(t, err)

	ctx := context.Background()
	dao := NewDAO(ctx, session)
	defer dao.RemoveAll(ctx)
	manager := NewManager(ctx, dao)

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

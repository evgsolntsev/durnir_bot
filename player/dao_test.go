package player

import (
	"context"
	"testing"

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

package player

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	DatabaseName   = "dbname"
	CollectionName = "players"
)

type DAO interface {
	FindOne(context.Context, idtype.Player) (*Player, error)
	Update(context.Context, *Player) error
	Insert(context.Context, *Player) (*Player, error)
	RemoveAll(context.Context) error
}

var _ DAO = (*defaultDAO)(nil)

type defaultDAO struct {
	collection *mgo.Collection
}

func NewDAO(ctx context.Context, session *mgo.Session) *defaultDAO {
	return &defaultDAO{
		collection: session.DB(DatabaseName).C(CollectionName),
	}
}

func (d *defaultDAO) FindOne(ctx context.Context, id idtype.Player) (*Player, error) {
	var result Player
	err := d.collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *defaultDAO) Update(ctx context.Context, player *Player) error {
	return d.collection.UpdateId(player.ID, player)
}

func (d *defaultDAO) Insert(ctx context.Context, player *Player) (*Player, error) {
	player.ID = idtype.NewPlayer()
	err := d.collection.Insert(player)
	if err != nil {
		return nil, err
	}
	return player, nil
}

func (d *defaultDAO) RemoveAll(ctx context.Context) error {
	_, err := d.collection.RemoveAll(bson.M{})
	return err
}

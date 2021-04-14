package fight

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	DatabaseName   = "dbname"
	CollectionName = "fights"
)

type DAO interface {
	FindOne(context.Context, idtype.Fight) (*Fight, error)
	FindOneByHex(context.Context, idtype.Hex) (*Fight, error)
	Update(context.Context, *Fight) error
	Insert(context.Context, *Fight) (*Fight, error)
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

func (d *defaultDAO) FindOne(ctx context.Context, id idtype.Fight) (*Fight, error) {
	var result Fight
	err := d.collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *defaultDAO) FindOneByHex(ctx context.Context, id idtype.Hex) (*Fight, error) {
	var result Fight
	err := d.collection.Find(bson.M{"hex": id}).One(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *defaultDAO) Update(ctx context.Context, fight *Fight) error {
	return d.collection.UpdateId(fight.ID, fight)
}

func (d *defaultDAO) Insert(ctx context.Context, fight *Fight) (*Fight, error) {
	fight.ID = idtype.NewFight()
	err := d.collection.Insert(fight)
	if err != nil {
		return nil, err
	}
	return fight, nil
}

func (d *defaultDAO) RemoveAll(ctx context.Context) error {
	_, err := d.collection.RemoveAll(bson.M{})
	return err
}

func (d *defaultDAO) RemoveOne(ctx context.Context, fight *Fight) error {
	_, err := d.collection.Remove(bson.M{"_id": fight.ID})
	return err
}

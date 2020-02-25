package fighter

import (
	"context"

	"github.com/evgsolntsev/durnir_bot/idtype"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

var (
	DatabaseName   = "dbname"
	CollectionName = "fighters"
)

type DAO interface {
	FindOne(context.Context, idtype.Fighter) (*Fighter, error)
	FindJoining(context.Context, idtype.Hex) ([]*Fighter, error)
	Update(context.Context, *Fighter) error
	Insert(context.Context, *Fighter) (*Fighter, error)
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

func (d *defaultDAO) FindOne(ctx context.Context, id idtype.Fighter) (*Fighter, error) {
	var result Fighter
	err := d.collection.Find(bson.M{"_id": id}).One(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (d *defaultDAO) FindJoining(ctx context.Context, hexID idtype.Hex) ([]*Fighter, error) {
	query := bson.M{
		"joinFight": true,
		"hex":       hexID,
	}
	var result []*Fighter
	err := d.collection.Find(query).All(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *defaultDAO) Update(ctx context.Context, fighter *Fighter) error {
	return d.collection.UpdateId(fighter.ID, fighter)
}

func (d *defaultDAO) Insert(ctx context.Context, fighter *Fighter) (*Fighter, error) {
	fighter.ID = idtype.NewFighter()
	err := d.collection.Insert(fighter)
	if err != nil {
		return nil, err
	}
	return fighter, nil
}

func (d *defaultDAO) RemoveAll(ctx context.Context) error {
	_, err := d.collection.RemoveAll(bson.M{})
	return err
}

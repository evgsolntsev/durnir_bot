package idtype

import (
	"math/rand"

	"github.com/globalsign/mgo/bson"
)

type Fighter bson.ObjectId

func (f Fighter) GetBSON() (interface{}, error) {
	return bson.ObjectId(f), nil
}

func NewFighter() Fighter {
	return Fighter(bson.NewObjectId())
}

type Player bson.ObjectId

func (p Player) GetBSON() (interface{}, error) {
	return bson.ObjectId(p), nil
}

func NewPlayer() Player {
	return Player(bson.NewObjectId())
}

type Fight bson.ObjectId

func (f Fight) GetBSON() (interface{}, error) {
	return bson.ObjectId(f), nil
}

func NewFight() Fight {
	return Fight(bson.NewObjectId())
}

var (
	_ bson.Getter = (*Player)(nil)
	_ bson.Getter = (*Fighter)(nil)
	_ bson.Getter = (*Fight)(nil)
)

type Hex int

func NewHex() Hex {
	return Hex(rand.Intn(1000))
}

var (
	ZeroPlayer  = Player(bson.ObjectId(""))
	ZeroFighter = Fighter(bson.ObjectId(""))
	ZeroFight   = Fight(bson.ObjectId(""))
	StartHex    = Hex(0)
)

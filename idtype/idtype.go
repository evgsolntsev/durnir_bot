package idtype

import (
	"math/rand"

	"github.com/globalsign/mgo/bson"
)

type Fighter bson.ObjectId

func NewFighter() Fighter {
	return Fighter(bson.NewObjectId())
}

type Player bson.ObjectId

func NewPlayer() Player {
	return Player(bson.NewObjectId())
}

type Fight bson.ObjectId

func NewFight() Fight {
	return Fight(bson.NewObjectId())
}

type Hex int

func NewHex() Hex {
	return Hex(rand.Intn(1000))
}

var (
	ZeroPlayer  = Player(bson.ObjectId(""))
	ZeroFighter = Fighter(bson.ObjectId(""))
)

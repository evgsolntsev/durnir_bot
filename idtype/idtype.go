package idtype

import (
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

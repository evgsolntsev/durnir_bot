package fighter

import (
	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Card struct {
	Type int `bson:"type"`
}

type Fighter struct {
	ID        idtype.Fighter `bson:"_id,omitempty"`
	Health    int            `bson:"health"`
	Mana      int            `bson:"mana"`
	Shield    int            `bson:"shield"`
	Will      int            `bson:"will"`
	Power     int            `bson:"power"`
	FearPower int            `bson:"fearPower"`
	Hex       int            `bson:"hex"`
	Deck      []Card         `bson:"deck"`
}

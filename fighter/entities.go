package fighter

import (
	"context"
	"math/rand"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

var cardTypeName = map[int]string{
	0: "Лечение",
	1: "Удар",
}

type Card struct {
	Type int `bson:"type"`
}

type Fighter struct {
	ID        idtype.Fighter `bson:"_id,omitempty"`
	Name      string         `bson:"name"`
	Health    int            `bson:"health"`
	Mana      int            `bson:"mana"`
	Shield    int            `bson:"shield"`
	Will      int            `bson:"will"`
	Power     int            `bson:"power"`
	FearPower int            `bson:"fearPower"`

	JoinFight bool           `bson:"joinFight"`
	Hex       idtype.Hex     `bson:"hex"`
	Deck      []Card         `bson:"deck"`
}

func (f *Fighter) GetCard(ctx context.Context) Card {
	i := rand.Intn(len(f.Deck))
	return f.Deck[i]
}

package fighter

import (
	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Fighter struct {
	ID        idtype.FighterID `json:"_id" bson:"_id"`
	Health    int              `json:"health" bson:"health"`
	Mana      int              `json:"mana" bson:"mana"`
	Shield    int              `json:"shield" bson:"shield"`
	Will      int              `json:"will" bson:"will"`
	Power     int              `json:"power" bson:"power"`
	FearPower int              `json:"fearPower" bson:"fearPower"`
}

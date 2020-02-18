package player

import "github.com/evgsolntsev/durnir_bot/idtype"

type Player struct {
	ID        idtype.Player   `bson:"_id"`
	Name      string          `bson:"string"`
	FighterID *idtype.Fighter `bson:"fighterId"`
}

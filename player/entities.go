package player

import (
	"github.com/evgsolntsev/durnir_bot/fighter"
	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Player struct {
	ID         idtype.Player   `bson:"_id"`
	TelegramId int64             `bson:"telegramId"`
	Name       string          `bson:"string"`
	FighterID  *idtype.Fighter `bson:"fighterId"`
	Gold       int             `bson:"gold"`
	Parts      []fighter.Part  `bson:"parts"`
}

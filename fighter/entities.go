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

type Card int

var (
	UnknownCard = Card(0)
	CardHeal    = Card(1)
	CardHit     = Card(2)
	CardSkip    = Card(3)
)

func (c Card) Name() string {
	switch c {
	case CardHeal:
		return "Лечение"
	case CardHit:
		return "Удар"
	case CardSkip:
		return "Пропуск"
	default:
		return "!@#$"
	}
}

type Fraction int

var (
	UnknownFraction  = Fraction(0)
	FractionPlayers  = Fraction(1)
	FractionMonsters = Fraction(2)
)

type Fighter struct {
	ID   idtype.Fighter `bson:"_id,omitempty"`
	Name string         `bson:"name"`

	Health    int `bson:"health"`
	Mana      int `bson:"mana"`
	Shield    int `bson:"shield"`
	Will      int `bson:"will"`
	Power     int `bson:"power"`
	FearPower int `bson:"fearPower"`

	Parts []Part `bson:"parts"`
	Gold  int    `bson:"gold"`

	JoinFight bool       `bson:"joinFight"`
	Hex       idtype.Hex `bson:"hex"`
	Deck      []Card     `bson:"deck"`
	Fraction  Fraction   `bson:"fraction"`
}

func (f *Fighter) GetCard(ctx context.Context) Card {
	i := rand.Intn(len(f.Deck))
	return f.Deck[i]
}

type Part int

var (
	UnknownPart = Part(0)
	PartBeak    = Part(1)
	PartBrain   = Part(2)
	PartWing    = Part(3)
	PartHand    = Part(4)
)

func (p Part) Value() int {
	switch p {
	case PartBeak:
		return 100
	case PartBrain:
		return 5000
	default:
		return 0
	}
}

func (p Part) Name() string {
	switch p {
	case PartBeak:
		return "Клюв"
	case PartBrain:
		return "Мозг"
	case UnknownPart:
		return "Неизвестная деталь"
	default:
		return "Неизвестная часть"
	}
}

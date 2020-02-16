package fight

import (
	"time"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Fight struct {
	ID          idtype.Fight   `bson:"_id"`
	Fighters    []FighterState `bson:"fighterIDs"`
	UpdatedTime time.Time      `bson:"updatedTime"`
	Started     bool           `bson:"started"`
	Hex         int            `bson:"hex"`
}

type Debuff int

var (
	DebuffPoisoned = Debuff(0)
)

type FighterState struct {
	ID      idtype.Fighter `bson:"fighter_id"`
	Health  int            `bson:"health"`
	Mana    int            `bson:"mana"`
	Debuffs []Debuff       `bson:"states,omitempty"`
}

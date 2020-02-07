package fight

import (
	"time"

	"github.com/evgsolntsev/durnir_bot/idtype"
)

type Fight struct {
	ID          idtype.Fight     `bson:"_id"`
	FighterIDs  []idtype.Fighter `bson:"fighterIDs"`
	UpdatedTime time.Time        `bson:"updatedTime"`
	Started     bool             `bson:"started"`
	Hex         int              `bson:"hex"`
}

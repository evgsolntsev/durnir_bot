package fighter

import (
	"testing"

	"github.com/globalsign/mgo/bson"
	"github.com/stretchr/testify/require"
)

func TestFighterMarshalling(t *testing.T) {
	f := &Fighter{
		Health:    100,
		Mana:      100,
		Shield:    2,
		Will:      3,
		Power:     4,
		FearPower: 3,
		Hex:       5,
		Deck:      []Card{{Type: 1}, {Type: 2}},
	}
	data, err := bson.Marshal(f)
	require.Nil(t, err)

	var newF Fighter
	err = bson.Unmarshal(data, &newF)
	require.Nil(t, err)
	require.Equal(t, *f, newF)
}

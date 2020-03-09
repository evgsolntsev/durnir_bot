package fighter

import (
	"context"
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
		Deck:      []Card{CardHeal, CardHeal},
		Parts:     []Part{},
	}
	data, err := bson.Marshal(f)
	require.Nil(t, err)

	var newF Fighter
	err = bson.Unmarshal(data, &newF)
	require.Nil(t, err)
	require.Equal(t, *f, newF)
}

func TestFighterGetCard(t *testing.T) {
	cards := []Card{CardHeal, CardHit, CardSkip}

	f := &Fighter{
		Deck: cards,
	}

	ctx := context.Background()
	result := make(map[Card]bool)
	for i := 0; i < 100; i++ {
		card := f.GetCard(ctx)
		result[card] = true
	}

	require.Len(t, result, 3)
}

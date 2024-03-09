package mtgapi

import (
	"time"

	client "github.com/MagicTheGathering/mtg-sdk-go"
	"github.com/rs/zerolog/log"
)

var sets []*client.Set
var lastUpdated time.Time

func GetAllSets() ([]*client.Set, error) {
	if len(sets) > 0 && lastUpdated.Add(time.Hour*24).After(time.Now()) {
		return sets, nil
	}
	newSets, err := client.NewSetQuery().All()
	if err != nil {
		return nil, err
	}
	lastUpdated = time.Now()
	sets = newSets
	return sets, nil
}

func GenerateBooster(set client.SetCode) ([]*client.Card, error) {
	log.Debug().Interface("set", set).Send()
	cards, err := set.GenerateBooster()
	if err != nil {
		return nil, err
	}
	return cards, nil
}

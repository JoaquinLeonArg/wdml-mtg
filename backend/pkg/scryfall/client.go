package scryfall

import (
	"context"
	"time"

	scryfallapi "github.com/BlueMonday/go-scryfall"
)

var sets []scryfallapi.Set
var lastUpdated time.Time
var client *scryfallapi.Client

func GetAllSets() ([]scryfallapi.Set, error) {
	if len(sets) > 0 && lastUpdated.Add(time.Hour*24).After(time.Now()) {
		return sets, nil
	}

	var err error
	if client == nil {
		client, err = scryfallapi.NewClient()
		if err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	newSets, err := client.ListSets(ctx)
	if err != nil {
		return nil, err
	}
	sets = newSets
	lastUpdated = time.Now()
	return sets, nil

}

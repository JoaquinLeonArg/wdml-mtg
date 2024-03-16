package scryfall

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/BlueMonday/go-scryfall"
	scryfallapi "github.com/BlueMonday/go-scryfall"
	lru "github.com/hashicorp/golang-lru/v2"
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
	for i, set := range sets {
		sets[i].Code = strings.ToUpper(set.Code)
	}
	lastUpdated = time.Now()
	return sets, nil

}

var setCards = make(map[string][]scryfallapi.Card)

func GetSetCards(set string) ([]scryfallapi.Card, error) {
	cachedCards, found := setCards[set]
	if found {
		return cachedCards, nil
	}

	var err error
	if client == nil {
		client, err = scryfallapi.NewClient()
		if err != nil {
			return nil, err
		}
	}

	page := 1
	ctx := context.Background()
	allCards := []scryfall.Card{}

	for {
		setData, err := client.SearchCards(ctx, fmt.Sprintf("set=%s", set), scryfallapi.SearchCardsOptions{Page: page})
		if err != nil {
			return nil, err
		}
		allCards = append(allCards, setData.Cards...)
		if !setData.HasMore {
			break
		}
		page += 1
	}

	setCards[set] = allCards

	return allCards, nil
}

var cachedPossibleCards *lru.Cache[string, []scryfallapi.Card]

func GetAllCardsByFilter(filter string) ([]scryfallapi.Card, error) {
	var err error
	if cachedPossibleCards == nil {
		cachedPossibleCards, err = lru.New[string, []scryfallapi.Card](64)
	}
	if err != nil {
		return nil, err
	}

	if cards, ok := cachedPossibleCards.Get(filter); ok {
		return cards, nil
	}

	if client == nil {
		client, err = scryfallapi.NewClient()
		if err != nil {
			return nil, err
		}
	}

	page := 1
	ctx := context.Background()
	allCards := []scryfall.Card{}

	for {
		setData, err := client.SearchCards(ctx, filter, scryfallapi.SearchCardsOptions{Page: page})
		if err != nil {
			return nil, err
		}
		allCards = append(allCards, setData.Cards...)
		if !setData.HasMore {
			break
		}
		page += 1
	}
	cachedPossibleCards.Add(filter, allCards)
	return allCards, nil
}

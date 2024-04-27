package scryfall

import (
	"context"
	"fmt"
	"slices"
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
		if err != nil && !strings.Contains(err.Error(), "not_found") {
			return nil, err
		}
		allCards = append(allCards, setData.Cards...)
		if !setData.HasMore {
			break
		}
		page += 1
	}
	if len(allCards) == 0 {
		var newFilter string
		if slices.Contains(strings.Split(filter, " "), "rarity:special") || slices.Contains(strings.Split(filter, " "), "rarity:mythic") {
			newFilter = strings.Replace(strings.Replace(filter, "rarity:special", "rarity:rare", 1), "rarity:mythic", "rarity:rare", 1)
			page = 1
			for {
				setData, err := client.SearchCards(ctx, newFilter, scryfallapi.SearchCardsOptions{Page: page})
				if err != nil {
					return nil, err
				}
				allCards = append(allCards, setData.Cards...)
				if !setData.HasMore {
					break
				}
				page += 1
			}
		}
		cachedPossibleCards.Add(newFilter, allCards)
	} else {
		cachedPossibleCards.Add(filter, allCards)

	}
	return allCards, nil
}

type CardsByIdentifier struct {
	Identifier scryfallapi.CardIdentifier
	Amount     int
}
type ScryfallCollectionRequest struct {
	Identifiers []scryfallapi.CardIdentifier `json:"identifiers"`
}

func GetAllCardsByIdentifiers(scryfallRequestBody ScryfallCollectionRequest) ([]scryfallapi.Card, error) {
	var err error

	if client == nil {
		client, err = scryfallapi.NewClient()
		if err != nil {
			return nil, err
		}
	}

	ctx := context.Background()

	response, err := client.GetCardsByIdentifiers(ctx, scryfallRequestBody.Identifiers)
	if err != nil {
		return nil, err
	}
	allCards := response.Data
	return allCards, nil
}

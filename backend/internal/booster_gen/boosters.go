package boostergen

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	scryfallapi "github.com/BlueMonday/go-scryfall"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/joaquinleonarg/wdml_mtg/backend/pkg/scryfall"
	"github.com/rs/zerolog/log"
)

type BoosterData struct {
	CardCount  int      `json:"cardCount"`
	Sets       []string `json:"sets"`
	DefaultSet string   `json:"defaultSet"`
	Slots      []struct {
		Options []Option `json:"options"`
		Count   int      `json:"count"`
	} `json:"slots"`
}

type Option struct {
	Rarity string `json:"rarity"`
	Weight int    `json:"weight"`
	Type   string `json:"type"`
	Layout string `json:"layout"`
	Set    string `json:"set"`
}

func GenerateBoosterFromJson(setCode string) ([]domain.CardData, error) {
	// Vanilla gen

	var boosterData BoosterData
	path, err := filepath.Abs("./internal/booster_gen/sets/" + strings.ToLower(setCode) + ".json")
	if err != nil {
		return nil, err
	}
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("Json read error")
	}
	err = json.Unmarshal(jsonData, &boosterData)
	if err != nil {
		return nil, fmt.Errorf("Json unmarshal error")
	}

	// log.Debug().Interface("data", boosterData).Send()

	cardList := make(map[string][]scryfallapi.Card)

	for _, sc := range boosterData.Sets {
		cards, err := scryfall.GetSetCards(sc)
		if err != nil {
			return nil, fmt.Errorf("scryfall error")
		}
		if _, ok := cardList[sc]; !ok {
			cardList[sc] = []scryfallapi.Card{}
		}
		cardList[sc] = append(cardList[sc], cards...)
	}

	if err != nil || len(cardList[strings.ToLower(setCode)]) == 0 {
		log.Debug().Str("set", setCode).Err(err).Msg("failed to generate booster pack")
		return nil, fmt.Errorf("No cards error")
	}

	boosterPack := make([]domain.CardData, 0, boosterData.CardCount)

	for _, slot := range boosterData.Slots {
		optionsByWeight := make(map[int]Option)
		currentWeight := 0
		for _, option := range slot.Options {
			currentWeight += option.Weight
			optionsByWeight[currentWeight] = option
		}
		for i := 0; i < slot.Count; i++ {

			chosenWeight := rand.Int() % currentWeight
			var chosenOption Option
			for w, option := range optionsByWeight {
				if chosenWeight < w {
					chosenOption = option
					break
				}
			}
			var set string

			if chosenOption.Set != "" {
				set = chosenOption.Set
			} else {
				set = boosterData.DefaultSet
			}
			possibleCards := make([]scryfallapi.Card, 0)
			for _, card := range cardList[set] {
				if card.Rarity == chosenOption.Rarity {
					if chosenOption.Type != "" && !strings.Contains(card.TypeLine, chosenOption.Type) {
						continue
					} else {
						if chosenOption.Layout != "" && card.Layout != scryfallapi.Layout(chosenOption.Layout) {
							continue
						} else {
							possibleCards = append(possibleCards, card)
						}
					}
				}
			}

			card := possibleCards[rand.Int()%len(possibleCards)]
			colors := []string{}
			for _, col := range card.Colors {
				colors = append(colors, string(col))
			}

			log.Debug().Interface("Selected card ", card.Name).Send()
			boosterPack = append(boosterPack,
				domain.CardData{
					Name:      card.Name,
					Typeline:  card.TypeLine,
					ManaValue: int(card.CMC),
					Colors:    colors,
					ImageURL:  card.ImageURIs.Large,
				},
			)
		}
	}

	log.Debug().Interface("booster", boosterPack).Send()

	return boosterPack, nil
}

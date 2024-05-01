package boostergen

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	scryfallapi "github.com/BlueMonday/go-scryfall"
	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/joaquinleonarg/wdml-mtg/backend/pkg/scryfall"
	"github.com/rs/zerolog/log"
)

type CardListsBySet map[string][]scryfallapi.Card

type BoosterDataGetter func(setCode string) (*domain.BoosterPack, error)

func GenerateBooster(setCode string, genFunc BoosterDataGetter) ([]domain.CardData, error) {
	boosterData, err := genFunc(setCode)
	if err != nil {
		return nil, err
	}
	boosterPack := make([]domain.CardData, 0, boosterData.CardCount)

	for _, slot := range boosterData.Slots {
		optionKeys := []int{}
		optionsByWeight := make(map[int]domain.Option)
		currentWeight := 0
		for _, option := range slot.Options {
			currentWeight += option.Weight
			optionsByWeight[currentWeight] = option
			optionKeys = append(optionKeys, currentWeight)
		}
		for i := 0; i < slot.Count; i++ {
			chosenOption := domain.Option{}
			if currentWeight > 0 {
				chosenWeight := rand.Int() % currentWeight
				for _, w := range optionKeys {
					if chosenWeight < w {
						chosenOption = optionsByWeight[w]
						break
					}
				}
			}

			filter := fmt.Sprintf("%s %s %s", boosterData.Filter, slot.Filter, chosenOption.Filter)

			cards, err := scryfall.GetAllCardsByFilter(filter)
			if err != nil || len(cards) == 0 {
				log.Debug().Str("set", setCode).Str("filter", filter).Err(err).Msg("failed to generate booster pack")
				return nil, fmt.Errorf("no cards error")
			}
			card := cards[rand.Int()%len(cards)]

			colors := []string{}
			for _, col := range card.Colors {
				colors = append(colors, string(col))
			}

			cardFront, cardBack := scryfall.GetImageFromFaces(card)
			boosterPack = append(boosterPack,
				domain.CardData{
					SetCode:         strings.ToUpper(card.Set),
					CollectorNumber: card.CollectorNumber,
					Name:            card.Name,
					Oracle:          card.OracleText,
					Rarity:          domain.CardRarity(card.Rarity),
					Types:           scryfall.ParseScryfallTypeline(card.TypeLine),
					ManaValue:       int(math.Floor(card.CMC)),
					ManaCost:        card.ManaCost,
					Colors:          colors,
					ImageURL:        cardFront,
					BackImageURL:    cardBack,
				},
			)
		}
	}

	return boosterPack, nil
}

func GetBoosterDataFromJson(setCode string) (*domain.BoosterPack, error) {
	var boosterData domain.BoosterPack
	path, err := filepath.Abs(fmt.Sprintf(("./internal/booster_gen/sets/%s.json"), strings.ToLower(setCode)))
	if err != nil {
		return nil, err
	}
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &boosterData)
	if err != nil {
		return nil, err
	}
	return &boosterData, nil
}

func GetBoosterDataFromDb(setCode string) (*domain.BoosterPack, error) {
	boosterPack, err := db.GetPackBySetCode(setCode)
	if err != nil {
		return nil, err
	}
	return boosterPack, nil
}

func GetBoosterDataPassthrough(boosterData domain.BoosterPack) BoosterDataGetter {
	return func(setCode string) (*domain.BoosterPack, error) {
		return &boosterData, nil
	}
}

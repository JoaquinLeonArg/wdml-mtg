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
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/joaquinleonarg/wdml_mtg/backend/pkg/scryfall"
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
		optionsByWeight := make(map[int]domain.Option)
		currentWeight := 0
		for _, option := range slot.Options {
			currentWeight += option.Weight
			optionsByWeight[currentWeight] = option
		}
		for i := 0; i < slot.Count; i++ {
			chosenOption := domain.Option{}
			if currentWeight > 0 {
				chosenWeight := rand.Int() % currentWeight
				for w, option := range optionsByWeight {
					if chosenWeight < w {
						chosenOption = option
						break
					}
				}
			}

			filter := fmt.Sprintf("%s %s game:paper", slot.Filter, chosenOption.Filter)

			cards, err := scryfall.GetAllCardsByFilter(filter)
			if err != nil || len(cards) == 0 {
				log.Debug().Str("set", setCode).Err(err).Msg("failed to generate booster pack")
				return nil, fmt.Errorf("no cards error")
			}
			card := cards[rand.Int()%len(cards)]

			colors := []string{}
			for _, col := range card.Colors {
				colors = append(colors, string(col))
			}

			cardFront, cardBack := scryfall.GetImageFromFaces(card)
			log.Debug().Interface(fmt.Sprintf("Selected card for slot %v", i), card.Name).Send()
			boosterPack = append(boosterPack,
				domain.CardData{
					SetCode:         strings.ToUpper(card.Set),
					CollectorNumber: card.CollectorNumber,
					Name:            card.Name,
					Oracle:          card.OracleText,
					Rarity:          domain.CardRarity(card.Rarity),
					Types:           scryfall.ParseScryfallTypeline(card.TypeLine),
					ManaValue:       int(math.Floor(card.CMC)),
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

// func CheckIfBoosterExists(setCode string) (bool, error) {
// 	path, err := filepath.Abs("./internal/booster_gen/sets/" + strings.ToLower(setCode) + ".json")
// 	if err != nil {
// 		return false, err
// 	}
// 	_, err = os.ReadFile(path)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

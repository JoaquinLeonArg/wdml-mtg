package scryfall

import (
	"strings"

	"github.com/BlueMonday/go-scryfall"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
)

func ParseScryfallTypeline(rawType string) []string {
	var types []string

	for _, rawType1 := range strings.Split(rawType, "—") {
		for _, rawType2 := range strings.Split(rawType1, " ") {
			rawType3 := strings.Trim(rawType2, " ")
			if rawType3 != "" && rawType3 != "//" {
				types = append(types, rawType3)
			}
		}
	}

	return types
}

func GetImageFromFaces(card scryfall.Card) (string, string) {
	if card.CardFaces != nil && card.CardFaces[0].ImageURIs.Normal != "" && card.CardFaces[1].ImageURIs.Normal != "" {
		return card.CardFaces[0].ImageURIs.Normal, card.CardFaces[1].ImageURIs.Normal
	} else {
		return card.ImageURIs.Normal, ""
	}
}

func GetCardDataFromScryCard(card scryfall.Card) domain.CardData {
	colors := []string{}
	for _, col := range card.Colors {
		colors = append(colors, string(col))
	}
	types := ParseScryfallTypeline(card.TypeLine)

	newCard := domain.CardData{
		SetCode:         strings.ToUpper(card.Set),
		CollectorNumber: card.CollectorNumber,
		Name:            card.Name,
		Oracle:          card.OracleText,
		Rarity:          domain.CardRarity(card.Rarity),
		Types:           types,
		ManaValue:       int(card.CMC),
		ManaCost:        card.ManaCost,
		Colors:          colors,
	}
	newCard.ImageURL, newCard.BackImageURL = GetImageFromFaces(card)
	return newCard
}

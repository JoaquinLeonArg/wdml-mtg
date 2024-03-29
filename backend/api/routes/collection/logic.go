package collection

import (
	"slices"
	"strconv"
	"strings"

	scryfallapi "github.com/BlueMonday/go-scryfall"
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	"github.com/joaquinleonarg/wdml_mtg/backend/pkg/scryfall"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCollectionCards(userID, tournamentID, filters string, count, page int) ([]domain.OwnedCard, int, error) {
	log.Debug().Str("filters", filters).Send()
	dbFilters := []db.CardFilter{}
	for _, filter := range strings.Split(filters, "+") {
		for _, filterOperation := range []db.CardFilterOperation{
			db.CardFilterOperationEq,
			db.CardFilterOperationLt,
			db.CardFilterOperationGt,
		} {
			for _, filterType := range []db.CardFilterType{
				db.CardFilterTypeName,
				db.CardFilterTypeTags,
				db.CardFilterTypeRarity,
				db.CardFilterTypeColor,
				db.CardFilterTypeTypes,
				db.CardFilterTypeOracle,
				db.CardFilterTypeSetCode,
				db.CardFilterTypeMV,
			} {
				if strings.Contains(filter, string(filterOperation)) {
					splitted := strings.Split(filter, string(filterOperation))
					if splitted[1] == "" {
						continue
					}
					if splitted[0] == string(filterType) {
						dbFilters = append(dbFilters, db.CardFilter{
							Type:      filterType,
							Operation: filterOperation,
							Value:     splitted[1],
						})
					}
					continue
				}
			}
		}
	}
	return db.GetCardsFromTournamentPlayer(userID, tournamentID, dbFilters, count, page)
}

func GetOwnedCardById(cardId string) (domain.OwnedCard, error) {
	return db.GetOwnedCardById(cardId)
}

func ImportCollection(importCardCsv [][]string, userID, tournamentID string) error {
	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return apiErrors.ErrBadRequest
	}
	dbUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return apiErrors.ErrBadRequest
	}
	cards := make([]domain.CardData, 0, len(importCardCsv))
	cardsBySetCode := make([]scryfall.CardsByIdentifier, 0, len(importCardCsv))
	cum := 0
	scryfallRequestBody := scryfall.ScryfallCollectionRequest{Identifiers: []scryfallapi.CardIdentifier{}}
	// Pos 3 = set, Pos 9 = collectors number
	for i, card := range importCardCsv {
		if i > 0 {
			amt, err := strconv.Atoi(card[0])
			if err != nil {
				return err
			}
			cardsBySetCode = append(cardsBySetCode, scryfall.CardsByIdentifier{Identifier: scryfallapi.CardIdentifier{Set: card[3], CollectorNumber: card[9]}, Amount: amt})
			scryfallRequestBody.Identifiers = append(scryfallRequestBody.Identifiers, scryfallapi.CardIdentifier{Set: card[3], CollectorNumber: card[9]})

			cum += 1
			// When cum == 75, query scryfall and add cards to slice of OwnedCard
			if cum == 75 {
				cum = 0

				// Pegada a scry
				scryCardData, err := scryfall.GetAllCardsByIdentifiers(scryfallRequestBody)
				if err != nil {
					return err
				}
				for _, card := range scryCardData {
					cardData := scryfall.GetCardDataFromScryCard(card)
					cards = append(cards, cardData)
				}
				scryfallRequestBody.Identifiers = []scryfallapi.CardIdentifier{}
			}
		}
	}
	if cum > 0 {
		cum = 0
		// Request to Scryfall
		scryCardData, err := scryfall.GetAllCardsByIdentifiers(scryfallRequestBody)
		if err != nil {
			return err
		}
		for _, card := range scryCardData {
			cardData := scryfall.GetCardDataFromScryCard(card)
			cards = append(cards, cardData)
		}
		scryfallRequestBody.Identifiers = []scryfallapi.CardIdentifier{}

	}

	ownedCards := make([]domain.OwnedCard, 0, len(cardsBySetCode))
	coinsToAdd := 0
	for _, cardIdent := range cardsBySetCode {
		var foundCard domain.CardData
		for _, cardData := range cards {
			if strings.ToLower(cardData.SetCode) == cardIdent.Identifier.Set && cardData.CollectorNumber == cardIdent.Identifier.CollectorNumber {
				foundCard = cardData
				if cardIdent.Amount > 4 && !slices.Contains(cardData.Types, "Basic") {
					coins := (cardIdent.Amount - 4)
					switch foundCard.Rarity {
					case "mythic":
						coinsToAdd += coins * domain.MYTHIC_TO_COIN
					case "rare":
						coinsToAdd += coins * domain.RARE_TO_COIN
					case "uncommon":
						coinsToAdd += coins * domain.UNCOMMON_TO_COIN
					case "common":
						coinsToAdd += coins * domain.COMMON_TO_COIN
					}
					cardIdent.Amount = 4
				}
			}
		}
		newOwnedCard := domain.OwnedCard{
			ID:           primitive.NewObjectID(),
			CardData:     foundCard,
			TournamentID: dbTournamentID,
			UserID:       dbUserID,
			Count:        cardIdent.Amount,
		}
		ownedCards = append(ownedCards, newOwnedCard)
	}
	db.AddCoinsToTournamentPlayer(coinsToAdd, userID, tournamentID)
	return db.ImportCollection(ownedCards)
}

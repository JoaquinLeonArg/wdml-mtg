package deck

import (
	"errors"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
)

func GetDeckById(deckID string) (*domain.Deck, []domain.OwnedCard, error) {
	return db.GetDeckByID(deckID)
}

func GetDecksForTournamentPlayer(tournamentID, userID string) ([]domain.Deck, error) {
	tournamentPlayer, err := db.GetTournamentPlayer(tournamentID, userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, apiErrors.ErrInternal
	}
	return db.GetDecksForTournamentPlayer(tournamentPlayer.ID.Hex())
}

func CreateEmptyDeck(ownerID, deckName, deckDescription, tournamentID string) error {
	tournamentPlayer, err := db.GetTournamentPlayer(tournamentID, ownerID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	createdDeck := domain.Deck{
		Name:               deckName,
		Description:        deckDescription,
		TournamentPlayerID: tournamentPlayer.ID,
	}

	err = db.CreateEmptyDeck(createdDeck)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return apiErrors.ErrDuplicatedResource
		}
		return apiErrors.ErrInternal
	}

	return nil
}

func AddOwnedCardToDeck(cardID string, deckID string, amount int, board domain.DeckBoard) error {
	return db.AddOwnedCardToDeck(
		cardID,
		deckID,
		amount,
		board,
	)
}

func RemoveCardFromDeck(ownedCardID, deckID string, board domain.DeckBoard, amount int) error {
	return db.RemoveDeckCardFromDeck(
		ownedCardID,
		deckID,
		board,
		amount,
	)
}

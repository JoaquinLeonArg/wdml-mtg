package deck

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
)

func GetDeckById(id string) (*domain.Deck, error) {
	return db.GetDeckById(id)
}

func GetDecksByTournamentPlayerId(id string) ([]domain.Deck, error) {
	return db.GetDecksByTournamentPlayerId(id)
}

func CreateEmptyDeck(ownerId string, req CreateEmptyDeckRequest) error {
	// Get TournamentPlayerId from params

	tournamentPlayer, err := db.GetTournamentPlayer(req.TournamentId, ownerId)
	if err != nil {
		return err
	}

	deck := domain.Deck{
		Name:               req.Deck.Name,
		Description:        req.Deck.Description,
		TournamentPlayerID: tournamentPlayer.ID,
	}

	return db.CreateEmptyDeck(deck)
}

func AddOwnedCardToDeck(addOwnedCardToDeckRequest AddOwnedCardToDeckRequest) error {
	return db.AddOwnedCardToDeck(
		addOwnedCardToDeckRequest.Card,
		addOwnedCardToDeckRequest.DeckId,
		addOwnedCardToDeckRequest.Amount,
		addOwnedCardToDeckRequest.Board,
	)
}

func RemoveCardFromDeck(addOwnedCardToDeckRequest RemoveCardFromDeckRequest) error {
	return db.RemoveDeckCardFromDeck(
		addOwnedCardToDeckRequest.Card,
		addOwnedCardToDeckRequest.Amount,
		addOwnedCardToDeckRequest.Board,
	)
}

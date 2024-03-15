package deck

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
)

func GetDeckById(deckId string) (*domain.Deck, error) {
	return db.GetDeckById(deckId)
}

func GetDecksByTournamentPlayerId(tournamentPlayerId string) ([]domain.Deck, error) {
	return db.GetDecksByTournamentPlayerId(tournamentPlayerId)
}

func CreateEmptyDeck(ownerId, deckName, deckDescription, tournamentId string) error {
	tournamentPlayer, err := db.GetTournamentPlayer(tournamentId, ownerId)
	if err != nil {
		return err
	}

	createdDeck := domain.Deck{
		Name:               deckName,
		Description:        deckDescription,
		TournamentPlayerID: tournamentPlayer.ID,
	}

	return db.CreateEmptyDeck(createdDeck)
}

func AddOwnedCardToDeck(cardId string, deckId string, amount int, board domain.DeckBoard) error {
	return db.AddOwnedCardToDeck(
		cardId,
		deckId,
		amount,
		board,
	)
}

func RemoveCardFromDeck(card domain.DeckCard, deckId string, amount int, board domain.DeckBoard) error {
	return db.RemoveDeckCardFromDeck(
		card,
		deckId,
		amount,
		board,
	)
}

package tournament

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTournamentPlayerByID(tournamentPlayerID string) (*domain.TournamentPlayer, error) {
	return db.GetTournamentPlayerByID(tournamentPlayerID)
}

func CreateTournamentPlayer(rawUserID string, createTournamentPlayerRequest CreateTournamentPlayerRequest) (string, error) {
	userID, err := primitive.ObjectIDFromHex(rawUserID)
	if err != nil {
		return "", db.ErrInvalidID
	}
	tournament, err := db.GetTournamentByInviteCode(createTournamentPlayerRequest.TournamentCode)
	if err != nil {
		return "", err
	}
	tournamentID, err := db.CreateTournamentPlayer(domain.TournamentPlayer{
		UserID:       userID,
		TournamentID: tournament.ID,
		AccessLevel:  domain.AccessLevelPlayer,
		GameResources: domain.GameResources{
			OwnedCards: []domain.Card{},
			Decks:      []domain.Deck{},
			Wildcards: domain.OwnedWildcards{
				CommonCount:      0,
				UncommonCount:    0,
				RareCount:        0,
				MythicRareCount:  0,
				MasterpieceCount: 0,
			},
			BoosterPacks: []domain.OwnedBoosterPack{},
			Rerolls:      0,
			Coins:        0,
		},
		TournamentPoints: 0,
	})
	if err != nil {
		return "", err
	}

	return tournamentID.Hex(), nil
}

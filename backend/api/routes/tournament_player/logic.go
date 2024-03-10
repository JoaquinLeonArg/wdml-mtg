package tournament_player

import (
	"errors"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTournamentPlayerByID(tournamentPlayerID string) (*domain.TournamentPlayer, error) {
	return db.GetTournamentPlayerByID(tournamentPlayerID)
}

func GetTournamentPlayersForUser(userID string) ([]domain.TournamentPlayer, error) {
	return db.GetTournamentPlayersForUser(userID)
}

func CreateTournamentPlayer(rawUserID string, createTournamentPlayerRequest CreateTournamentPlayerRequest) (string, error) {
	userID, err := primitive.ObjectIDFromHex(rawUserID)
	if err != nil {
		return "", apiErrors.ErrInternal
	}
	tournament, err := db.GetTournamentByInviteCode(createTournamentPlayerRequest.TournamentCode)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return "", apiErrors.ErrNotFound
		}
		return "", err
	}
	tournamentID, err := db.CreateTournamentPlayer(domain.TournamentPlayer{
		UserID:       userID,
		TournamentID: tournament.ID,
		AccessLevel:  domain.AccessLevelPlayer,
		GameResources: domain.GameResources{
			Decks: []domain.Deck{},
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

package tournament_player

import (
	"errors"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTournamentPlayers(tournamentID string) ([]domain.TournamentPlayer, []domain.User, error) {
	tournament_players, users, err := db.GetTournamentPlayers(tournamentID)

	// Redact sensitive information
	for index := range users {
		users[index].Password = nil
		users[index].Email = ""
	}

	return tournament_players, users, err
}

func GetTournamentPlayerByID(tournamentPlayerID string) (*domain.TournamentPlayer, error) {
	return db.GetTournamentPlayerByID(tournamentPlayerID)
}

func GetTournamentPlayerByAuth(tournamentID, userID string) (*domain.TournamentPlayer, error) {
	return db.GetTournamentPlayer(tournamentID, userID)
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
		return "", apiErrors.ErrInternal
	}
	// TODO: Consolidate this creation with the other one
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
		if errors.Is(err, db.ErrNotFound) {
			return "", apiErrors.ErrNotFound
		}
		if errors.Is(err, db.ErrAlreadyExists) {
			return "", apiErrors.ErrDuplicatedResource
		}
		return "", apiErrors.ErrInternal
	}

	return tournamentID.Hex(), nil
}

func AddCoinsToTournamentPlayer(tournamentPlayerID string, coins int) error {
	tPlayer, err := db.GetTournamentPlayerByID(tournamentPlayerID)
	if err != nil {
		return err
	}
	return db.AddCoinsToTournamentPlayer(coins, tPlayer.UserID.Hex(), tPlayer.TournamentID.Hex())
}

func AddPointsToTournamentPlayer(tournamentPlayerID string, coins int) error {
	tPlayer, err := db.GetTournamentPlayerByID(tournamentPlayerID)
	if err != nil {
		return err
	}
	return db.AddPointsToTournamentPlayer(coins, tPlayer.UserID.Hex(), tPlayer.TournamentID.Hex())
}

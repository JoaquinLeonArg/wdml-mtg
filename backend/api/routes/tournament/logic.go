package tournament

import (
	"errors"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
)

func GetTournamentByID(tournamentID string) (*domain.Tournament, error) {
	tournament, err := db.GetTournamentByID(tournamentID)
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	return tournament, nil
}

func GetTournamentsForUser(userID string) ([]domain.Tournament, error) {
	tournaments, err := db.GetTournamentsForUser(userID)
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	return tournaments, nil
}

func GetTournamentPlayers(tournamentID string) ([]domain.TournamentPlayer, []domain.User, error) {
	tournament_players, users, err := db.GetTournamentPlayers(tournamentID)

	// Redact sensitive information
	for index := range users {
		users[index].Password = nil
		users[index].Email = ""
	}

	return tournament_players, users, err
}

func CreateTournament(tournament domain.Tournament) (string, error) {
	tournamentID, err := db.CreateTournament(tournament)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return "", apiErrors.ErrDuplicatedResource
		}
		if errors.Is(err, db.ErrNotFound) {
			return "", apiErrors.ErrNotFound
		}
		return "", apiErrors.ErrInternal
	}
	_, err = db.CreateTournamentPlayer(
		domain.TournamentPlayer{
			UserID:       tournament.OwnerID,
			TournamentID: tournamentID,
			AccessLevel:  domain.AccessLevelAdministrator,
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
		},
	)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return "", apiErrors.ErrDuplicatedResource
		}
		if errors.Is(err, db.ErrNotFound) {
			return "", apiErrors.ErrNotFound
		}
		return "", apiErrors.ErrInternal
	}
	return tournamentID.Hex(), nil
}

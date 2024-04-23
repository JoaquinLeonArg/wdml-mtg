package user

import (
	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
)

func GetTournamentsForUser(userID string) ([]domain.Tournament, error) {
	tournaments, err := db.GetTournamentsForUser(userID)
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	return tournaments, nil
}

func GetTournamentPlayersForUser(userID string) ([]domain.TournamentPlayer, error) {
	tournament_players, err := db.GetTournamentPlayersForUser(userID)
	if err != nil {
		return nil, apiErrors.ErrInternal
	}
	return tournament_players, nil
}

package tournament_post

import (
	"errors"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
)

func GetTournamentPosts(tournamentID string) ([]domain.TournamentPost, error) {
	tournamentPosts, err := db.GetAllTournamentPosts(tournamentID)
	if err != nil {
		return nil, apiErrors.ErrInternal
	}

	return tournamentPosts, nil
}

func AddTournamentPost(userID, tournamentID string, tournamentPost domain.TournamentPost) error {
	// Check if tournament player has permissions
	tournament_player, err := db.GetTournamentPlayer(tournamentID, userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	tournamentPost.TournamentPlayerID = tournament_player.ID

	if tournament_player.AccessLevel != domain.AccessLevelAdministrator && tournament_player.AccessLevel != domain.AccessLevelModerator {
		return apiErrors.ErrUnauthorized
	}

	err = db.CreateTournamentPost(tournamentPost, tournamentID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	return nil
}

func DeleteTournamentPost(userID, tournamentID, tournamentPostID string) error {
	// Check if tournament player has permissions
	tournament_player, err := db.GetTournamentPlayer(tournamentID, userID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	if tournament_player.AccessLevel != domain.AccessLevelAdministrator && tournament_player.AccessLevel != domain.AccessLevelModerator {
		return apiErrors.ErrUnauthorized
	}

	err = db.DeleteTournamentPost(tournamentID, tournamentPostID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}

	return nil
}

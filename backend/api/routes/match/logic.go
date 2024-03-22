package match

import (
	"errors"

	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
)

func CreateMatch(match domain.Match) error {
	err := db.CreateMatch(match)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}
	return nil
}

func UpdateMatch(match domain.Match) error {
	err := db.UpdateMatch(match)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}
	return nil
}

func GetMatchesFromSeason(seasonID string, onlyPending bool) ([]domain.Match, error) {
	matches, err := db.GetMatchesFromSeason(seasonID, onlyPending)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, apiErrors.ErrInternal
	}
	return matches, nil
}

func GetMatchesFromPlayer(playerID string, onlyPending bool, count, page int) ([]domain.Match, error) {
	matches, err := db.GetMatchesFromPlayer(playerID, onlyPending, count, page)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, apiErrors.ErrNotFound
		}
		return nil, apiErrors.ErrInternal
	}
	return matches, nil
}

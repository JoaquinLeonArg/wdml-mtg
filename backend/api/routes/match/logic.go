package match

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
)

func CreateMatch(match domain.Match) error {
	return db.CreateMatch(match)
}

func UpdateMatch(match domain.Match) error {
	return db.UpdateMatch(match)
}

func GetMatchesFromSeason(seasonID string, onlyPending bool, count, page int) ([]domain.Match, error) {
	return db.GetMatchesFromSeason(seasonID, onlyPending, count, page)
}

func GetMatchesFromPlayer(playerID string, onlyPending bool, count, page int) ([]domain.Match, error) {
	return db.GetMatchesFromPlayer(playerID, onlyPending, count, page)
}

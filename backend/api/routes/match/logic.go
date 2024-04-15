package match

import (
	"errors"

	"github.com/joaquinleonarg/wdml-mtg/backend/db"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml-mtg/backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateMatch(seasonID string, gamemode domain.Gamemode, tournamentPlayerIDs []string) error {
	if len(tournamentPlayerIDs) < 2 {
		return apiErrors.ErrBadRequest
	}
	_, err := db.GetSeasonByID(seasonID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}
	playersData := []domain.MatchPlayerData{}
	for _, playerID := range tournamentPlayerIDs {
		id, err := primitive.ObjectIDFromHex(playerID)
		if err != nil {
			return apiErrors.ErrBadRequest
		}
		_, err = db.GetTournamentPlayerByID(playerID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return apiErrors.ErrNotFound
			}
			return apiErrors.ErrInternal
		}
		playersData = append(playersData, domain.MatchPlayerData{
			TournamentPlayerID: id,
			Wins:               0,
			Tags:               []string{},
		})
	}
	match := domain.Match{
		PlayersData: playersData,
		GamesPlayed: 0,
		Gamemode:    gamemode,
		Completed:   false,
	}
	err = db.CreateMatch(seasonID, match)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apiErrors.ErrNotFound
		}
		return apiErrors.ErrInternal
	}
	return nil
}

func UpdateMatch(matchID string, playersPoints map[string]int, gamesPlayed int, completed bool) error {
	err := db.UpdateMatch(matchID, playersPoints, gamesPlayed, completed)
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

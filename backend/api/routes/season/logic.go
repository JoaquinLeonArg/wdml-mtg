package season

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllSeasons(tournamentID string) ([]domain.Season, error) {
	return db.GetAllSeasons(tournamentID)
}

func GetSeasonByID(seasonID string) (*domain.Season, error) {
	return db.GetSeasonByID(seasonID)
}

func CreateEmptySeason(name, description, tournamentID string) error {
	dbTournamentID, err := primitive.ObjectIDFromHex(tournamentID)
	if err != nil {
		return apiErrors.ErrInternal
	}
	return db.CreateEmptySeason(domain.Season{
		Name:         name,
		Description:  description,
		TournamentID: dbTournamentID,
	})
}

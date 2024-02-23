package tournament

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
)

func GetTournamentByID(tournamentID string) (*domain.Tournament, error) {
	return db.GetTournamentByID(tournamentID)
}

func CreateTournament(tournament domain.Tournament) (string, error) {
	tournamentID, err := db.CreateTournament(tournament)
	if err != nil {
		return "", err
	}
	return tournamentID.Hex(), nil
}

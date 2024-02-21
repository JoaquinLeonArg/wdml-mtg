package tournament

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
)

func GetTournamentByID(tournamentID string) (*domain.Tournament, error) {
	return db.GetTournamentByID(tournamentID)
}

func CreateNewTournament(tournament domain.Tournament) error {
	return db.CreateTournament(tournament)
}

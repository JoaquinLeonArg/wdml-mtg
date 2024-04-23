package user

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/users").Subrouter()
	r.HandleFunc("/{user_id}/tournaments", GetTournamentsForUserHandler).Methods(http.MethodGet)
	r.HandleFunc("/{user_id}/tournament-players", GetTournamentsForUserHandler).Methods(http.MethodGet)
}

// ENDPOINT: Get all tournaments for a user
type GetTournamentsForUserResponse struct {
	Tournaments []domain.Tournament `json:"tournaments"`
}

func GetTournamentsForUserHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(string)
	if userID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Get the tournament
	tournaments, err := GetTournamentsForUser(userID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournaments")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentsForUserResponse{Tournaments: tournaments}))
}

// ENDPOINT: Get tournament players for a user
type GetTournamentPlayersForUserResponse struct {
	TournamentPlayers []domain.TournamentPlayer `json:"tournament_players"`
}

func GetTournamentPlayersForUserHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from url
	userID, ok := mux.Vars(r)["user_id"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Try to get tournament players
	tournamentPlayers, err := GetTournamentPlayersForUser(userID)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament players")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentPlayersForUserResponse{TournamentPlayers: tournamentPlayers}))
}

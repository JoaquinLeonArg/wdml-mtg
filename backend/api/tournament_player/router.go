package tournament

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r.HandleFunc("/tournament_player/{tournamentPlayerId}", GetTournamentPlayerHandler).Methods(http.MethodGet)
	r.HandleFunc("/tournament", CreateTournamentPlayerHandler).Methods(http.MethodPost)
}

type GetTournamentPlayerResponse struct {
	TournamentPlayer domain.TournamentPlayer `json:"tournament_player"`
}

func GetTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get id
	vars := mux.Vars(r)
	tournamentPlayerID, ok := vars["tournamentPlayerId"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Try to get tournament player
	tournamentPlayer, err := GetTournamentPlayerByID(tournamentPlayerID)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament player")
		w.Write(response.NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentPlayerResponse{TournamentPlayer: *tournamentPlayer}))
}

type CreateTournamentPlayerRequest struct {
	TournamentCode string `json:"tournament_code"`
}

type CreateTournamentPlayerResponse struct {
	TournamentID string `json:"tournament_id"`
}

func CreateTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user id from request context
	userID, ok := r.Context().Value("user_id").(string)
	if userID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Decode body data
	var createTournamentPlayerRequest CreateTournamentPlayerRequest
	err := json.NewDecoder(r.Body).Decode(&createTournamentPlayerRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.Write(response.NewErrorResponse(err))
		return
	}

	tournamentID, err := CreateTournamentPlayer(userID, createTournamentPlayerRequest)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		w.Write(response.NewErrorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(CreateTournamentPlayerResponse{TournamentID: tournamentID}))
}

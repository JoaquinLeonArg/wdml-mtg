package tournament_player

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/auth"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/tournament_player").Subrouter()

	r.HandleFunc("", GetAllTournamentPlayersHandler).Methods(http.MethodGet)
	r.HandleFunc("/me", GetTournamentPlayersFromAuthHandler).Methods(http.MethodGet)
	r.HandleFunc("", CreateTournamentPlayerHandler).Methods(http.MethodPost)
	r.HandleFunc("/{tournament_player_id}/coins", AddCoinsToTournamentPlayerHandler).Methods(http.MethodPost)
	r.HandleFunc("/{tournament_player_id}/points", AddPointsToTournamentPlayerHandler).Methods(http.MethodPost)
}

// ENDPOINT: Get a tournament's players
type GetTournamentPlayersResponse struct {
	TournamentPlayers []domain.TournamentPlayer `json:"tournament_players"`
	Users             []domain.User             `json:"users"`
}

func GetAllTournamentPlayersHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from url
	tournamentID, ok := mux.Vars(r)["tournament_id"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Get the tournament players
	tournament_players, users, err := GetTournamentPlayers(tournamentID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentPlayersResponse{TournamentPlayers: tournament_players, Users: users}))
}

// ENDPOINT: Create tournament player on a tournament
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tournamentID, err := CreateTournamentPlayer(userID, createTournamentPlayerRequest)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(CreateTournamentPlayerResponse{TournamentID: tournamentID}))
}

type GetTournamentPlayerFromAuthResponse struct {
	TournamentPlayer domain.TournamentPlayer `json:"tournament_player"`
}

func GetTournamentPlayersFromAuthHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from context
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get tournament ID from url
	tournamentID, ok := mux.Vars(r)["tournament_id"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Try to get tournament players
	tournamentPlayers, err := GetTournamentPlayerByAuth(tournamentID, userID)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament players")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentPlayerFromAuthResponse{TournamentPlayer: *tournamentPlayers}))
}

type AddCoinsToTournamentPlayerRequest struct {
	Coins int `json:"coins"`
}

func AddCoinsToTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament player ID from url
	tournamentPlayerID, ok := mux.Vars(r)["tournament_player_id"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Decode body data
	var req AddCoinsToTournamentPlayerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	err = AddCoinsToTournamentPlayer(tournamentPlayerID, req.Coins)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(nil))
}

type AddPointsToTournamentPlayerRequest struct {
	Points int `json:"points"`
}

func AddPointsToTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament player ID from url
	tournamentPlayerID, ok := mux.Vars(r)["tournament_player_id"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Decode body data
	var req AddPointsToTournamentPlayerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	err = AddPointsToTournamentPlayer(tournamentPlayerID, req.Points)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(nil))
}

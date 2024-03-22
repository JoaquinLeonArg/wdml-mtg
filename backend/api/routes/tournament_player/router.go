package tournament_player

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/routes/auth"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	apiErrors "github.com/joaquinleonarg/wdml_mtg/backend/errors"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/tournament_player").Subrouter()
	r.HandleFunc("/boosters", GetPacksForTournamentPlayerHandler).Methods(http.MethodGet)
	r.HandleFunc("/user/{userID}", GetTournamentPlayersForUserHandler).Methods(http.MethodGet)
	r.HandleFunc("", GetTournamentPlayersFromAuthHandler).Methods(http.MethodGet)
	r.HandleFunc("/tournament", GetTournamentPlayer).Methods(http.MethodGet)
	r.HandleFunc("", CreateTournamentPlayerHandler).Methods(http.MethodPost)
	r.HandleFunc("/add/coins", AddCoinsToTournamentPlayerHandler).Methods(http.MethodPost)
	r.HandleFunc("/add/points", AddPointsToTournamentPlayerHandler).Methods(http.MethodPost)
}

type GetPacksForTournamentPlayerResponse struct {
	BoosterPacks []domain.OwnedBoosterPack `json:"booster_packs"`
}

func GetPacksForTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get user ID from request context
	userID, ok := r.Context().Value("user_id").(string)
	if userID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Get booster packs
	packs, err := GetBoosterPacksForTournamentPlayer(tournamentID, userID)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get booster packs")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetPacksForTournamentPlayerResponse{BoosterPacks: packs}))
}

type GetTournamentPlayerResponse struct {
	TournamentPlayer domain.TournamentPlayer `json:"tournament_player"`
}

func GetTournamentPlayer(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get user ID from request context
	userID, ok := r.Context().Value("user_id").(string)
	if userID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Get tournament players for this user
	// TODO: Filter on the DB
	tournamentPlayers, err := GetTournamentPlayersForUser(userID)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament player")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, tournamentPlayer := range tournamentPlayers {
		if tournamentPlayer.TournamentID.Hex() == tournamentID {
			w.WriteHeader(http.StatusOK)
			w.Write(response.NewDataResponse(GetTournamentPlayerResponse{TournamentPlayer: tournamentPlayer}))
			return
		}
	}

	w.Write(response.NewErrorResponse(apiErrors.ErrNotFound))
	w.WriteHeader(http.StatusNotFound)
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

type GetTournamentPlayersForUserResponse struct {
	TournamentPlayers []domain.TournamentPlayer `json:"tournament_players"`
}

func GetTournamentPlayersForUserHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get id
	vars := mux.Vars(r)
	userID, ok := vars["userID"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
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

type GetTournamentPlayersFromAuthResponse struct {
	TournamentPlayers []domain.TournamentPlayer `json:"tournament_players"`
}

func GetTournamentPlayersFromAuthHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get id
	userID, err := auth.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
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

type AddCoinsToTournamentPlayerRequest struct {
	Coins int `json:"coins"`
}

type EmptyResponse struct {
}

func AddCoinsToTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament player ID from request context
	tPlayerID := r.URL.Query().Get("tournament_player_id")
	if tPlayerID == "" {
		http.Error(w, "", http.StatusBadRequest)
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

	err = AddCoinsToTournamentPlayer(tPlayerID, req.Coins)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(EmptyResponse{}))
}

type AddPointsToTournamentPlayerRequest struct {
	Points int `json:"points"`
}

func AddPointsToTournamentPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament player ID from request context
	tPlayerID := r.URL.Query().Get("tournament_player_id")
	if tPlayerID == "" {
		http.Error(w, "", http.StatusBadRequest)
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

	err = AddPointsToTournamentPlayer(tPlayerID, req.Points)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(EmptyResponse{}))
}

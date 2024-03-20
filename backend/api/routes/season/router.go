package season

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/season").Subrouter()
	r.HandleFunc("/all", GetAllSeasonsHandler).Methods(http.MethodGet)
	r.HandleFunc("", GetSeasonByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("", CreateEmptySeasonHandler).Methods(http.MethodPost)
	r.HandleFunc("/match", CreateMatchHandler).Methods(http.MethodPost)
	r.HandleFunc("/match", UpdateMatchHandler).Methods(http.MethodPut)
	r.HandleFunc("/match", GetMatchesFromSeasonHandler).Methods(http.MethodGet)
	r.HandleFunc("/match/by-player", GetMatchesFromPlayerHandler).Methods(http.MethodGet)
}

type GetSeasonsResponse struct {
	Seasons []domain.Season `json:"seasons"`
}
type GetMatchesResponse struct {
	Matches []domain.Match `json:"matches"`
}

type EmptyResponse struct{}

func GetAllSeasonsHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get season
	seasons, err := GetAllSeasons(tournamentID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get seasons")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetSeasonsResponse{Seasons: seasons}))
}

func GetSeasonByIDHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()
	seasonID := r.URL.Query().Get("season_id")
	if seasonID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	season, err := GetSeasonByID(seasonID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get deck data")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(season))
}

type CreateEmptySeasonRequest struct {
	TournamentID string `bson:"tournament_id" json:"tournament_id"`
	Name         string `bson:"name" json:"name"`
	Description  string `bson:"description" json:"description"`
}

func CreateEmptySeasonHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Decode body data
	var createEmptySeasonRequest CreateEmptySeasonRequest
	err := json.NewDecoder(r.Body).Decode(&createEmptySeasonRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Get tournament ID from body
	tournamentID := createEmptySeasonRequest.TournamentID
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Create a season with no matches
	err = CreateEmptySeason(
		createEmptySeasonRequest.Name,
		createEmptySeasonRequest.Description,
		tournamentID,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(EmptyResponse{}))
}

type CreateMatchRequest struct {
	SeasonID     string `json:"season_id"`
	TournamentID string `json:"tournament_id"`
	PlayerAID    string `json:"player_a_id"`
	PlayerBID    string `json:"player_b_id"`
}

func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Decode body data
	var createMatchRequest CreateMatchRequest
	err := json.NewDecoder(r.Body).Decode(&createMatchRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Get tournament ID from body
	tournamentID := createMatchRequest.TournamentID
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Add a match to given season
	err = CreateMatch(
		createMatchRequest.SeasonID,
		tournamentID,
		createMatchRequest.PlayerAID,
		createMatchRequest.PlayerBID,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(EmptyResponse{}))
}

type UpdateMatchRequest struct {
	SeasonID string       `json:"season_id"`
	Match    domain.Match `json:"match"`
}

func UpdateMatchHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Decode body data
	var updateMatchRequest UpdateMatchRequest
	err := json.NewDecoder(r.Body).Decode(&updateMatchRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Update a given match with results
	err = UpdateMatch(
		updateMatchRequest.SeasonID,
		updateMatchRequest.Match,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(EmptyResponse{}))
}

func GetMatchesFromSeasonHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get season ID from query
	seasonID := r.URL.Query().Get("season_id")
	if seasonID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}
	onlyPending := false
	onlyPendingQuery := r.URL.Query().Get("pending")
	if onlyPendingQuery == "true" {
		onlyPending = true
	}

	// Get matches
	matches, err := GetMatchesFromSeason(seasonID, onlyPending)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get seasons")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetMatchesResponse{Matches: matches}))
}

func GetMatchesFromPlayerHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}
	// Get player ID from query
	playerID := r.URL.Query().Get("player_id")
	if playerID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}
	onlyPending := false
	onlyPendingQuery := r.URL.Query().Get("pending")
	if onlyPendingQuery == "true" {
		onlyPending = true
	}

	// Get matches
	matches, err := GetMatchesFromPlayer(tournamentID, playerID, onlyPending)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get seasons")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetMatchesResponse{Matches: matches}))
}

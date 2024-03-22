package match

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/match").Subrouter()

	r.HandleFunc("", CreateMatchHandler).Methods(http.MethodPost)
	r.HandleFunc("", UpdateMatchHandler).Methods(http.MethodPut)
	r.HandleFunc("", GetMatchesFromSeasonHandler).Methods(http.MethodGet)
	r.HandleFunc("/by-player", GetMatchesFromPlayerHandler).Methods(http.MethodGet)
}

type GetMatchesResponse struct {
	Matches []domain.Match `json:"matches"`
}

type EmptyResponse struct{}

type CreateMatchRequest struct {
	SeasonID  string          `json:"season_id"`
	Gamemode  domain.Gamemode `json:"gamemode"`
	PlayerIDs []string        `json:"player_ids"`
}

func CreateMatchHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Decode body data
	var req CreateMatchRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Add a match to given season
	err = CreateMatch(req.SeasonID, req.Gamemode, req.PlayerIDs)
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
	PlayersPoints map[string]int `json:"players_points"`
	GamesPlayed   int            `json:"games_played"`
	Completed     bool           `json:"completed"`
}

func UpdateMatchHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get Match ID from query
	matchID := r.URL.Query().Get("match_id")
	if matchID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

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
		matchID,
		updateMatchRequest.PlayersPoints,
		updateMatchRequest.GamesPlayed,
		updateMatchRequest.Completed,
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

	// Get Season ID from query
	// TODO: Also query by tournament id
	seasonID := r.URL.Query().Get("season_id")
	if seasonID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}
	onlyPending := false
	onlyPendingQuery := r.URL.Query().Get("pending")
	if onlyPendingQuery != "" {
		val, err := strconv.ParseBool(onlyPendingQuery)
		if err != nil {
			log.Debug().
				Msg("failed to read pending from query")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		onlyPending = val
	}

	// Get matches
	matches, err := GetMatchesFromSeason(seasonID, onlyPending)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get matches")
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

	// Get player ID from query
	playerID := r.URL.Query().Get("tournament_player_id")
	if playerID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}
	onlyPending := false
	onlyPendingQuery := r.URL.Query().Get("pending")
	if onlyPendingQuery != "" {
		val, err := strconv.ParseBool(onlyPendingQuery)
		if err != nil {
			log.Debug().
				Msg("failed to read count from query")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		onlyPending = val
	}

	page := 0
	pageQuery := r.URL.Query().Get("page")
	if pageQuery != "" {
		val, err := strconv.Atoi(pageQuery)
		if err != nil {
			log.Debug().
				Msg("failed to read page from query")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		page = val
	}

	count := 0
	countQuery := r.URL.Query().Get("count")
	if countQuery != "" {
		val, err := strconv.Atoi(countQuery)
		if err != nil {
			log.Debug().
				Msg("failed to read count from query")
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		count = max(val, 75)
	}

	// Get matches
	matches, err := GetMatchesFromPlayer(playerID, onlyPending, count, page)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get matches")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetMatchesResponse{Matches: matches}))
}

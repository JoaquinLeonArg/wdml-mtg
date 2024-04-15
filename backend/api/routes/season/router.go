package season

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/season").Subrouter()
	r.HandleFunc("/all", GetAllSeasonsHandler).Methods(http.MethodGet)
	r.HandleFunc("", GetSeasonByIDHandler).Methods(http.MethodGet)
	r.HandleFunc("", CreateEmptySeasonHandler).Methods(http.MethodPost)
}

type GetSeasonsResponse struct {
	Seasons []domain.Season `json:"seasons"`
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

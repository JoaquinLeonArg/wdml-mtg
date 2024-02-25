package tournament

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/tournament").Subrouter()
	r.HandleFunc("/{tournamentId}", GetTournamentHandler).Methods(http.MethodGet)
	r.HandleFunc("", CreateTournamentHandler).Methods(http.MethodPost)
}

type GetTournamentHandlerResponse struct {
	Tournament domain.Tournament `json:"tournament"`
}

func GetTournamentHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get id
	vars := mux.Vars(r)
	tournamentID, ok := vars["tournamentId"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Try to get tournament
	tournament, err := GetTournamentByID(tournamentID)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentHandlerResponse{Tournament: *tournament}))
}

type CreateTournamentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateTournamentResponse struct {
	TournamentID string `json:"tournament_id"`
}

func CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user id from request context
	rawOwnerID, ok := r.Context().Value("user_id").(string)
	if rawOwnerID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	ownerID, err := primitive.ObjectIDFromHex(rawOwnerID)
	if err != nil {
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Decode body data
	var createTournamentRequest CreateTournamentRequest
	err = json.NewDecoder(r.Body).Decode(&createTournamentRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tournamentID, err := CreateTournament(domain.Tournament{
		OwnerID:     ownerID,
		Name:        createTournamentRequest.Name,
		Description: createTournamentRequest.Description,
	})

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(CreateTournamentResponse{TournamentID: tournamentID}))
}

package tournament

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/tournament/tournament_player"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/tournaments").Subrouter()
	tournament_player.RegisterEndpoints(r.PathPrefix("/{tournament_id}/tournament-players").Subrouter())

	r.HandleFunc("/{tournament_id}", GetTournamentHandler).Methods(http.MethodGet)
	r.HandleFunc("", CreateTournamentHandler).Methods(http.MethodPost)
}

type GetTournamentResponse struct {
	Tournament domain.Tournament `json:"tournament"`
}

// GetTournament godoc
//
//	@Summary Get tournament
//	@Description Returns a tournament given it's id.
//	@Tags tournament
//	@Accept json
//	@Produce json
func GetTournamentHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from url
	tournamentID, ok := mux.Vars(r)["tournament_id"]
	if !ok {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	// Get the tournament
	tournament, err := GetTournamentByID(tournamentID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentResponse{Tournament: *tournament}))
}

// ENDPOINT: Create a new tournament
type CreateTournamentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateTournamentResponse struct {
	TournamentID string `json:"tournament_id"`
}

func CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from context
	rawOwnerID, ok := r.Context().Value("user_id").(string)
	if rawOwnerID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Get the DB object ID for the user
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

	// Create the tournament
	// TODO: Have this struct be created on the logic layer
	tournamentID, err := CreateTournament(domain.Tournament{
		OwnerID:     ownerID,
		Name:        createTournamentRequest.Name,
		Description: createTournamentRequest.Description,
	})
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		w.Write(response.NewErrorResponse(err))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(CreateTournamentResponse{TournamentID: tournamentID}))
}

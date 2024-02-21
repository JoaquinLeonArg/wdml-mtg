package tournament

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterEndpoints(r *mux.Router) {
	r.HandleFunc("/tournament/{tournamentId}", GetTournamentHandler).Methods(http.MethodGet)
	r.HandleFunc("/tournament", CreateTournamentHandler).Methods(http.MethodPost)
}

func GetTournamentHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get id
	vars := mux.Vars(r)
	tournamentID, ok := vars["tournamentId"]
	if !ok {
		http.Error(w, "missing tournament id", http.StatusBadRequest)
	}

	// Try to get tournament
	tournament, err := GetTournamentByID(tournamentID)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encode tournament data
	data, err := json.Marshal(tournament)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

type CreateTournamentRequest struct {
}

func CreateTournamentHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user id from request context
	rawOwnerID, ok := r.Context().Value("user_id").(string)
	if rawOwnerID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "auth failed", http.StatusForbidden)
		return
	}

	ownerID, err := primitive.ObjectIDFromHex(rawOwnerID)

	// Decode body data
	var createTournamentRequest CreateTournamentRequest
	err = json.NewDecoder(r.Body).Decode(&createTournamentRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = CreateNewTournament(domain.Tournament{
		OwnerID:     ownerID,
		Title:       createTournamentRequest.Title,
		Description: createTournamentRequest.Description,
	})

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Encode tournament data
	data, err := json.Marshal(tournament)

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

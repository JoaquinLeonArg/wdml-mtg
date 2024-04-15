package tournament

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/tournament").Subrouter()
	r.HandleFunc("", GetTournamentHandler).Methods(http.MethodGet)
	r.HandleFunc("/user", GetTournamentsForUserHandler).Methods(http.MethodGet)
	r.HandleFunc("", CreateTournamentHandler).Methods(http.MethodPost)
	r.HandleFunc("/tournament_player", GetTournamentPlayersHandler).Methods(http.MethodGet)
}

//
// ENDPOINT: Get a tournament by it's id
//

type GetTournamentResponse struct {
	Tournament domain.Tournament `json:"tournament"`
}

func GetTournamentHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
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

//
// ENDPOINT: Get all tournaments for a user
//

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

//
// ENDPOINT: Get a tournament's players
//

type GetTournamentPlayersResponse struct {
	TournamentPlayers []domain.TournamentPlayer `json:"tournament_players"`
	Users             []domain.User             `json:"users"`
}

func GetTournamentPlayersHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
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

//
// ENDPOINT: Create a new tournament
//

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

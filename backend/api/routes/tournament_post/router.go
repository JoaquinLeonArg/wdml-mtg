package tournament_post

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/joaquinleonarg/wdml-mtg/backend/domain"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/tournament_post").Subrouter()
	r.HandleFunc("", GetTournamentPostsHandler).Methods(http.MethodGet)
	r.HandleFunc("", AddTournamentPostHandler).Methods(http.MethodPost)
	r.HandleFunc("/remove", DeleteTournamentPostHandler).Methods(http.MethodPost)
}

//
// ENDPOINT: Get all posts on a tournament
//

type GetTournamentPostsResponse struct {
	TournamentPosts []domain.TournamentPost `json:"tournament_posts"`
}

func GetTournamentPostsHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Get all available booster packs for this tournament
	tournamentPosts, err := GetTournamentPosts(tournamentID)
	if err != nil {
		log.Debug().Err(err).Msg("failed to get tournament posts")
		w.WriteHeader(http.StatusInternalServerError)
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(GetTournamentPostsResponse{TournamentPosts: tournamentPosts}))
}

//
// ENDPOINT: Add tournament post
//

type AddTournamentPostRequest struct {
	TournamentPost domain.TournamentPost `json:"tournament_post"`
}

type AddTournamentPostResponse struct{}

func AddTournamentPostHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from request context
	userID, ok := r.Context().Value("user_id").(string)
	if userID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Decode body data
	var req AddTournamentPostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Add the tournament post
	err = AddTournamentPost(userID, tournamentID, req.TournamentPost)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(AddTournamentPostResponse{}))
}

//
// ENDPOINT: Delete a tournament post
//

type DeleteTournamentPostRequest struct {
	TournamentPostID string `json:"tournament_post_id"`
}

type DeleteTournamentPostResponse struct{}

func DeleteTournamentPostHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()

	// Get user ID from request context
	userID, ok := r.Context().Value("user_id").(string)
	if userID == "" || !ok {
		log.Debug().
			Msg("failed to read user id from context")
		http.Error(w, "", http.StatusForbidden)
		return
	}

	// Decode body data
	var req DeleteTournamentPostRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Get tournament ID from query
	tournamentID := r.URL.Query().Get("tournament_id")
	if tournamentID == "" {
		http.Error(w, "", http.StatusBadRequest)
	}

	// Delete the tournament post
	err = DeleteTournamentPost(userID, tournamentID, req.TournamentPostID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Send response back
	w.WriteHeader(http.StatusOK)
	w.Write(response.NewDataResponse(DeleteTournamentPostResponse{}))
}

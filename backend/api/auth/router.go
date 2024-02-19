package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r.HandleFunc("/auth/login", Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/register", Register).Methods(http.MethodPost)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Logic

	res, _ := json.Marshal("okay")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()
	// Decode body data
	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Try to create user
	err = CreateUser(registerRequest)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to create user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	log.Debug().
		Str("username", registerRequest.Username).
		Str("email", registerRequest.Email).
		Msg("created new user")
	w.WriteHeader(http.StatusOK)
	w.Write(nil)

}

package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/response"
	"github.com/rs/zerolog/log"
)

func RegisterEndpoints(r *mux.Router) {
	r = r.PathPrefix("/auth").Subrouter()
	r.HandleFunc("/login", LoginHandler).Methods(http.MethodPost)
	r.HandleFunc("/register", RegisterHandler).Methods(http.MethodPost)
	r.HandleFunc("/check", CheckHandler).Methods(http.MethodGet)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct{}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()
	// Decode body data
	var loginRequest LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Try to login user
	jwt, err := LoginUser(loginRequest)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to login user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}
	log.Debug().
		Str("username", loginRequest.Username).
		Msg("logged in user")
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		MaxAge:   36000,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
	w.Write(response.NewDataResponse(LoginResponse{}))
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct{}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	log := log.With().Ctx(r.Context()).Str("path", r.URL.Path).Logger()
	// Decode body data
	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		log.Debug().
			Err(err).
			Msg("failed to read request body")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}

	// Try to create user
	err = CreateUser(registerRequest)

	// Write response
	if err != nil {
		log.Debug().Err(err).Msg("failed to create user")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response.NewErrorResponse(err))
		return
	}
	log.Debug().
		Str("username", registerRequest.Username).
		Str("email", registerRequest.Email).
		Msg("created new user")
	w.Write(response.NewDataResponse(RegisterResponse{}))
	w.WriteHeader(http.StatusOK)
}

func CheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(response.NewDataResponse(nil))
	w.WriteHeader(http.StatusOK)
}

package auth

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(r *mux.Router) {
	r.HandleFunc("/auth/login", Login).Methods(http.MethodPost)
	r.HandleFunc("/auth/register", Register).Methods(http.MethodPost)
}

func Login(w http.ResponseWriter, r *http.Request) {
	// TODO: Logic

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	// Decode body data
	var registerRequest RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&registerRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	CreateUser(registerRequest)

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

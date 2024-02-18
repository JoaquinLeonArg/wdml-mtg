package users

import (
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterEndpoints(r *mux.Router) {
	r.HandleFunc("/users", GetUsers).Methods(http.MethodGet)
	r.HandleFunc("/users/{user_id}", GetUserById).Methods(http.MethodGet)
	r.HandleFunc("/users", CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/users", PatchUser).Methods(http.MethodPatch)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// TODO: Logic

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func GetUserById(w http.ResponseWriter, r *http.Request) {
	// mux.Vars(r)["user_id"]
	// TODO: Logic

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Logic
	// r.Body.Read(&body)

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

func PatchUser(w http.ResponseWriter, r *http.Request) {
	// TODO: Logic
	// r.Body.Read(&body)

	w.WriteHeader(http.StatusOK)
	w.Write(nil)
}

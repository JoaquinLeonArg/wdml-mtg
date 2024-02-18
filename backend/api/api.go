package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/auth"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/users"
	"github.com/joaquinleonarg/wdml_mtg/backend/config"
)

func StartServer() {
	r := mux.NewRouter()
	r.Use(auth.AuthMiddleware)

	auth.RegisterEndpoints(r)
	users.RegisterEndpoints(r)

	http.ListenAndServe(fmt.Sprintf(":%v", config.Config.ApiPort), r)
}

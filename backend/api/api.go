package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
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

	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // TODO: Set to a more sensible value for security reasons
	credentialsOk := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	http.ListenAndServe(fmt.Sprintf(":%v", config.Config.ApiPort), handlers.CORS(originsOk, credentialsOk, headersOk, methodsOk)(r))
}

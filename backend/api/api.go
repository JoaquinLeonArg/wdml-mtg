package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/routes/auth"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/routes/tournament"
	"github.com/joaquinleonarg/wdml_mtg/backend/api/routes/tournament_player"
	"github.com/joaquinleonarg/wdml_mtg/backend/config"
)

func StartServer() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	authRouter := router.NewRoute().Subrouter()
	auth.RegisterEndpoints(authRouter)

	router.Use(auth.AuthMiddleware)
	tournament.RegisterEndpoints(router)
	tournament_player.RegisterEndpoints(router)

	originsOk := handlers.AllowedOrigins([]string{"http://localhost:3000"}) // TODO: Set to a more sensible value for security reasons
	credentialsOk := handlers.AllowCredentials()
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, _ := route.GetPathTemplate()
		met, _ := route.GetMethods()
		fmt.Println(tpl, met)
		return nil
	})

	http.ListenAndServe(fmt.Sprintf(":%v", config.Config.ApiPort), handlers.CORS(originsOk, credentialsOk, headersOk, methodsOk)(router))
}

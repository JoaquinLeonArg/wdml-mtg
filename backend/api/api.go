package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/auth"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/boosterpacks"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/collection"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/deck"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/match"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/season"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/tournament"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/tournament_player"
	"github.com/joaquinleonarg/wdml-mtg/backend/api/routes/tournament_post"
	"github.com/joaquinleonarg/wdml-mtg/backend/config"
)

func StartServer() {
	router := mux.NewRouter().PathPrefix("/api").Subrouter()

	authRouter := router.NewRoute().Subrouter()
	auth.RegisterEndpoints(authRouter)

	router.Use(auth.AuthMiddleware)
	tournament.RegisterEndpoints(router)
	tournament_player.RegisterEndpoints(router)
	boosterpacks.RegisterEndpoints(router)
	collection.RegisterEndpoints(router)
	deck.RegisterEndpoints(router)
	season.RegisterEndpoints(router)
	match.RegisterEndpoints(router)
	tournament_post.RegisterEndpoints(router)

	originsOk := handlers.AllowedOrigins([]string{config.Config.CorsOrigin})
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

package main

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/api"
	"github.com/joaquinleonarg/wdml_mtg/backend/config"
	"github.com/rs/zerolog/log"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Panic().
			Err(err).
			Msg("failed to load config")
	}
	log.Info().
		Int("port", config.Config.ApiPort).
		Msg("starting server")
	api.StartServer()
}

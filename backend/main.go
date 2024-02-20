package main

import (
	"github.com/joaquinleonarg/wdml_mtg/backend/api"
	"github.com/joaquinleonarg/wdml_mtg/backend/config"
	"github.com/joaquinleonarg/wdml_mtg/backend/db"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	err := config.Load()
	if err != nil {
		log.Panic().
			Err(err).
			Msg("failed to load config")
	}

	err = db.InitDBConnection()
	if err != nil {
		log.Panic().
			Err(err).
			Msg("failed to init db connection")
	}

	log.Info().
		Int("port", config.Config.ApiPort).
		Msg("starting server")
	api.StartServer()
}

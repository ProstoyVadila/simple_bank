package main

import (
	"database/sql"

	"github.com/ProstoyVadila/simple_bank/api"
	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	setLogger(zerolog.InfoLevel)

	log.Info().Msg("Loading config")
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Can't load config file")
	}

	log.Info().Msg("Connecting to database")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to db")
	}

	setGinMode(config)
	store := db.NewStore(conn)
	server := api.NewServer(store)

	log.Info().Msg("Starting server")
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't start server")
	}
}

func setLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = utils.TimeFormat
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func setGinMode(config utils.Config) {
	gin.SetMode(config.GinMode)
}

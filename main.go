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

func init() {
	setLogger(zerolog.InfoLevel)
}

func main() {
	log.Info().Msg("loading config")
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("can't load config file")
	}

	log.Info().Msg("connecting to database")
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to db")
	}

	setGinMode(config)
	store := db.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("can't create server")
	}

	log.Info().Msg("starting server")
	log.Info().Msg("Listening and serving HTTP on " + config.ServerAddress)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("can't start server")
	}
}

func setLogger(level zerolog.Level) {
	zerolog.TimeFieldFormat = utils.TimeFormat
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func setGinMode(config utils.Config) {
	gin.SetMode(config.GinMode)
}

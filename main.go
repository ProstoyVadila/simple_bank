package main

import (
	"database/sql"
	"log"

	"github.com/ProstoyVadila/simple_bank/api"
	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/utils"

	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config file", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Can't connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Can't start server", err)
	}
}

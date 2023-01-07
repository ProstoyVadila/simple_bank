package main

import (
	"database/sql"
	"log"

	"github.com/ProstoyVadila/simple_bank/api"
	db "github.com/ProstoyVadila/simple_bank/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://admni:stopmining@localhost:5432/data?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Can't connect to db", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(serverAddress)
	if err != nil {
		log.Fatal("Can't start server", err)
	}
}

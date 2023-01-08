package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/ProstoyVadila/simple_bank/utils"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := utils.LoadConfig("../..")
	if err != nil {
		log.Fatalf("Can't load config file: %v", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("Can't connect to database: %v", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}

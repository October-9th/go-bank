package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/October-9th/simple-bank/util"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("Couldn't load config file: ", err)
	}
	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to database: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())

}

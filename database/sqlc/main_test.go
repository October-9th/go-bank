package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Println("Couldn't load env file")
	}
	dbDriver, dbSource := os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE")
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("Couldn't connect to database: ", err)
	}
	testQueries = New(testDB)
	os.Exit(m.Run())
}

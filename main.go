package main

import (
	"database/sql"
	"log"

	"github.com/October-9th/simple-bank/api"
	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/October-9th/simple-bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Couldn't load config file: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Couldn't connect to database: ", err)
	}

	store := sqlc.NewStore(conn)
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}
	log.Println("HTTP server served at:", config.ServerAddress)
	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("Couldn't start server: ", err)
	}

}

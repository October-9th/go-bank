package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/October-9th/simple-bank/api"
	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/October-9th/simple-bank/gapi"
	"github.com/October-9th/simple-bank/pb"
	"github.com/October-9th/simple-bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	runGrpcServer(config, store)

}
func runGrpcServer(config util.Config, store sqlc.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't create server", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoBankServer(grpcServer, server)
	// Register the reflection for the gprc server for gRPC client to explore what RPCs are available on the server
	// and how to call them
	reflection.Register(grpcServer)

	// Define the listener
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("Couldn't create listener: ", err)
	}
	log.Printf("GRPC server served at %s", config.GRPCServerAddress)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal("Couldn't start server: ", err)
	}
}
func runGinServer(config util.Config, store sqlc.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("Couldn't create server: ", err)
	}

	log.Println("HTTP server served at:", config.HTTPServerAddress)
	if err = server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("Couldn't start server: ", err)
	}
}

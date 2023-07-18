package gapi

import (
	"fmt"

	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/October-9th/simple-bank/pb"
	"github.com/October-9th/simple-bank/token"
	"github.com/October-9th/simple-bank/util"
	"github.com/gin-gonic/gin"
)

// Server serves HTTP request for banking service
type Server struct {
	pb.UnimplementedGoBankServer
	config     util.Config
	store      sqlc.Store
	tokenMaker token.Maker
}

// NewServer create a new GRPC server to server gRPC request and setup routing.
func NewServer(config util.Config, store sqlc.Store) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmectricKey)
	if err != nil {
		return nil, fmt.Errorf("couldn't create token maker: %v", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	gin.SetMode(gin.ReleaseMode)

	return server, nil
}

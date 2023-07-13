package api

import (
	"fmt"

	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/October-9th/simple-bank/token"
	"github.com/October-9th/simple-bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for banking service
type Server struct {
	config     util.Config
	store      sqlc.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

// NewServer create a new HTTP server and setup routing.
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
	// Register custom validator with Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// Add routes to router
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Routes for handler account api request
	router.POST("/api/v1/users", server.CreateUser)
	router.POST("/api/v1/users/login", server.loginUser)
	router.POST("/api/v1/users/renew_access", server.renewAccessToken)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/api/v1/accounts", server.createAccount)
	authRoutes.GET("/api/v1/accounts/:id", server.getAccount)
	authRoutes.GET("/api/v1/accounts", server.getListAccount)
	authRoutes.PUT("/api/v1/accounts", server.updateAccount)
	authRoutes.DELETE("/api/v1/accounts/:id", server.deleteAccount)

	// Routes for hanlder transfer api request
	authRoutes.POST("/api/v1/transfers", server.createTransfer)

	// Routes for handler usr api request

	server.router = router

}

// Start runs the HTTP server on a specific address and listening for API requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

package api

import (
	"github.com/October-9th/simple-bank/database/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server serves HTTP request for banking service
type Server struct {
	store  sqlc.Store
	router *gin.Engine
}

// NewServer create a new HTTP server and setup routing.
func NewServer(store sqlc.Store) *Server {
	server := &Server{store: store}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.SetTrustedProxies(nil)

	// Register custom validator with Gin
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// Add routes to router

	// Routes for handler account api request
	router.POST("/api/v1/accounts", server.createAccount)
	router.GET("/api/v1/accounts/:id", server.getAccount)
	router.GET("/api/v1/accounts", server.getListAccount)
	router.PUT("/api/v1/accounts", server.updateAccount)
	router.DELETE("/api/v1/accounts/:id", server.deleteAccount)

	// Routes for hanlder transfer api request
	router.POST("/api/v1/transfers", server.createTransfer)

	// Routes for handler usr api request
	router.POST("/api/v1/users", server.CreateUser)

	server.router = router
	return server
}

// Start runs the HTTP server on a specific address and listening for API requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

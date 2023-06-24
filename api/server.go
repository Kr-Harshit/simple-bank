package api

import (
	"log"

	db "github.com/KHarshit1203/simple-bank/db/gen"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := v.RegisterValidation("currency", validCurrency); err != nil {
			log.Fatalf("unable to register currency validator, %v", err)
		}
	}

	server.registerRoutes(router)

	server.router = router
	return server
}

func (server *Server) registerRoutes(router *gin.Engine) {
	router.POST("/api/accounts", server.createAccount)
	router.GET("/api/accounts/:id", server.getAccount)
	router.GET("/api/accounts", server.listAccount)
	router.DELETE("/api/accounts/:id", server.deleteAccount)
	router.DELETE("/api/accounts/purge/:owner-id", server.purgeUserAccounts)
	router.POST("/api/transfer", server.createTransfer)

	router.POST("/api/users", server.createUser)
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

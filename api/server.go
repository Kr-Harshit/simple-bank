package api

import (
	db "github.com/KHarshit1203/simple-bank/db/gen"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	server.registerRoutes(router)

	server.router = router
	return server
}

func (server *Server) registerRoutes(router *gin.Engine) {
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)
	router.DELETE("/accounts/purge/:owner-id", server.purgeUserAccounts)
}

// Start runs the HTTP server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errroResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

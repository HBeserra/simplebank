package api

import (
	db "github.com/HBeserra/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	Store  *db.Store
	router *gin.Engine
}

// NewServer create a new HTTP server and setup routing
func NewServer(store *db.Store) *Server {
	server := Server{
		Store:  store,
		router: gin.Default(),
	}

	server.router.POST("/accounts", server.createAccount)
	server.router.GET("/accounts/:id", server.getAccount)
	server.router.GET("/accounts", server.listAccounts)
	return &server
}

// Start runs the HTTP server on a specific address
func (s Server) Start(addr string) error {
	return s.router.Run(addr)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}

package api

import (
	"github.com/ProstoyVadila/simple_bank/api/middleware"
	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type ApiError struct {
	Field string
	Msg   string
}

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	server.router = gin.New()

	// register middleware
	middleware.Register(server.router)

	// register routes
	server.setRoutes()

	// register various validators for gin
	registerValidators()

	return server
}

// Start http server
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// errorResponse wraps error messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

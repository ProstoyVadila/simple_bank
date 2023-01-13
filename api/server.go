package api

import (
	"fmt"

	"github.com/ProstoyVadila/simple_bank/api/middleware"
	db "github.com/ProstoyVadila/simple_bank/db/sqlc"
	"github.com/ProstoyVadila/simple_bank/token"
	"github.com/ProstoyVadila/simple_bank/utils"
	"github.com/gin-gonic/gin"
)

const (
	maxEventsPerSec = 1000
	maxBurstSize    = 20
)

type Server struct {
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
	config     utils.Config
}

func NewServer(config utils.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPaseto(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}
	server.router = gin.New()

	server.setMiddlewares()
	server.setRoutes()
	server.setValidators()

	return server, nil
}

// Start http server
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

// setMiddlewares adds middlewares to router
func (s *Server) setMiddlewares() {
	s.router.Use(
		middleware.Recovery(),
		middleware.DefaultLogger(),
		middleware.CORS(),
		middleware.Throttling(maxEventsPerSec, maxBurstSize),
		middleware.Errors(),
	)
}

// errorResponse wraps error messages
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

package api

import "github.com/gin-gonic/gin"

func (s *Server) setRoutes() *gin.Engine {
	router := gin.Default()

	// accounts
	s.router.POST("/accounts", s.createAccount)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.GET("/accounts", s.listAccount)
	s.router.DELETE("/accouts/:id", s.deleteAccount)

	// transfers
	s.router.POST("/transfers", s.createTransfer)

	// users
	s.router.POST("/users", s.createUser)

	return router
}

package api

import "github.com/gin-gonic/gin"

func routes(s *Server) *gin.Engine {
	router := gin.Default()

	router.POST("/accounts", s.createAccount)
	router.GET("/accounts/:id", s.getAccount)
	router.GET("/accounts", s.listAccount)
	router.DELETE("/accouts/:id", s.deleteAccount)

	return router
}

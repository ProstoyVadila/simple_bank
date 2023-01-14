package api

import "github.com/ProstoyVadila/simple_bank/api/middleware"

// setRoutes adds routes to router
func (s *Server) setRoutes() {
	// users
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
	s.router.GET("/users", s.getUser)

	// routes under an authorization
	authRouters := s.router.Group("/").Use(middleware.Auth(s.tokenMaker))
	// accounts
	authRouters.POST("/accounts", s.createAccount)
	authRouters.GET("/accounts/:id", s.getAccount)
	authRouters.GET("/accounts", s.listAccount)

	// transfers
	authRouters.POST("/transfers", s.createTransfer)

}

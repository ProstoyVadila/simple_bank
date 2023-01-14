package api

// setRoutes adds routes to router
func (s *Server) setRoutes() {
	// accounts
	s.router.POST("/accounts", s.createAccount)
	s.router.GET("/accounts/:id", s.getAccount)
	s.router.GET("/accounts", s.listAccount)
	s.router.DELETE("/accouts/:id", s.deleteAccount)

	// transfers
	s.router.POST("/transfers", s.createTransfer)

	// users
	s.router.POST("/users", s.createUser)
	s.router.POST("/users/login", s.loginUser)
	s.router.GET("/users", s.getUser)
}

package app

func (s *Server) routes() {
	router := s.router

	// Authorization group
	authorized := router.Group("/")
	authorized.Use(s.ginContextToContext())
	authorized.Use(s.authGuard(s.resolver.authService))
	{
		// graphql route
		authorized.POST("/query", gqlHandler(s.resolver))
	}

}

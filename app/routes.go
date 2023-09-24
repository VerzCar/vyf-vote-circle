package app

func (s *Server) routes() {
	router := s.router

	// Service group
	v1 := router.Group("/v1/api/vote-circle")
	v1.Use(s.ginContextToContext())

	// Authorization group
	authorized := v1.Group("/")
	authorized.Use(s.authGuard(s.authService))
	{
		// circle
		authorized.GET("/circle", s.Circle())
		authorized.POST("/circle", s.CreateCircle())
		authorized.PUT("/circle", s.UpdateCircle())

		// votes
		authorized.POST("/vote", s.CreateVote())

		// rankings
		authorized.GET("/rankings", s.Rankings())
	}
}

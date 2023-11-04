package app

func (s *Server) routes() {
	router := s.router

	// Service group
	v1 := router.Group("/v1/api/vote-circle")

	// Authorization group
	authorized := v1.Group("/")
	authorized.Use(s.authGuard(s.authService))
	{
		// circle
		authorized.GET("/circle/:circleId", s.Circle())
		authorized.GET("/circle/:circleId/eligible", s.EligibleToBeInCircle())
		authorized.GET("/circles", s.Circles())
		authorized.POST("/circle", s.CreateCircle())
		authorized.PUT("/circle", s.UpdateCircle())
		authorized.PUT("/circle/to-global", s.AddToGlobalCircle())

		// votes
		authorized.POST("/vote", s.CreateVote())

		// rankings
		authorized.GET("/rankings/:circleId", s.Rankings())

		// Upload group
		upload := authorized.Group("/upload")
		upload.PUT("/circle-img", s.UploadCircleImage())
	}
}

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
		authorized.GET("/circles/of-interest", s.CirclesOfInterest())
		authorized.GET("/circles/:name", s.CirclesByName())
		authorized.POST("/circle", s.CreateCircle())
		authorized.PUT("/circle", s.UpdateCircle())
		authorized.DELETE("/circle/:circleId", s.DeleteCircle())
		authorized.PUT("/circle/to-global", s.AddToGlobalCircle())

		// circle voters
		authorized.GET("/circle-voters/:circleId", s.CircleVoters())
		authorized.POST("/circle-voters/:circleId/join", s.CircleVoterJoinCircle())
		authorized.DELETE("/circle-voters/:circleId/leave", s.CircleVoterLeaveCircle())
		authorized.POST("/circle-voters/:circleId/add", s.CircleVoterAddToCircle())

		// circle candidates
		authorized.GET("/circle-candidates/:circleId", s.CircleCandidates())
		authorized.POST("/circle-candidates/:circleId/commitment", s.CircleCandidateCommitment())
		authorized.POST("/circle-candidates/:circleId/join", s.CircleCandidateJoinCircle())
		authorized.DELETE("/circle-candidates/:circleId/leave", s.CircleCandidateLeaveCircle())
		authorized.POST("/circle-candidates/:circleId/add", s.CircleCandidateAddToCircle())
		authorized.GET("/circle-candidates/:circleId/voted-by", s.CircleCandidateVotedBy())

		// votes
		authorized.POST("/vote/:circleId", s.CreateVote())
		authorized.POST("/vote/revoke/:circleId", s.RevokeVote())

		// rankings
		authorized.GET("/rankings/:circleId", s.Rankings())

		// user option
		authorized.GET("/user-option", s.UserOption())

		// ably token
		authorized.GET("/token/ably", s.TokenAbly())

		// Upload group
		upload := authorized.Group("/upload")
		upload.PUT("/circle-img/:circleId", s.UploadCircleImage())
	}
}

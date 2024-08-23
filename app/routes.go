package app

func (s *Server) routes() {
	router := s.router

	// Service group
	v1 := router.Group("/v1/api/vote-circle")

	// Authorization group
	authorized := v1.Group("")
	authorized.Use(s.authGuard(s.authService))
	{
		// circle group
		circle := authorized.Group("/circle")
		circle.GET("/:circleId", s.Circle())
		circle.GET("/:circleId/eligible", s.EligibleToBeInCircle())
		circle.POST("", s.CreateCircle())
		circle.PUT("/:circleId", s.UpdateCircle())
		circle.DELETE("/:circleId", s.DeleteCircle())
		circle.PUT("/to-global", s.AddToGlobalCircle())

		// circles group
		circles := authorized.Group("/circles")
		circles.GET("/of-interest", s.CirclesOfInterest())
		circles.GET("/:name", s.CirclesByName())
		circles.GET("", s.Circles())
		circles.GET("/open-commitments", s.CirclesOpenCommitments())

		// circle voters group
		circleVoters := authorized.Group("/circle-voters")
		circleVoters.GET("/:circleId", s.CircleVoters())
		circleVoters.POST("/:circleId/join", s.CircleVoterJoinCircle())
		circleVoters.DELETE("/:circleId/leave", s.CircleVoterLeaveCircle())
		circleVoters.POST("/:circleId/add", s.CircleVotersAddToCircle())
		circleVoters.POST("/:circleId/remove", s.CircleVoterRemoveFromCircle())

		// circle candidates group
		circleCandidates := authorized.Group("/circle-candidates")
		circleCandidates.GET("/:circleId", s.CircleCandidates())
		circleCandidates.POST("/:circleId/commitment", s.CircleCandidateCommitment())
		circleCandidates.POST("/:circleId/join", s.CircleCandidateJoinCircle())
		circleCandidates.DELETE("/:circleId/leave", s.CircleCandidateLeaveCircle())
		circleCandidates.POST("/:circleId/add", s.CircleCandidatesAddToCircle())
		circleCandidates.POST("/:circleId/remove", s.CircleCandidateRemoveFromCircle())
		circleCandidates.GET("/:circleId/voted-by", s.CircleCandidateVotedBy())

		// vote group
		vote := authorized.Group("/vote")
		vote.POST("/:circleId", s.CreateVote())
		vote.POST("/revoke/:circleId", s.RevokeVote())

		// rankings group
		rankings := authorized.Group("/rankings")
		rankings.GET("/:circleId", s.Rankings())
		rankings.GET("/last-viewed", s.RankingsLastViewed())

		// user option
		authorized.GET("/user-option", s.UserOption())

		// ably token
		authorized.GET("/token/ably", s.TokenAbly())

		// Upload group
		upload := authorized.Group("/upload")
		upload.PUT("/circle-img/:circleId", s.UploadCircleImage())
		upload.DELETE("/circle-img/:circleId", s.DeleteCircleImage())
	}
}

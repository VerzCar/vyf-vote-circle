package app

import (
	"fmt"
	"github.com/VerzCar/vyf-lib-awsx"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/sanitizer"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	router                             *gin.Engine
	authService                        awsx.AuthService
	circleService                      api.CircleService
	circleUploadService                api.CircleUploadService
	rankingService                     api.RankingService
	rankingSubscriptionService         api.RankingSubscriptionService
	voteService                        api.VoteService
	circleVoterService                 api.CircleVoterService
	circleCandidateService             api.CircleCandidateService
	circleCandidateSubscriptionService api.CircleCandidateSubscriptionService
	tokenService                       api.TokenService
	validate                           sanitizer.Validator
	config                             *config.Config
	log                                logger.Logger
}

func NewServer(
	router *gin.Engine,
	authService awsx.AuthService,
	circleService api.CircleService,
	circleUploadService api.CircleUploadService,
	rankingService api.RankingService,
	rankingSubscriptionService api.RankingSubscriptionService,
	voteService api.VoteService,
	circleVoterService api.CircleVoterService,
	circleCandidateService api.CircleCandidateService,
	circleCandidateSubscriptionService api.CircleCandidateSubscriptionService,
	tokenService api.TokenService,
	validate sanitizer.Validator,
	config *config.Config,
	log logger.Logger,
) *Server {
	server := &Server{
		router:                             router,
		authService:                        authService,
		circleService:                      circleService,
		circleUploadService:                circleUploadService,
		rankingService:                     rankingService,
		rankingSubscriptionService:         rankingSubscriptionService,
		voteService:                        voteService,
		circleVoterService:                 circleVoterService,
		circleCandidateService:             circleCandidateService,
		circleCandidateSubscriptionService: circleCandidateSubscriptionService,
		tokenService:                       tokenService,
		validate:                           validate,
		config:                             config,
		log:                                log,
	}

	server.routes()

	return server
}

func (s *Server) Run() error {
	port := fmt.Sprintf(":%s", s.config.Port)
	err := s.router.Run(port)

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}

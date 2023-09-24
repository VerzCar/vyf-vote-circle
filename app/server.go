package app

import (
	"github.com/VerzCar/vyf-lib-awsx"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"log"
)

type Server struct {
	router                     *gin.Engine
	authService                awsx.AuthService
	circleService              api.CircleService
	rankingService             api.RankingService
	rankingSubscriptionService api.RankingSubscriptionService
	voteService                api.VoteService
	validate                   *validator.Validate
	config                     *config.Config
	log                        logger.Logger
}

func NewServer(
	router *gin.Engine,
	authService awsx.AuthService,
	circleService api.CircleService,
	rankingService api.RankingService,
	rankingSubscriptionService api.RankingSubscriptionService,
	voteService api.VoteService,
	validate *validator.Validate,
	config *config.Config,
	log logger.Logger,
) *Server {
	server := &Server{
		router:                     router,
		authService:                authService,
		circleService:              circleService,
		rankingService:             rankingService,
		rankingSubscriptionService: rankingSubscriptionService,
		voteService:                voteService,
		validate:                   validate,
		config:                     config,
		log:                        log,
	}

	server.routes()

	return server
}

func (s *Server) Run() error {
	err := s.router.Run(":8080")

	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v", err)
		return err
	}

	return nil
}

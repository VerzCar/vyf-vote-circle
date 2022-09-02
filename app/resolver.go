package app

//go:generate go run github.com/99designs/gqlgen
import (
	"github.com/go-playground/validator/v10"
	"gitlab.vecomentman.com/libs/awsx"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
)

type Resolver struct {
	authService                awsx.AuthService
	circleService              api.CircleService
	rankingService             api.RankingService
	rankingSubscriptionService api.RankingSubscriptionService
	voteService                api.VoteService
	validate                   *validator.Validate
	config                     *config.Config
	log                        logger.Logger
}

func NewResolver(
	authService awsx.AuthService,
	circleService api.CircleService,
	rankingService api.RankingService,
	rankingSubscriptionService api.RankingSubscriptionService,
	voteService api.VoteService,
	validate *validator.Validate,
	config *config.Config,
	log logger.Logger,
) *Resolver {
	return &Resolver{
		authService:                authService,
		circleService:              circleService,
		rankingService:             rankingService,
		rankingSubscriptionService: rankingSubscriptionService,
		voteService:                voteService,
		validate:                   validate,
		config:                     config,
		log:                        log,
	}
}

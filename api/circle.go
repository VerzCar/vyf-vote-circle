package api

import (
	"context"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	routerContext "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router/ctx"
)

type CircleService interface {
	Circle(
		ctx context.Context,
		circleId int64,
	) (*model.Circle, error)
}

type CircleRepository interface {
	CircleById(id int64) (*model.Circle, error)
}

type circleService struct {
	storage CircleRepository
	config  *config.Config
	log     logger.Logger
}

func NewCircleService(
	circleRepo CircleRepository,
	config *config.Config,
	log logger.Logger,
) CircleService {
	return &circleService{
		storage: circleRepo,
		config:  config,
		log:     log,
	}
}

func (c *circleService) Circle(
	ctx context.Context,
	circleId int64,
) (*model.Circle, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	c.log.Infof(authClaims.Subject)

	circle, err := c.storage.CircleById(circleId)

	c.log.Info(circle)
	return circle, nil
}

package api

import (
	"context"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	routerContext "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router/ctx"
)

type RankingService interface {
	Rankings(
		ctx context.Context,
		circleId int64,
	) ([]*model.Ranking, error)
}

type RankingRepository interface {
	RankingsByCircleId(circleId int64) ([]*model.Ranking, error)
}

type rankingService struct {
	storage RankingRepository
	config  *config.Config
	log     logger.Logger
}

func NewRankingService(
	circleRepo RankingRepository,
	config *config.Config,
	log logger.Logger,
) RankingService {
	return &rankingService{
		storage: circleRepo,
		config:  config,
		log:     log,
	}
}

func (c *rankingService) Rankings(
	ctx context.Context,
	circleId int64,
) ([]*model.Ranking, error) {
	_, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	rankings, err := c.storage.RankingsByCircleId(circleId)

	if err != nil {
		return nil, err
	}

	return rankings, nil
}

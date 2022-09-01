package api

import (
	"context"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database"
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
	Votes(circleId int64) ([]*model.Vote, error)
}

type RankingCache interface {
	RankingList(
		ctx context.Context,
		circleId int64,
	) ([]*model.Ranking, error)
	ExistsRankingListForCircle(
		ctx context.Context,
		circleId int64,
	) (bool, error)
	BuildRankingList(
		ctx context.Context,
		circleId int64,
		votes []*model.Vote,
	) error
}

type rankingService struct {
	storage RankingRepository
	cache   RankingCache
	config  *config.Config
	log     logger.Logger
}

func NewRankingService(
	circleRepo RankingRepository,
	cache RankingCache,
	config *config.Config,
	log logger.Logger,
) RankingService {
	return &rankingService{
		storage: circleRepo,
		cache:   cache,
		config:  config,
		log:     log,
	}
}

// Rankings from the circle with the given circle id.
// It returns always the ranking list from the cache.
// It first checks whether some votes already exists for this circle in the cache
// otherwise it will build up the cache with the votes for this circle.
// If the circle has"nt any votes, an empty list will be returned.
func (c *rankingService) Rankings(
	ctx context.Context,
	circleId int64,
) ([]*model.Ranking, error) {
	_, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	exists, err := c.cache.ExistsRankingListForCircle(ctx, circleId)

	if err != nil {
		c.log.Errorf("error for check if ranking list exists for circle with id %d: %s", circleId, err)
		return nil, err
	}

	if !exists {
		isEmpty, err := c.buildRankingList(ctx, circleId)

		if err != nil {
			return nil, err
		}

		if isEmpty {
			return make([]*model.Ranking, 0), nil
		}
	}

	rankings, err := c.cache.RankingList(ctx, circleId)

	if err != nil {
		return nil, err
	}

	return rankings, nil
}

// buildRankingList for the given circle.
// Returns true if the circle does not contain any votes
// (has an empty ranking list), otherwise false or an error if any occurs.
func (c *rankingService) buildRankingList(
	ctx context.Context,
	circleId int64,
) (bool, error) {
	votes, err := c.storage.Votes(circleId)

	switch {
	case err != nil && !database.RecordNotFound(err):
		{
			c.log.Errorf("error building up cache for circle id %d: %s", circleId, err)
			return false, err
		}
	case database.RecordNotFound(err):
		{
			return true, nil
		}
	default:
		{
			err := c.cache.BuildRankingList(ctx, circleId, votes)
			return false, err
		}
	}
}

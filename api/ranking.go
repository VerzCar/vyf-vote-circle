package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
)

type RankingService interface {
	Rankings(
		ctx context.Context,
		circleId int64,
	) ([]*model.RankingResponse, error)
	LastViewedRankings(
		ctx context.Context,
	) ([]*model.RankingLastViewedResponse, error)
}

type RankingRepository interface {
	CircleById(id int64) (*model.Circle, error)
	RankingsByCircleId(circleId int64) ([]*model.Ranking, error)
	Votes(circleId int64) ([]*model.Vote, error)
	CountsVotesOfCandidateByCircleId(circleId int64, candidateId int64) (int64, error)
	RankingByCircleId(circleId int64, identityId string) (*model.Ranking, error)
	CreateNewRankingLastViewed(
		circleId int64,
		identityId string,
	) (*model.RankingLastViewed, error)
	RankingsLastViewedByUserIdentityId(identityId string) ([]*model.RankingLastViewed, error)
}

type RankingCache interface {
	RankingList(
		ctx context.Context,
		circleId int64,
		fromRanking *model.RankingResponse,
	) ([]*model.RankingResponse, error)
	ExistsRankingListForCircle(
		ctx context.Context,
		circleId int64,
	) (bool, error)
	BuildRankingList(
		ctx context.Context,
		circleId int64,
		rankingCacheItems []*model.RankingCacheItem,
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
// If the circle hasn't any votes, an empty list will be returned.
func (c *rankingService) Rankings(
	ctx context.Context,
	circleId int64,
) ([]*model.RankingResponse, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return nil, err
	}

	if !circle.IsEditable() {
		rankings, err := c.storage.RankingsByCircleId(circleId)

		if err != nil && !database.RecordNotFound(err) {
			return nil, err
		}

		if database.RecordNotFound(err) {
			return make([]*model.RankingResponse, 0), nil
		}

		return c.mapRankingToRankingResponse(rankings), nil
	}

	exists, err := c.cache.ExistsRankingListForCircle(ctx, circleId)

	if err != nil {
		c.log.Errorf("error for check if ranking list exists for circle with id %d: %s", circleId, err)
		return nil, err
	}

	if !exists {
		isEmpty, err := c.buildCacheRankingList(ctx, circleId)

		if err != nil {
			return nil, err
		}

		if isEmpty {
			return make([]*model.RankingResponse, 0), nil
		}
	}

	rankings, err := c.cache.RankingList(ctx, circleId, nil)

	if err != nil {
		return nil, err
	}

	_, _ = c.storage.CreateNewRankingLastViewed(circleId, authClaims.Subject)

	return rankings, nil
}

func (c *rankingService) LastViewedRankings(
	ctx context.Context,
) ([]*model.RankingLastViewedResponse, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	lastViewedRankings, err := c.storage.RankingsLastViewedByUserIdentityId(authClaims.Subject)

	if err != nil {
		return nil, err
	}
	fmt.Println()
	fmt.Println(lastViewedRankings)
	fmt.Println()

	lastViewedRankingsResponse := make([]*model.RankingLastViewedResponse, 0)

	for _, ranking := range lastViewedRankings {
		circlePaginated := &model.CirclePaginatedResponse{
			Name:        ranking.Circle.Name,
			Description: ranking.Circle.Description,
			ImageSrc:    ranking.Circle.ImageSrc,
			Stage:       ranking.Circle.Stage,
			ID:          ranking.Circle.ID,
			Active:      ranking.Circle.Active,
		}

		res := &model.RankingLastViewedResponse{
			CreatedAt: ranking.CreatedAt,
			UpdatedAt: ranking.UpdatedAt,
			Circle:    circlePaginated,
			ID:        ranking.ID,
		}

		lastViewedRankingsResponse = append(lastViewedRankingsResponse, res)
	}

	return lastViewedRankingsResponse, nil
}

// buildCacheRankingList for the given circle.
// Returns true if the circle does not contain any votes
// (has an empty ranking list), otherwise false or an error if any occurs.
func (c *rankingService) buildCacheRankingList(
	ctx context.Context,
	circleId int64,
) (bool, error) {
	votes, err := c.storage.Votes(circleId)

	switch {
	case err != nil && !database.RecordNotFound(err):
		{
			c.log.Errorf("error building up ranking list for circle id %d: %s", circleId, err)
			return false, err
		}
	case database.RecordNotFound(err):
		{
			return true, nil
		}
	default:
		{
			var rankingCacheItems []*model.RankingCacheItem
			for _, vote := range votes {
				voteCount, err := c.storage.CountsVotesOfCandidateByCircleId(circleId, vote.Candidate.ID)

				if err != nil {
					c.log.Errorf("error getting vote count for candidate id %d: %s", vote.Candidate.ID, err)
					return false, err
				}

				ranking, err := c.storage.RankingByCircleId(circleId, vote.Candidate.Candidate)

				if err != nil {
					c.log.Errorf("error reading by circle id %d ranking: %s", circleId, err)
					return false, err
				}

				rankingCacheItem := &model.RankingCacheItem{
					Candidate: vote.Candidate,
					Ranking:   ranking,
					VoteCount: voteCount,
				}
				rankingCacheItems = append(rankingCacheItems, rankingCacheItem)
			}

			err := c.cache.BuildRankingList(ctx, circleId, rankingCacheItems)

			return false, err
		}
	}
}

func (c *rankingService) mapRankingToRankingResponse(rankings []*model.Ranking) []*model.RankingResponse {
	var responses []*model.RankingResponse
	for _, ranking := range rankings {
		response := &model.RankingResponse{
			CreatedAt:    ranking.CreatedAt,
			UpdatedAt:    ranking.UpdatedAt,
			IdentityID:   ranking.IdentityID,
			Placement:    ranking.Placement,
			ID:           ranking.ID,
			CandidateID:  0,
			Number:       ranking.Number,
			Votes:        ranking.Votes,
			IndexedOrder: 0,
			CircleID:     ranking.CircleID,
		}
		responses = append(responses, response)
	}

	return responses
}

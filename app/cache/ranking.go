package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/go-redis/redis/v8"
	"time"
)

func (c *redisCache) UpsertRanking(
	ctx context.Context,
	circleId int64,
	candidate *model.CircleCandidate,
	ranking *model.Ranking,
	votes int64,
) (*model.RankingResponse, error) {
	rankingScore := &model.RankingScore{
		VoteCount:      votes,
		UserIdentityId: candidate.Candidate,
	}

	if err := c.setRankingScore(ctx, circleId, candidate, ranking, rankingScore); err != nil {
		return nil, err
	}

	placementNumber, rankingPlacementIndex, err := c.rankingPlacementNumberWithIndex(
		ctx,
		circleId,
		rankingScore.UserIdentityId,
	)

	if err != nil {
		return nil, err
	}

	rankingRes := populateRanking(
		ranking.ID,
		circleId,
		candidate.ID,
		rankingScore,
		rankingPlacementIndex,
		placementNumber,
		ranking.CreatedAt,
		ranking.UpdatedAt,
	)

	return rankingRes, nil
}

func (c *redisCache) RemoveRanking(
	ctx context.Context,
	circleId int64,
	candidate *model.CircleCandidate,
) error {
	if err := c.removeRanking(ctx, circleId, candidate); err != nil {
		return err
	}

	return nil
}

// RankingList of the current cached ranking for the circle
func (c *redisCache) RankingList(
	ctx context.Context,
	circleId int64,
) ([]*model.RankingResponse, error) {
	circleRankingKey := circleRankingKey(circleId)

	rankingScores, err := c.rankingScores(ctx, circleRankingKey)

	if err != nil {
		c.log.Errorf(
			"error getting ranking scores: for circle key %s: %s",
			circleRankingKey,
			err,
		)
		return nil, err
	}

	rankingList := make([]*model.RankingResponse, 0)
	placementNumber := int64(0)
	var voteCount int64

	for placementIndex, rankingScore := range rankingScores {
		rankingUserCandidateKey := circleUserCandidateKey(circleId, rankingScore.UserIdentityId)
		rankingUserCandidate, err := c.rankingUserCandidate(ctx, rankingUserCandidateKey)

		if err != nil {
			c.log.Errorf(
				"error getting ranking user %s candidate: for circle key %s: %s",
				rankingScore.UserIdentityId,
				circleRankingKey,
				err,
			)
			return nil, err
		}

		if voteCount != rankingScore.VoteCount {
			voteCount = rankingScore.VoteCount
			placementNumber++
		}

		rankingList = append(
			rankingList,
			populateRanking(
				rankingUserCandidate.RankingID,
				circleId,
				rankingUserCandidate.CandidateID,
				rankingScore,
				int64(placementIndex),
				placementNumber,
				rankingUserCandidate.CreatedAt,
				rankingUserCandidate.UpdatedAt,
			),
		)
	}

	return rankingList, nil
}

// ExistsRankingListForCircle with given circle id checks whether a
// ranking list for this circle is in cache.
// Returns true if exists in cache, otherwise false.
func (c *redisCache) ExistsRankingListForCircle(
	ctx context.Context,
	circleId int64,
) (bool, error) {
	circleRankingKey := circleRankingKey(circleId)

	result := c.redis.Exists(ctx, circleRankingKey)

	if result.Err() != nil {
		return false, result.Err()
	}

	return result.Val() > 0, nil
}

// BuildRankingList from votes for the circle id.
func (c *redisCache) BuildRankingList(
	ctx context.Context,
	circleId int64,
	rankingCacheItems []*model.RankingCacheItem,
) error {
	for _, item := range rankingCacheItems {
		_, err := c.UpsertRanking(ctx, circleId, item.Candidate, item.Ranking, item.VoteCount)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *redisCache) setRankingScore(
	ctx context.Context,
	circleId int64,
	candidate *model.CircleCandidate,
	ranking *model.Ranking,
	rankingScore *model.RankingScore,
) error {
	key := circleRankingKey(circleId)

	_, err := c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			pipeSetRankingScore(ctx, pipe, key, rankingScore)
			candidateKey := circleUserCandidateKey(circleId, candidate.Candidate)
			pipeSetUserCandidate(ctx, pipe, candidateKey, candidate, ranking)
			return nil
		},
	)

	if err != nil {
		c.log.Errorf(
			"could not set ranking in transaction for circle key %s: %s",
			key,
			err,
		)
		return err
	}

	return nil
}

func (c *redisCache) rankingPlacementNumberWithIndex(
	ctx context.Context,
	circleId int64,
	member string,
) (int64, int64, error) {
	key := circleRankingKey(circleId)

	cmds, err := c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			pipe.ZScore(ctx, key, member)
			pipe.ZRevRank(ctx, key, member)
			return nil
		},
	)

	if err != nil {
		c.log.Errorf(
			"could not read placement number and index for member %s for circle key %s: %s",
			member,
			key,
			err,
		)
		return 0, 0, err
	}

	placementNumber := int64(cmds[0].(*redis.FloatCmd).Val())
	placementIndex := cmds[1].(*redis.IntCmd).Val()

	return placementNumber, placementIndex, nil
}

func (c *redisCache) removeRanking(
	ctx context.Context,
	circleId int64,
	candidate *model.CircleCandidate,
) error {
	key := circleRankingKey(circleId)

	_, err := c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			pipeRemoveRankingScore(ctx, pipe, key, candidate.Candidate)
			candidateKey := circleUserCandidateKey(circleId, candidate.Candidate)
			pipeRemoveUserCandidate(ctx, pipe, candidateKey)
			return nil
		},
	)

	if err != nil {
		c.log.Errorf(
			"could not remove ranking in transaction for circle key %s: %s",
			key,
			err,
		)
		return err
	}

	return nil
}

// Ranking scores of the given key as a list
func (c *redisCache) rankingScores(
	ctx context.Context,
	key string,
) ([]*model.RankingScore, error) {
	result := c.redis.ZRevRangeWithScores(ctx, key, 0, -1)

	switch {
	case errors.Is(result.Err(), redis.Nil):
		return nil, nil
	case result.Err() != nil:
		return nil, result.Err()
	default:
		var rankingScores []*model.RankingScore
		for _, z := range result.Val() {
			rankingScore := &model.RankingScore{
				VoteCount:      int64(z.Score),
				UserIdentityId: z.Member.(string),
			}
			rankingScores = append(rankingScores, rankingScore)
		}
		return rankingScores, nil
	}
}

// Ranking score of the given key and member.
// Returns 0 if the key or member does not exist
func (c *redisCache) rankingScore(
	ctx context.Context,
	key string,
	member string,
) (int64, error) {
	result := c.redis.ZScore(ctx, key, member)

	switch {
	case errors.Is(result.Err(), redis.Nil):
		return 0, nil
	case result.Err() != nil:
		return 0, result.Err()
	default:
		return int64(result.Val()), nil
	}
}

// Identifies the current index (placement)
// of the given score that the member is.
// Returns 0 if the key or member does not exist.
func (c *redisCache) rankingPlacementIndex(
	ctx context.Context,
	key string,
	member string,
) (int64, error) {
	result := c.redis.ZRevRank(ctx, key, member)

	switch {
	case errors.Is(result.Err(), redis.Nil):
		return 0, nil
	case result.Err() != nil:
		return 0, result.Err()
	default:
		return result.Val(), nil
	}
}

func (c *redisCache) rankingUserCandidate(
	ctx context.Context,
	key string,
) (*model.RankingUserCandidate, error) {
	var userCandidate model.RankingUserCandidate

	err := c.redis.HGetAll(ctx, key).Scan(&userCandidate)

	if err != nil {
		return nil, err
	}

	return &userCandidate, nil
}

// Sets the ranking score for the given key in cache
func pipeSetRankingScore(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	rankingScore *model.RankingScore,
) {
	members := &redis.Z{
		Score:  float64(rankingScore.VoteCount),
		Member: rankingScore.UserIdentityId,
	}
	pipe.ZAdd(ctx, key, members)
}

// Removes the member for the given key in cache
func pipeRemoveRankingScore(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	member string,
) {
	pipe.ZRem(ctx, key, member)
}

// Ranking score of the given key and member.
// Returns 0 if the key or member does not exist
func pipeRankingScore(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	member string,
) *redis.FloatCmd {
	return pipe.ZScore(ctx, key, member)
}

func pipeSetUserCandidate(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	candidate *model.CircleCandidate,
	ranking *model.Ranking,
) {
	pipe.HSet(ctx, key, "candidateId", candidate.ID)
	pipe.HSet(ctx, key, "rankingId", ranking.ID)
	pipe.HSet(ctx, key, "createdAt", ranking.CreatedAt)
	pipe.HSet(ctx, key, "updatedAt", ranking.UpdatedAt)
}

func pipeRemoveUserCandidate(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
) {
	pipe.HDel(ctx, key, "candidateId", "rankingId", "createdAt", "updatedAt")
}

func circleRankingKey(circleId int64) string {
	return fmt.Sprintf("cirlce:%d:ranking", circleId)
}

func circleUserCandidateKey(circleId int64, identityId string) string {
	return fmt.Sprintf("circle:%d:%s", circleId, identityId)
}

func populateRanking(
	id int64,
	circleId int64,
	candidateId int64,
	rankingScore *model.RankingScore,
	placementIndex int64,
	placementNumber int64,
	createdAt time.Time,
	updatedAt time.Time,
) *model.RankingResponse {
	return &model.RankingResponse{
		ID:           id,
		CandidateID:  candidateId,
		IdentityID:   rankingScore.UserIdentityId,
		Number:       placementNumber,
		Votes:        rankingScore.VoteCount,
		IndexedOrder: placementIndex,
		Placement:    model.PlacementNeutral,
		CircleID:     circleId,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}
}

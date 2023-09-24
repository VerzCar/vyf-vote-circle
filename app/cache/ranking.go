package cache

import (
	"context"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/go-redis/redis/v8"
	"strconv"
	"time"
)

func (c *redisCache) UpdateRanking(
	ctx context.Context,
	circleId int64,
	identityId string,
	votes int64,
) error {
	rankingScore := &model.RankingScore{
		VoteCount:      votes,
		UserIdentityId: identityId,
	}

	circleRankingKey := circleRankingKey(circleId)

	err := c.setRankingScore(ctx, circleRankingKey, rankingScore)

	if err != nil {
		c.log.Errorf(
			"error setting ranking score: for circle key %s with ranking score %v: %s",
			circleRankingKey,
			rankingScore,
			err,
		)
		return err
	}

	return nil
}

// RankingList of the current cached ranking for the circle
func (c *redisCache) RankingList(
	ctx context.Context,
	circleId int64,
) ([]*model.Ranking, error) {
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

	var rankingList []*model.Ranking
	placementNumber := int64(0)
	var voteCount int64

	for index, rankingScore := range rankingScores {
		if voteCount != rankingScore.VoteCount {
			voteCount = rankingScore.VoteCount
			placementNumber++
		}
		rankingList = append(
			rankingList,
			populateRanking(
				circleId,
				int64(index)+1,
				rankingScore,
				placementNumber,
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
	votes []*model.Vote,
) error {
	circleRankingKey := circleRankingKey(circleId)

	for _, vote := range votes {
		identityId := vote.Elected.Voter
		score, err := c.rankingScore(ctx, circleRankingKey, identityId)

		if err != nil {
			c.log.Errorf(
				"error getting ranking score for voter %s: for circle key %s: %s",
				identityId,
				circleRankingKey,
				err,
			)
			return err
		}

		err = c.UpdateRanking(ctx, circleId, identityId, score+1)

		if err != nil {
			return err
		}
	}

	return nil
}

// setRankingScore sets the ranking score for the given key in cache
func (c *redisCache) setRankingScore(ctx context.Context, key string, rankingScore *model.RankingScore) error {
	members := &redis.Z{
		Score:  float64(rankingScore.VoteCount),
		Member: rankingScore.UserIdentityId,
	}

	return c.redis.ZAdd(ctx, key, members).Err()
}

// rankingScores of the given key as a list
func (c *redisCache) rankingScores(
	ctx context.Context,
	key string,
) ([]*model.RankingScore, error) {
	result := c.redis.ZRevRangeWithScores(ctx, key, 0, -1)

	switch {
	case result.Err() == redis.Nil:
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

// rankingScore of the given key and member.
// Returns 0 if the key or member does not exist
func (c *redisCache) rankingScore(
	ctx context.Context,
	key string,
	member string,
) (int64, error) {
	result := c.redis.ZScore(ctx, key, member)

	switch {
	case result.Err() == redis.Nil:
		return 0, nil
	case result.Err() != nil:
		return 0, result.Err()
	default:
		return int64(result.Val()), nil
	}
}

func circleRankingKey(id int64) string {
	return "circle:" + strconv.FormatInt(id, 10) + ":ranking"
}

func populateRanking(
	circleId int64,
	rankingId int64,
	rankingScore *model.RankingScore,
	placementNumber int64,
) *model.Ranking {
	return &model.Ranking{
		ID:         rankingId,
		IdentityID: rankingScore.UserIdentityId,
		Number:     placementNumber,
		Votes:      rankingScore.VoteCount,
		Placement:  model.PlacementNeutral,
		CircleID:   circleId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"strconv"
)

func (c *redisCache) UpdateRanking(
	ctx context.Context,
	circleId int64,
	identityId model.UserIdentityId,
	votes int64,
) error {
	rankingScore := &model.RankingScore{
		VoteCount:      model.VoteCount(votes),
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

func circleRankingKey(id int64) string {
	return "circle:" + strconv.FormatInt(id, 10) + ":ranking"
}

// setRankingScore sets the ranking score for the given key in cache
func (c *redisCache) setRankingScore(ctx context.Context, key string, rankingScore *model.RankingScore) error {
	members := &redis.Z{
		Score:  float64(rankingScore.VoteCount),
		Member: rankingScore.UserIdentityId,
	}

	return c.redis.ZAdd(ctx, key, members).Err()
}

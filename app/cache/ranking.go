package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
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

	for _, rankingScore := range rankingScores {
		if voteCount != rankingScore.VoteCount {
			voteCount = rankingScore.VoteCount
			placementNumber++
		}
		rankingList = append(rankingList, populateRanking(circleId, rankingScore, placementNumber))
	}

	return rankingList, nil
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

func populateRanking(
	circleId int64,
	rankingScore *model.RankingScore,
	placementNumber int64,
) *model.Ranking {
	return &model.Ranking{
		IdentityID: rankingScore.UserIdentityId,
		Number:     placementNumber,
		Votes:      rankingScore.VoteCount,
		Placement:  "",
		CircleID:   circleId,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

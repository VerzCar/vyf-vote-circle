package cache

import (
	"context"
	"errors"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/go-redis/redis/v8"
	"strconv"
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

	key := circleRankingKey(circleId)

	var rankingPlacementIndex int64

	_, err := c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			err := pipeSetRankingScore(ctx, pipe, key, rankingScore)

			if err != nil {
				c.log.Errorf(
					"error setting ranking score: for circle key %s with ranking score %v: %s",
					key,
					rankingScore,
					err,
				)
				return err
			}

			score, err := pipeRankingScore(ctx, pipe, key, rankingScore.UserIdentityId)

			if err != nil {
				c.log.Errorf(
					"error getting ranking score for candidate %s: for circle key %s: %s",
					rankingScore.UserIdentityId,
					key,
					err,
				)
				return err
			}

			rankingPlacementIndex, err = pipeRankingPlacementIndex(
				ctx,
				pipe,
				key,
				strconv.FormatInt(score, 10),
				"",
			)

			if err != nil {
				c.log.Errorf(
					"error getting ranking score for candidate %s: for circle key %s: %s",
					rankingScore.UserIdentityId,
					key,
					err,
				)
				return err
			}

			candidateKey := circleUserCandidateKey(circleId, candidate.Candidate)
			err = pipeSetUserCandidate(ctx, pipe, candidateKey, candidate, ranking)

			if err != nil {
				c.log.Errorf(
					"error setting user candidate key candidate id %d: for circle id %s: %s",
					candidate.ID,
					candidate.Candidate,
					err,
				)
				return err
			}

			return nil
		},
	)

	if err != nil {
		c.log.Errorf(
			"could not update ranking in transaction: %s",
			err,
		)
		return nil, err
	}

	rankingRes := populateRanking(
		ranking.ID,
		circleId,
		candidate.ID,
		rankingScore,
		rankingPlacementIndex+1,
		ranking.CreatedAt,
		ranking.UpdatedAt,
	)

	return rankingRes, nil
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

	for _, rankingScore := range rankingScores {
		if voteCount != rankingScore.VoteCount {
			voteCount = rankingScore.VoteCount
			placementNumber++
		}

		rankingList = append(
			rankingList,
			populateRanking(
				circleId,
				0,
				0,
				rankingScore,
				placementNumber,
				time.Now(),
				time.Now(),
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
		identityId := vote.Candidate.Candidate
		_, err := c.rankingScore(ctx, circleRankingKey, identityId)

		if err != nil {
			c.log.Errorf(
				"error getting ranking score for voter %s: for circle key %s: %s",
				identityId,
				circleRankingKey,
				err,
			)
			return err
		}

		//_, err = c.UpsertRanking(ctx, circleId, identityId, 0, score+1)
		//
		//if err != nil {
		//	return err
		//}
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

// Sets the ranking score for the given key in cache
func pipeSetRankingScore(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	rankingScore *model.RankingScore,
) error {
	members := &redis.Z{
		Score:  float64(rankingScore.VoteCount),
		Member: rankingScore.UserIdentityId,
	}

	return pipe.ZAdd(ctx, key, members).Err()
}

// Ranking score of the given key and member.
// Returns 0 if the key or member does not exist
func pipeRankingScore(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	member string,
) (int64, error) {
	result := pipe.ZScore(ctx, key, member)

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
// Returns 0 if the key or member does not exist
func pipeRankingPlacementIndex(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	min string,
	max string,
) (int64, error) {
	result := pipe.ZCount(ctx, key, min, max)

	switch {
	case errors.Is(result.Err(), redis.Nil):
		return 0, nil
	case result.Err() != nil:
		return 0, result.Err()
	default:
		return result.Val(), nil
	}
}

func pipeSetUserCandidate(
	ctx context.Context,
	pipe redis.Pipeliner,
	key string,
	candidate *model.CircleCandidate,
	ranking *model.Ranking,
) error {
	err := pipe.HSet(ctx, key, "candidateId", candidate.ID).Err()

	if err != nil {
		return err
	}

	err = pipe.HSet(ctx, key, "rankingId", ranking.ID).Err()

	if err != nil {
		return err
	}

	err = pipe.HSet(ctx, key, "createdAt", ranking.CreatedAt).Err()

	if err != nil {
		return err
	}

	return pipe.HSet(ctx, key, "updatedAt", ranking.UpdatedAt).Err()
}

func circleRankingKey(id int64) string {
	return "circle:" + strconv.FormatInt(id, 10) + ":ranking"
}

func circleUserCandidateKey(circleId int64, identityId string) string {
	return "circle:" + strconv.FormatInt(circleId, 10) + ":" + identityId
}

func populateRanking(
	id int64,
	circleId int64,
	candidateId int64,
	rankingScore *model.RankingScore,
	placementNumber int64,
	createdAt time.Time,
	updatedAt time.Time,
) *model.RankingResponse {
	return &model.RankingResponse{
		ID:          id,
		CandidateID: candidateId,
		IdentityID:  rankingScore.UserIdentityId,
		Number:      placementNumber,
		Votes:       rankingScore.VoteCount,
		Placement:   model.PlacementNeutral,
		CircleID:    circleId,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
}

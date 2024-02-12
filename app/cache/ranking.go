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
	// TODO: put this in transaction
	rankingScore, err := c.setRankingScore(ctx, circleId, candidate, ranking, votes)

	if err != nil {
		return nil, err
	}

	rankingPlacementIndex, highestVotedMember, err := c.rankingIndexWithLastestVoteCountMember(
		ctx,
		circleId,
		rankingScore.UserIdentityId,
		votes,
	)

	if err != nil {
		return nil, err
	}

	key := circleRankingKey(circleId)
	highestVotedMemberIndex, err := c.rankingPlacementIndex(ctx, key, highestVotedMember)

	if err != nil {
		c.log.Errorf(
			"could not read index for member %s for circle key %s: %s",
			highestVotedMember,
			key,
			err,
		)
		return nil, err
	}

	placementNumber := highestVotedMemberIndex + 1

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
	// TODO: put this in transaction
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
	key := circleRankingKey(circleId)

	rankingScores, err := c.rankingScores(ctx, key)

	if err != nil {
		c.log.Errorf(
			"error getting ranking scores: for circle key %s: %s",
			key,
			err,
		)
		return nil, err
	}

	rankingList, err := c.rankingList(ctx, circleId, rankingScores)

	if err != nil {
		return nil, err
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
	key := circleRankingKey(circleId)

	result := c.redis.Exists(ctx, key)

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
		_, err := c.setRankingScore(ctx, circleId, item.Candidate, item.Ranking, item.VoteCount)

		if err != nil {
			return err
		}
	}

	// TODO: set expiration based on circle inactive time
	expirationDuration := time.Duration(72) * time.Hour
	_ = c.setExpiration(ctx, circleRankingKey(circleId), expirationDuration)

	return nil
}

func (c *redisCache) setRankingScore(
	ctx context.Context,
	circleId int64,
	candidate *model.CircleCandidate,
	ranking *model.Ranking,
	votes int64,
) (*model.RankingScore, error) {
	key := circleRankingKey(circleId)
	rankingScore := &model.RankingScore{
		VoteCount:      votes,
		UserIdentityId: candidate.Candidate,
	}

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
		return nil, err
	}

	return rankingScore, nil
}

func (c *redisCache) rankingList(
	ctx context.Context,
	circleId int64,
	rankingScores []*model.RankingScore,
) ([]*model.RankingResponse, error) {
	key := circleRankingKey(circleId)

	cmds, err := c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			for _, rankingScore := range rankingScores {
				rankingUserCandidateKey := circleUserCandidateKey(circleId, rankingScore.UserIdentityId)
				pipe.HGetAll(ctx, rankingUserCandidateKey)
			}
			return nil
		},
	)

	if err != nil {
		c.log.Errorf(
			"could not get user canidates of ranking for circle key %s: %s",
			key,
			err,
		)
		return nil, err
	}

	rankingList := make([]*model.RankingResponse, 0)
	placementNumber := int64(0)
	var voteCount int64

	for placementIndex, rankingScore := range rankingScores {
		var rankingUserCandidate model.RankingUserCandidate

		err = cmds[placementIndex].(*redis.StringStringMapCmd).Scan(&rankingUserCandidate)

		if err != nil {
			c.log.Errorf(
				"error getting ranking user %s candidate: for circle key %s: %s",
				rankingScore.UserIdentityId,
				key,
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

func (c *redisCache) rankingIndexWithLastestVoteCountMember(
	ctx context.Context,
	circleId int64,
	member string,
	voteCount int64,
) (int64, string, error) {
	key := circleRankingKey(circleId)

	cmds, err := c.redis.Pipelined(
		ctx, func(pipe redis.Pipeliner) error {
			pipe.ZRevRank(ctx, key, member)
			rangeArgs := redis.ZRangeArgs{
				Key:     key,
				Start:   voteCount,
				Stop:    voteCount,
				ByScore: true,
				ByLex:   false,
				Rev:     true,
				Offset:  0,
				Count:   1,
			}
			pipe.ZRangeArgs(ctx, rangeArgs)
			return nil
		},
	)

	if err != nil {
		c.log.Errorf(
			"could not read ranking index for member %s for circle key %s: %s",
			member,
			key,
			err,
		)
		return 0, "", err
	}

	placementIndex := cmds[0].(*redis.IntCmd).Val()
	highestMember := cmds[1].(*redis.StringSliceCmd).Val()[0]

	return placementIndex, highestMember, nil
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

func (c *redisCache) setExpiration(
	ctx context.Context,
	key string,
	expiration time.Duration,
) error {
	err := c.redis.Expire(ctx, key, expiration).Err()

	if err != nil {
		return err
	}

	return nil
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
	return fmt.Sprintf("circle:%d:ranking", circleId)
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

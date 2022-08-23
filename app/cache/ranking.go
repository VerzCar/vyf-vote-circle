package cache

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"strconv"
)

func (c *redisCache) UpdateRanking(
	ctx context.Context,
	circleId int64,
	identityID string,
) error {

	rankingMap, err := c.getRanking(ctx, circleId, identityID)

	if err != nil {
		return err
	}
	rankingMap[identityID].Votes = rankingMap[identityID].Votes + 1

	index, err := c.findVotesRankingPlacement(ctx, circleId, rankingMap[identityID].Votes)

	if err != nil {
		return err
	}

	if index == -1 {
		index, err = c.pushToVotesRankingPlacement(ctx, circleId, rankingMap[identityID].Votes)
		if err != nil {
			return err
		}
	}

	rankingMap[identityID].Number = int(index + 1)

	_, err = c.setRanking(ctx, circleId, rankingMap)

	if err != nil {
		return err
	}

	return nil
}

func (c *redisCache) getRanking(
	ctx context.Context,
	circleId int64,
	identityID string,
) (model.RankingMap, error) {
	hashField := hashField(circleId)
	rankingMap := model.RankingMap{
		identityID: &model.Ranking{
			IdentityID: identityID,
			CircleID:   circleId,
		},
	}
	_, err := c.getHashJson(ctx, hashField, identityID, rankingMap)

	if err != nil {
		c.log.Errorf(
			"error reading ranking: reading field %s with current key %s: %s",
			hashField,
			identityID,
			err,
		)
		return nil, err
	}

	return rankingMap, nil
}

func (c *redisCache) setRanking(
	ctx context.Context,
	circleId int64,
	rankingMap model.RankingMap,
) (model.RankingMap, error) {
	hashField := hashField(circleId)

	err := c.setHashJson(ctx, hashField, rankingMap)

	if err != nil {
		c.log.Errorf(
			"error reading ranking: reading field %s with ranking map %v: %s",
			hashField,
			rankingMap,
			err,
		)
		return nil, err
	}

	return rankingMap, nil
}

func (c *redisCache) findVotesRankingPlacement(
	ctx context.Context,
	circleId int64,
	voteCount int,
) (int64, error) {
	entry, err := c.getIndexInList(ctx, circleVoteRankingListKey(circleId), string(rune(voteCount)))

	if err != nil {
		c.log.Errorf(
			"error reading key %s in list for vote count %d: %s",
			circleVoteRankingListKey(circleId),
			voteCount,
			err,
		)
		return -1, err
	}

	if !entry.Exists {
		return -1, fmt.Errorf("company cannot be verified anymore. new token required")
	}

	return entry.Val, nil
}

func (c *redisCache) pushToVotesRankingPlacement(
	ctx context.Context,
	circleId int64,
	voteCount int,
) (int64, error) {
	index, err := c.pushToListEnd(ctx, circleVoteRankingListKey(circleId), voteCount)

	if err != nil {
		c.log.Errorf(
			"error push element to list with key %s for vote count %d: %s",
			circleVoteRankingListKey(circleId),
			voteCount,
			err,
		)
		return 0, err
	}

	return index, nil
}

func hashField(id int64) string {
	return "circle_" + strconv.FormatInt(id, 10)
}

func circleVoteRankingListKey(id int64) string {
	return "circle_" + strconv.FormatInt(id, 10) + "vote_ranking_list"
}

func isNumberNil(value *int64) bool {
	return value == nil
}

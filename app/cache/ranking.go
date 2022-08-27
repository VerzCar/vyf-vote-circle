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
	identityId model.UserIdentityId,
	votes int64,
) error {

	rankingMap, err := c.getRankingMap(ctx, circleId, identityId)

	if err != nil {
		return err
	}

	ranking := rankingMap[identityId]

	ranking.Votes = votes

	placementNumber, err := c.setUserPlacement(ctx, circleId, identityId, model.VoteCount(votes))

	if err != nil {
		return err
	}

	ranking.Number = placementNumber

	err = c.setRanking(ctx, circleId, identityId, ranking)

	if err != nil {
		return err
	}

	return nil
}

func (c *redisCache) getRankingMap(
	ctx context.Context,
	circleId int64,
	identityId model.UserIdentityId,
) (model.RankingMap, error) {
	hashField := circleHashField(circleId)
	rankingMap := model.RankingMap{
		identityId: &model.Ranking{
			IdentityID: identityId,
			CircleID:   circleId,
		},
	}
	_, err := c.getHashJson(ctx, hashField, string(identityId), rankingMap[identityId])

	if err != nil {
		c.log.Errorf(
			"error reading ranking: reading field %s with current key %s: %s",
			hashField,
			identityId,
			err,
		)
		return nil, err
	}

	return rankingMap, nil
}

func (c *redisCache) setRanking(
	ctx context.Context,
	circleId int64,
	identityId model.UserIdentityId,
	ranking *model.Ranking,
) error {
	hashField := circleHashField(circleId)

	encodedRanking, err := ranking.MarshalBinary()
	if err != nil {
		c.log.Errorf(
			"error encoding ranking by setting ranking: for hash field %s map %v: %s",
			hashField,
			encodedRanking,
			err,
		)
		return err
	}

	rankingMap := map[string]interface{}{identityId.String(): encodedRanking}

	err = c.setHashMap(ctx, hashField, rankingMap)

	if err != nil {
		c.log.Errorf(
			"error setting ranking: setting field %s with ranking map %v: %s",
			hashField,
			rankingMap,
			err,
		)
		return err
	}

	return nil
}

func (c *redisCache) setUserPlacement(
	ctx context.Context,
	circleId int64,
	identityId model.UserIdentityId,
	voteCount model.VoteCount,
) (model.PlacementNumber, error) {
	hashField := circleVotesUserPlacementHashField(circleId)
	userPlacementMap := model.UserPlacementMap{}

	var placementNumber model.PlacementNumber
	placementNumber = 1
	previousVoteCount := voteCount - 1

	// try to find an existing entry for the previous vote
	entry, err := c.getHashJson(ctx, hashField, strconv.FormatInt(int64(previousVoteCount), 10), &userPlacementMap)

	if err != nil {
		c.log.Errorf(
			"error reading user placement in votes: reading field %s with key %d: %s",
			hashField,
			voteCount,
			err,
		)
		return 0, err
	}
	// if an entry exists for this count of vote
	// try to find the user id in the map and get the current
	// placing number from it.
	if entry.Exists {
		var userExist bool
		placementNumber, userExist = userPlacementMap[identityId]

		// count the placement one down for the next higher placement
		if placementNumber != 1 {
			placementNumber--
		}
		// the user must always exist, otherwise something has gone wrong
		if userExist {
			// delete the entry from the map
			// and overwrite the vote count hash entry
			delete(userPlacementMap, identityId)

			encodedUserPlacementMap, err := userPlacementMap.MarshalBinary()
			if err != nil {
				c.log.Errorf(
					"error encoding user placement map for vote count for previous vote: for hash field %s map %v: %s",
					hashField,
					userPlacementMap,
					err,
				)
				return 0, err
			}

			voteCountMap := map[string]interface{}{previousVoteCount.String(): encodedUserPlacementMap}
			err = c.setHashMap(ctx, hashField, voteCountMap)

			if err != nil {
				c.log.Errorf(
					"error setting vote count for previous vote: for hash field %s map %v: %s",
					hashField,
					userPlacementMap,
					err,
				)
				return 0, err
			}
		} else {
			return 0, fmt.Errorf(
				"user does not exist in user placement map - internal error: hashfield %s, voteCount %d, map %v",
				hashField,
				voteCount,
				userPlacementMap,
			)
		}
	}

	userPlacementMap = model.UserPlacementMap{}
	// try to find an existing entry for the next (current) vote
	entry, err = c.getHashJson(ctx, hashField, string(rune(voteCount)), &userPlacementMap)

	if err != nil {
		c.log.Errorf(
			"error getting vote count for current vote: for hash field %s map %v: %s",
			hashField,
			userPlacementMap,
			err,
		)
		return 0, err
	}

	// entry exists, get the placement from one of the existing users',
	// with one of the other users' placement number.
	// If it does not exist, set the placement number from the previous
	// found placement or default place 1.
	if entry.Exists {
		for key := range userPlacementMap {
			placementNumber = userPlacementMap[key]
			break
		}
	}

	userPlacementMap[identityId] = placementNumber
	encodedUserPlacementMap, err := userPlacementMap.MarshalBinary()
	if err != nil {
		c.log.Errorf(
			"error encoding user placement map for vote count for previous vote: for hash field %s map %v: %s",
			hashField,
			userPlacementMap,
			err,
		)
		return 0, err
	}

	voteCountMap := map[string]interface{}{voteCount.String(): encodedUserPlacementMap}
	err = c.setHashMap(ctx, hashField, voteCountMap)

	if err != nil {
		c.log.Errorf(
			"error setting vote count for previous vote: for hash field %s map %v: %s",
			hashField,
			userPlacementMap,
			err,
		)
		return 0, err
	}

	return placementNumber, nil
}

func circleHashField(id int64) string {
	return "circle_" + strconv.FormatInt(id, 10)
}

func circleVotesUserPlacementHashField(id int64) string {
	return "circle_" + strconv.FormatInt(id, 10) + "_votes_user_placement"
}

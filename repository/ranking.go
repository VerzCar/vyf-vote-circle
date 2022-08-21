package repository

import (
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database"
)

// RankingsByCircleId gets all rankings by the given circle id
func (s *storage) RankingsByCircleId(circleId int64) ([]*model.Ranking, error) {
	var rankings []*model.Ranking
	err := s.db.Where(&model.Ranking{CircleID: circleId}).
		Find(&rankings).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading rankings by circle id %d: %s", circleId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("rankings with circle id %d not found: %s", circleId, err)
		return nil, err
	}

	return rankings, nil
}

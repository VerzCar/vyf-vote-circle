package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
)

// CreateNewRanking based on given model.Ranking model
func (s *storage) CreateNewRanking(ranking *model.Ranking) (*model.Ranking, error) {
	if err := s.db.Create(ranking).Error; err != nil {
		s.log.Infof("error creating ranking: %s", err)
		return nil, err
	}

	return ranking, nil
}

func (s *storage) UpdateRanking(ranking *model.Ranking) (*model.Ranking, error) {
	if err := s.db.Model(ranking).Updates(ranking).Error; err != nil {
		s.log.Errorf("error updating ranking: %s", err)
		return nil, err
	}

	return ranking, nil
}

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

func (s *storage) RankingByCircleId(circleId int64, identityId string) (*model.Ranking, error) {
	ranking := &model.Ranking{}
	err := s.db.Where(&model.Ranking{IdentityID: identityId, CircleID: circleId}).
		First(ranking).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading ranking by circle id %d for user %s: %s", circleId, identityId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("ranking for user %s in circle %d not found: %s", identityId, circleId, err)
		return nil, err
	}

	return ranking, nil
}

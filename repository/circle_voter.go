package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
)

// CreateNewCircleVoter based on given CircleVoter model
func (s *storage) CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error) {
	if err := s.db.Create(voter).Error; err != nil {
		s.log.Infof("error creating circle voter: %s", err)
		return nil, err
	}

	return voter, nil
}

// CircleVoterByCircleId returns the queried circle voter in
// the circle based on the given circle id
func (s *storage) CircleVoterByCircleId(circleId int64, voterId string) (*model.CircleVoter, error) {
	circleVoter := &model.CircleVoter{}
	err := s.db.Where(&model.CircleVoter{Voter: voterId, CircleID: circleId}).
		First(circleVoter).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle voter by circle id %d: %s", circleId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circle voter with id %s in circle %d not found: %s", voterId, circleId, err)
		return nil, err
	}

	return circleVoter, nil
}

// IsVoterInCircle determines if the user exists in the circle voters list
func (s *storage) IsVoterInCircle(
	userIdentityId string,
	circle *model.Circle,
) (bool, error) {
	var count int64
	err := s.db.Model(&model.CircleVoter{}).
		Where(&model.CircleVoter{Voter: userIdentityId, CircleID: circle.ID}).
		Count(&count).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle voter id %d by circle id %d: %s", userIdentityId, circle.ID, err)
		return false, err
	case database.RecordNotFound(err):
		s.log.Infof("voter with id %d in circle id %d not found: %s", userIdentityId, circle.ID, err)
		return false, err
	}

	return count > 0, nil
}

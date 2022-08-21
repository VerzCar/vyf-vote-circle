package repository

import (
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database"
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

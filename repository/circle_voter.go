package repository

import (
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
)

// CreateNewCircleVoter based on given CircleVoter model
func (s *storage) CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error) {
	if err := s.db.Create(voter).Error; err != nil {
		s.log.Infof("error creating circle voter: %s", err)
		return nil, err
	}

	return voter, nil
}

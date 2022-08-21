package repository

import (
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
)

// CreateNewVote creates a new vote for the circle
func (s *storage) CreateNewVote(
	voterId int64,
	electedId int64,
	circleId int64,
) (*model.Vote, error) {
	vote := &model.Vote{
		VoterRefer:   voterId,
		ElectedRefer: electedId,
		CircleID:     circleId,
		CircleRefer:  &circleId,
	}
	if err := s.db.Create(vote).Error; err != nil {
		s.log.Infof("error creating vote in circle %d: %s", circleId, err)
		return nil, err
	}

	return vote, nil
}

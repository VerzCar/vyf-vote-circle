package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
)

// CreateNewVote creates a new vote for the circle
func (s *storage) CreateNewVote(
	voterId int64,
	electedId int64,
	circleId int64,
) (*model.Vote, error) {
	vote := &model.Vote{
		VoterRefer:     voterId,
		CandidateRefer: electedId,
		CircleID:       circleId,
		CircleRefer:    &circleId,
	}
	if err := s.db.Create(vote).Error; err != nil {
		s.log.Infof("error creating vote in circle %d: %s", circleId, err)
		return nil, err
	}

	return vote, nil
}

// ElectedVoterCountsByCircleId gets the number of votes for the elected id
func (s *storage) ElectedVoterCountsByCircleId(circleId int64, electedId int64) (int64, error) {
	var votes []*model.Vote
	err := s.db.Where(&model.Vote{CircleID: circleId, CircleRefer: &circleId, CandidateRefer: electedId}).
		Find(&votes).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading votes for elected user %d by circle id %d: %s", electedId, circleId, err)
		return 0, err
	case database.RecordNotFound(err):
		s.log.Infof("votes for elected user %d with circle id %d not found: %s", electedId, circleId, err)
		return 0, err
	}

	return int64(len(votes)), nil
}

// VoterElectedByCircleId query the given voter and elected id for the circle id and get the first result
func (s *storage) VoterElectedByCircleId(
	circleId int64,
	voterId int64,
	electedId int64,
) (*model.Vote, error) {
	vote := &model.Vote{}
	err := s.db.Where(
		&model.Vote{
			VoterRefer:     voterId,
			CandidateRefer: electedId,
			CircleID:       circleId,
			CircleRefer:    &circleId,
		},
	).First(vote).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf(
			"error reading vote for voter %d and elected %d by circle id %d: %s",
			voterId,
			electedId,
			circleId,
			err,
		)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof(
			"voter %s and elected %s by circle id %d not found: %s",
			voterId,
			electedId,
			circleId,
			err,
		)
		return nil, err
	}

	return vote, nil
}

// Votes gets all votes for the given circle id
func (s *storage) Votes(circleId int64) ([]*model.Vote, error) {
	var votes []*model.Vote
	err := s.db.Where(&model.Vote{CircleID: circleId, CircleRefer: &circleId}).Find(&votes).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading votes by circle id %d: %s", circleId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("votes with circle id %d not found: %s", circleId, err)
		return nil, err
	}

	return votes, nil
}

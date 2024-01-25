package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
)

// CreateNewVote creates a new vote for the circle
func (s *storage) CreateNewVote(
	voterId int64,
	candidateId int64,
	circleId int64,
) (*model.Vote, error) {
	vote := &model.Vote{
		VoterRefer:     voterId,
		CandidateRefer: candidateId,
		CircleID:       circleId,
		CircleRefer:    &circleId,
	}
	if err := s.db.Create(vote).Error; err != nil {
		s.log.Infof("error creating vote in circle %d: %s", circleId, err)
		return nil, err
	}

	return vote, nil
}

// Gets the number of votes for the candidate id
func (s *storage) CountsVotesOfCandidateByCircleId(circleId int64, candidateId int64) (int64, error) {
	var count int64
	err := s.db.Model(&model.Vote{}).
		Where(&model.Vote{CircleID: circleId, CircleRefer: &circleId, CandidateRefer: candidateId}).
		Count(&count).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading votes for candidate id %d by circle id %d: %s", candidateId, circleId, err)
		return 0, err
	case database.RecordNotFound(err):
		s.log.Infof("votes for candidate id %d with circle id %d not found: %s", candidateId, circleId, err)
		return 0, err
	}

	return count, nil
}

// Query the given voter and candidate id for the circle id and get the first result
func (s *storage) VoterCandidateByCircleId(
	circleId int64,
	voterId int64,
	candidateId int64,
) (*model.Vote, error) {
	vote := &model.Vote{}
	err := s.db.Where(
		&model.Vote{
			VoterRefer:     voterId,
			CandidateRefer: candidateId,
			CircleID:       circleId,
			CircleRefer:    &circleId,
		},
	).First(vote).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf(
			"error reading vote for voter %d and candidate %d by circle id %d: %s",
			voterId,
			candidateId,
			circleId,
			err,
		)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof(
			"voter %d and candidate %d by circle id %d not found: %s",
			voterId,
			candidateId,
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

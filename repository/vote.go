package repository

import (
	"context"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/cache"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Creates a new vote and updates all necessary tables
// in a transaction.
func (s *storage) CreateNewVote(
	ctx context.Context,
	circleId int64,
	voter *model.CircleVoter,
	candidate *model.CircleCandidate,
	upsertRankingCache cache.UpsertRankingCacheCallback,
) (*model.RankingResponse, int64, error) {
	vote := &model.Vote{
		VoterRefer:     voter.ID,
		CandidateRefer: candidate.ID,
		CircleID:       circleId,
		CircleRefer:    &circleId,
	}
	voteCount := int64(0)
	ranking := &model.Ranking{}
	var cachedRanking *model.RankingResponse

	err := s.db.Transaction(
		func(tx *gorm.DB) error {
			// create vote
			err := tx.Model(&model.Vote{}).Create(vote).Error

			if err != nil {
				s.log.Errorf("error creating vote in circle %d: %s", circleId, err)
				return err
			}

			// update the voters meta information
			voter.VotedFor = &candidate.Candidate
			err = tx.Model(voter).Update("voted_for", candidate.Candidate).Error

			if err != nil {
				s.log.Errorf("error updating voter id %d for circle id %d: %s", voter.ID, circleId, err)
				return err
			}

			err = tx.Model(&model.Vote{}).
				Where(&model.Vote{CircleID: circleId, CircleRefer: &circleId, CandidateRefer: candidate.ID}).
				Count(&voteCount).
				Error

			if err != nil {
				s.log.Errorf("error reading votes for candidate id %d by circle id %d: %s", candidate.ID, circleId, err)
				return err
			}

			ranking, err = s.txUpsertRanking(tx, circleId, voteCount, candidate)

			if err != nil {
				return err
			}

			cachedRanking, err = upsertRankingCache(ctx, circleId, candidate, ranking, voteCount)

			if err != nil {
				return err
			}

			// update ranking with newly indexed order
			err = tx.Model(&model.Ranking{ID: cachedRanking.ID}).
				Update("number", cachedRanking.Number).
				Error

			if err != nil {
				s.log.Errorf("error updating ranking for ranking id %d: %s", ranking.ID, err)
				return err
			}

			return nil
		},
	)

	if err != nil {
		s.log.Error("error creating vote: %s", err)
		return nil, 0, err
	}

	return cachedRanking, voteCount, nil
}

// Deletes a new vote and updates all necessary tables
// in a transaction.
func (s *storage) DeleteVote(
	ctx context.Context,
	circleId int64,
	vote *model.Vote,
	voter *model.CircleVoter,
	upsertRankingCache cache.UpsertRankingCacheCallback,
	removeRankingCache cache.RemoveRankingCacheCallback,
) (*model.RankingResponse, int64, error) {
	voteCount := int64(0)
	ranking := &model.Ranking{}
	var cachedRanking *model.RankingResponse

	err := s.db.Transaction(
		func(tx *gorm.DB) error {
			// delete vote
			err := tx.Model(&model.Vote{}).Delete(&model.Vote{ID: vote.ID}).Error

			if err != nil {
				s.log.Errorf("error deleting vote id %d: %s", vote.ID, err)
				return err
			}

			// update the voters meta information
			voter.VotedFor = nil
			err = tx.Model(voter).Update("voted_for", nil).Error

			if err != nil {
				s.log.Errorf("error updating voter id %d for circle id %d: %s", voter.ID, circleId, err)
				return err
			}

			err = tx.Model(&model.Vote{}).
				Where(&model.Vote{CircleID: circleId, CircleRefer: &circleId, CandidateRefer: vote.Candidate.ID}).
				Count(&voteCount).
				Error

			if err != nil {
				s.log.Errorf(
					"error reading votes for candidate id %d by circle id %d: %s",
					vote.Candidate.ID,
					circleId,
					err,
				)
				return err
			}

			// if still has votes update ranking
			if voteCount > 0 {
				ranking, err = s.txUpsertRanking(tx, circleId, voteCount, &vote.Candidate)

				if err != nil {
					return err
				}

				cachedRanking, err = upsertRankingCache(ctx, circleId, &vote.Candidate, ranking, voteCount)

				if err != nil {
					return err
				}

				// update ranking with newly indexed order
				err = tx.Model(&model.Ranking{ID: cachedRanking.ID}).
					Update("number", cachedRanking.Number).
					Error

				if err != nil {
					s.log.Errorf("error updating ranking for ranking id %d: %s", ranking.ID, err)
					return err
				}

				return nil
			}

			// if it does not have any votes delete ranking
			err = tx.Where(&model.Ranking{IdentityID: vote.Candidate.Candidate, CircleID: circleId}).
				First(ranking).
				Error

			if err != nil && !database.RecordNotFound(err) {
				s.log.Errorf(
					"error reading ranking by circle id %d for user %s: %s",
					circleId,
					vote.Candidate.Candidate,
					err,
				)
				return err
			}

			// Ranking exists
			if err == nil {
				err = tx.Model(&model.Ranking{}).
					Delete(&model.Ranking{}, ranking.ID).
					Error

				if err != nil {
					s.log.Errorf("error deleting ranking: %s", err)
					return err
				}
			}

			err = removeRankingCache(ctx, circleId, &vote.Candidate)

			if err != nil {
				return err
			}

			cachedRanking = &model.RankingResponse{
				ID: ranking.ID,
			}

			return nil
		},
	)

	if err != nil {
		s.log.Error("error deleting vote: %s", err)
		return nil, 0, err
	}

	return cachedRanking, voteCount, nil
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

// VoteByCircleId returns the queried vote in
// the circle based on the given voter id
func (s *storage) VoteByCircleId(
	circleId int64,
	voterId int64,
) (*model.Vote, error) {
	vote := &model.Vote{}
	err := s.db.Preload(clause.Associations).
		Where(&model.Vote{VoterRefer: voterId, CircleID: circleId}).
		First(vote).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading vote for voter id %d by circle id %d: %s", voterId, circleId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("vote with id %s in circle %d not found: %s", voterId, circleId, err)
		return nil, err
	}

	return vote, nil
}

// Determines if the voter already voted in the circle
func (s *storage) HasVoterVotedForCircle(
	circleId int64,
	voterId int64,
) (bool, error) {
	var count int64
	err := s.db.Model(&model.Vote{}).
		Where(
			&model.Vote{
				VoterRefer:  voterId,
				CircleID:    circleId,
				CircleRefer: &circleId,
			},
		).Count(&count).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf(
			"error reading vote for voter %d by circle id %d: %s",
			voterId,
			circleId,
			err,
		)
		return false, err
	case database.RecordNotFound(err):
		s.log.Infof(
			"vote for voter %d by circle id %d not found: %s",
			voterId,
			circleId,
			err,
		)
		return false, err
	}

	return count > 0, nil
}

// Votes gets all votes for the given circle id
func (s *storage) Votes(circleId int64) ([]*model.Vote, error) {
	var votes []*model.Vote
	err := s.db.Preload(clause.Associations).
		Where(&model.Vote{CircleID: circleId, CircleRefer: &circleId}).Find(&votes).Error

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

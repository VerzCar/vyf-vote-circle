package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
)

// based on given CircleCandidate model
func (s *storage) CreateNewCircleCandidate(candidate *model.CircleCandidate) (*model.CircleCandidate, error) {
	if err := s.db.Create(candidate).Error; err != nil {
		s.log.Infof("error creating circle candidate: %s", err)
		return nil, err
	}

	return candidate, nil
}

// UpdateCircleCandidate update circle voter based on given circle model
func (s *storage) UpdateCircleCandidate(candidate *model.CircleCandidate) (*model.CircleCandidate, error) {
	if err := s.db.Save(candidate).Error; err != nil {
		s.log.Errorf("error updating candidate: %s", err)
		return nil, err
	}

	return candidate, nil
}

// CircleCandidateByCircleId returns the queried circle candidate in
// the circle based on the given circle id
func (s *storage) CircleCandidateByCircleId(
	circleId int64,
	candidateId string,
) (*model.CircleCandidate, error) {
	circleCandidate := &model.CircleCandidate{}
	err := s.db.Where(&model.CircleCandidate{Candidate: candidateId, CircleID: circleId}).
		First(circleCandidate).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle candidate by circle id %d: %s", circleId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circle candidate with id %s in circle %d not found: %s", candidateId, circleId, err)
		return nil, err
	}

	return circleCandidate, nil
}

// IsCandidateInCircle determines if the user exists in the circle candidates list
func (s *storage) IsCandidateInCircle(
	userIdentityId string,
	circleId int64,
) (bool, error) {
	var count int64
	err := s.db.Model(&model.CircleCandidate{}).
		Where(&model.CircleCandidate{Candidate: userIdentityId, CircleID: circleId}).
		Count(&count).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle candidate id %s by circle id %d: %s", userIdentityId, circleId, err)
		return false, err
	case database.RecordNotFound(err):
		s.log.Infof("candidate with id %s in circle id %d not found: %s", userIdentityId, circleId, err)
		return false, err
	}

	return count > 0, nil
}

// Gets all the voters for the given circle id that matches the filter
// If no filter is provided all circle voters for the circle will be returned.
func (s *storage) CircleCandidatesFiltered(
	circleId int64,
	userIdentityId string,
	filterBy *model.CircleCandidatesFilterBy,
) ([]*model.CircleCandidate, error) {
	var circleCandidates []*model.CircleCandidate

	tx := s.db.Model(&model.CircleCandidate{}).
		Where(&model.CircleCandidate{CircleID: circleId}).
		Limit(100).
		Order("updated_at desc")

	if filterBy.Commitment != nil {
		tx.Where(&model.CircleCandidate{Commitment: *filterBy.Commitment})
	}

	if filterBy.HasBeenVoted != nil {
		if *filterBy.HasBeenVoted {
			tx.Where("voted_from IS NOT NULL")
		} else {
			tx.Where("voted_from IS NULL")
		}
	}

	if shouldContainUser := filterBy.ShouldContainUser != nil; !shouldContainUser {
		tx.Not(&model.CircleCandidate{Candidate: userIdentityId})
	}

	err := tx.Find(&circleCandidates).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle candidates: %s", err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circle candidates not found: %s", err)
		return nil, err
	}

	return circleCandidates, nil
}

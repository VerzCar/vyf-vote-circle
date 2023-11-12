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

// UpdateCircleVoter update circle voter based on given circle model
func (s *storage) UpdateCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error) {
	if err := s.db.Save(voter).Error; err != nil {
		s.log.Errorf("error updating circle: %s", err)
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

// CircleVotersFiltered gets all the voters for the given circle id that matches the filter
// If no filter is provided all circle voters for the circle will be returned.
func (s *storage) CircleVotersFiltered(
	circleId int64,
	userIdentityId string,
	filterBy *model.CircleVotersFilterBy,
) ([]*model.CircleVoter, error) {
	var circleVoters []*model.CircleVoter

	tx := s.db.Model(&model.CircleVoter{}).
		Where(&model.CircleVoter{CircleID: circleId}).
		Limit(100).
		Order("updated_at desc")

	if filterBy.Commitment != nil {
		tx.Where(&model.CircleVoter{Commitment: *filterBy.Commitment})
	}

	if filterBy.HasBeenVoted != nil {
		if *filterBy.HasBeenVoted {
			tx.Where("voted_from IS NOT NULL")
		} else {
			tx.Where("voted_from IS NULL")
		}
	}

	if shouldContainUser := filterBy.ShouldContainUser != nil; !shouldContainUser {
		tx.Not(&model.CircleVoter{Voter: userIdentityId})
	}

	err := tx.Find(&circleVoters).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle voters: %s", err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circle voters not found: %s", err)
		return nil, err
	}

	return circleVoters, nil
}

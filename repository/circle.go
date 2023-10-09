package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CircleById gets the circle by id
func (s *storage) CircleById(id int64) (*model.Circle, error) {
	circle := &model.Circle{}
	err := s.db.Preload(clause.Associations).
		Where(&model.Circle{ID: id}).First(circle).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle by id %d: %s", id, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circle with id %d not found: %s", id, err)
		return nil, err
	}

	return circle, nil
}

// Circles gets all the circles that have been create from the user
func (s *storage) Circles(userIdentityId string) ([]*model.Circle, error) {
	var circles []*model.Circle
	err := s.db.Preload(clause.Associations).
		Where(&model.Circle{CreatedFrom: userIdentityId}).
		Limit(int(s.config.Circle.MaxAmountPerUser)).
		Find(&circles).
		Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circles for user id %s: %s", userIdentityId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("circles with user id %s not found: %s", userIdentityId, err)
		return nil, err
	}

	return circles, nil
}

// UpdateCircle update circle based on given circle model
func (s *storage) UpdateCircle(circle *model.Circle) (*model.Circle, error) {
	if err := s.db.Save(circle).Error; err != nil {
		s.log.Errorf("error updating circle: %s", err)
		return nil, err
	}

	return circle, nil
}

// CreateNewCircle based on given circle model.
// The associations that come with it, will be created in the transaction accordingly.
func (s *storage) CreateNewCircle(circle *model.Circle) (*model.Circle, error) {
	err := s.db.Transaction(
		func(tx *gorm.DB) error {
			err := tx.Model(circle).Omit(clause.Associations).Create(circle).Error

			if err != nil {
				s.log.Error("error creating circle entry: %s", err)
				return err
			}

			circleVoters := circle.Voters

			for _, voter := range circleVoters {
				voter.CircleID = circle.ID
				voter.CircleRefer = &circle.ID
			}

			err = tx.Model(&model.CircleVoter{}).Create(circleVoters).Error

			if err != nil {
				s.log.Error("error creating circle voters entry: %s", err)
				return err
			}

			circle.Voters = circleVoters
			return nil
		},
	)

	if err != nil {
		s.log.Error("error creating circle: %s", err)
		return nil, err
	}

	return circle, nil
}

// CountCirclesOfUser determines how many circles the user already obtains
func (s *storage) CountCirclesOfUser(
	userIdentityId string,
) (int64, error) {
	var count int64
	err := s.db.Model(&model.Circle{}).
		Where(&model.Circle{CreatedFrom: userIdentityId}).
		Count(&count).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading circle count by user id %s: %s", userIdentityId, err)
		return 0, err
	case database.RecordNotFound(err):
		s.log.Infof("user with id %s in circles not found: %s", userIdentityId, err)
		return 0, err
	}

	return count, nil
}

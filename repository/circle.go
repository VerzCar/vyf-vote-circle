package repository

import (
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database"
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

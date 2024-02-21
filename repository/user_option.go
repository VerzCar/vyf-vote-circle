package repository

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/database"
)

// CreateNewUserOption based on given UserOption model
func (s *storage) CreateNewUserOption(option *model.UserOption) (*model.UserOption, error) {
	if err := s.db.Create(option).Error; err != nil {
		s.log.Infof("error creating user option: %s", err)
		return nil, err
	}

	return option, nil
}

func (s *storage) DeleteUserOption(optionId int64) error {
	if err := s.db.Model(&model.UserOption{}).Delete(&model.UserOption{}, optionId).Error; err != nil {
		s.log.Errorf("error deleting option: %s", err)
		return err
	}

	return nil
}

func (s *storage) UserOptionByUserIdentityId(userIdentityId string) (*model.UserOption, error) {
	option := &model.UserOption{}
	err := s.db.Where(&model.UserOption{IdentityID: userIdentityId}).
		First(option).Error

	switch {
	case err != nil && !database.RecordNotFound(err):
		s.log.Errorf("error reading user option by id %s: %s", userIdentityId, err)
		return nil, err
	case database.RecordNotFound(err):
		s.log.Infof("option for userIdentity with id %s not found: %s", userIdentityId, err)
		return nil, err
	}

	return option, nil
}

package api

import (
	"context"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
)

type UserOptionService interface {
	UserOption(
		ctx context.Context,
	) (*model.UserOptionResponse, error)
}

type UserOptionRepository interface {
	UserOptionByUserIdentityId(userIdentityId string) (*model.UserOption, error)
}

type userOptionService struct {
	storage UserOptionRepository
	config  *config.Config
	log     logger.Logger
}

func NewUserOptionService(
	userOptionRepo UserOptionRepository,
	config *config.Config,
	log logger.Logger,
) UserOptionService {
	return &userOptionService{
		storage: userOptionRepo,
		config:  config,
		log:     log,
	}
}

func (c *userOptionService) UserOption(
	ctx context.Context,
) (*model.UserOptionResponse, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	option, err := c.storage.UserOptionByUserIdentityId(authClaims.Subject)

	if err != nil && !database.RecordNotFound(err) {
		c.log.Errorf("error ruding query of user option")
		return c.defaultUserOption(), nil
	}

	if database.RecordNotFound(err) {
		return c.defaultUserOption(), nil
	}

	optionResponse := &model.UserOptionResponse{
		MaxCircles:    option.MaxCircles,
		MaxVoters:     option.MaxVoters,
		MaxCandidates: option.MaxCandidates,
		PrivateOption: model.UserPrivateOptionResponse{
			MaxVoters:     option.MaxPrivateVoters,
			MaxCandidates: option.MaxPrivateCandidates,
		},
	}

	return optionResponse, nil
}

func (c *userOptionService) defaultUserOption() *model.UserOptionResponse {
	return &model.UserOptionResponse{
		MaxCircles:    int(c.config.Circle.MaxAmountPerUser),
		MaxVoters:     c.config.Circle.MaxVoters,
		MaxCandidates: c.config.Circle.MaxCandidates,
		PrivateOption: model.UserPrivateOptionResponse{
			MaxVoters:     c.config.Circle.Private.MaxVoters,
			MaxCandidates: c.config.Circle.Private.MaxCandidates,
		},
	}
}

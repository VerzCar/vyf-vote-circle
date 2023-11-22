package api

import (
	"context"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
)

type CircleVoterService interface {
	CircleVotersFiltered(
		ctx context.Context,
		circleId int64,
		filterBy *model.CircleVotersFilterBy,
	) ([]*model.CircleVoter, *model.CircleVoter, error)
	CircleVoterCommitment(
		ctx context.Context,
		circleId int64,
		commitment model.Commitment,
	) (*model.Commitment, error)
}

type CircleVoterRepository interface {
	CircleVotersFiltered(
		circleId int64,
		userIdentityId string,
		filterBy *model.CircleVotersFilterBy,
	) ([]*model.CircleVoter, error)
	CircleVoterByCircleId(circleId int64, voterId string) (*model.CircleVoter, error)
	UpdateCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
}

type circleVoterService struct {
	storage CircleVoterRepository
	config  *config.Config
	log     logger.Logger
}

func NewCircleVoterService(
	circleVoterRepo CircleVoterRepository,
	config *config.Config,
	log logger.Logger,
) CircleVoterService {
	return &circleVoterService{
		storage: circleVoterRepo,
		config:  config,
		log:     log,
	}
}

func (c *circleVoterService) CircleVotersFiltered(
	ctx context.Context,
	circleId int64,
	filterBy *model.CircleVotersFilterBy,
) ([]*model.CircleVoter, *model.CircleVoter, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, nil, err
	}

	voters, err := c.storage.CircleVotersFiltered(circleId, authClaims.Subject, filterBy)

	if err != nil {
		return nil, nil, err
	}

	voter, err := c.storage.CircleVoterByCircleId(circleId, authClaims.Subject)

	if err != nil {
		return nil, nil, err
	}

	return voters, voter, nil
}

func (c *circleVoterService) CircleVoterCommitment(
	ctx context.Context,
	circleId int64,
	commitment model.Commitment,
) (*model.Commitment, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	voter, err := c.storage.CircleVoterByCircleId(circleId, authClaims.Subject)

	if err != nil {
		c.log.Errorf("error voter id %s not in circle: %s", authClaims.Subject, err)
		return nil, err
	}

	voter.Commitment = commitment
	_, err = c.storage.UpdateCircleVoter(voter)

	if err != nil {
		return nil, err
	}

	return &voter.Commitment, nil
}

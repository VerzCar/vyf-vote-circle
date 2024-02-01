package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
)

type CircleVoterService interface {
	CircleVotersFiltered(
		ctx context.Context,
		circleId int64,
		filterBy *model.CircleVotersFilterBy,
	) ([]*model.CircleVoter, *model.CircleVoter, error)
	CircleVoterJoinCircle(
		ctx context.Context,
		circleId int64,
	) (*model.CircleVoter, error)
	CircleVoterLeaveCircle(
		ctx context.Context,
		circleId int64,
	) error
}

type CircleVoterRepository interface {
	CircleVotersFiltered(
		circleId int64,
		filterBy *model.CircleVotersFilterBy,
	) ([]*model.CircleVoter, error)
	CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	CircleVoterByCircleId(circleId int64, voterId string) (*model.CircleVoter, error)
	IsVoterInCircle(userIdentityId string, circleId int64) (bool, error)
	CircleById(id int64) (*model.Circle, error)
	DeleteCircleVoter(voterId int64) error
}

type CircleVoterSubscription interface {
	CircleVoterChangedEvent(
		ctx context.Context,
		circleId int64,
		event *model.CircleVoterChangedEvent,
	) error
}

type circleVoterService struct {
	storage      CircleVoterRepository
	subscription CircleVoterSubscription
	config       *config.Config
	log          logger.Logger
}

func NewCircleVoterService(
	circleVoterRepo CircleVoterRepository,
	subscription CircleVoterSubscription,
	config *config.Config,
	log logger.Logger,
) CircleVoterService {
	return &circleVoterService{
		storage:      circleVoterRepo,
		subscription: subscription,
		config:       config,
		log:          log,
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

	voters, err := c.storage.CircleVotersFiltered(circleId, filterBy)

	if err != nil {
		return nil, nil, err
	}

	voter, err := c.storage.CircleVoterByCircleId(circleId, authClaims.Subject)

	if database.RecordNotFound(err) {
		return voters, nil, nil
	}

	if err != nil && !database.RecordNotFound(err) {
		return nil, nil, err
	}

	return voters, voter, nil
}

func (c *circleVoterService) CircleVoterJoinCircle(
	ctx context.Context,
	circleId int64,
) (*model.CircleVoter, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return nil, err
	}

	if circle.Private {
		err = fmt.Errorf("user cannot join private circle")
		return nil, err
	}

	isVoterInCircle, err := c.storage.IsVoterInCircle(authClaims.Subject, circleId)

	if err != nil {
		return nil, err
	}

	if isVoterInCircle {
		err = fmt.Errorf("user is already as voter in the circle")
		return nil, err
	}

	circleVoter := &model.CircleVoter{
		Voter:       authClaims.Subject,
		Circle:      circle,
		CircleRefer: &circle.ID,
		Commitment:  model.CommitmentCommitted,
	}
	voter, err := c.storage.CreateNewCircleVoter(circleVoter)

	if err != nil {
		c.log.Errorf("error adding voter to circle id %d: %s", circleId, err)
		return nil, err
	}

	voterEvent := createVoterChangedEvent(model.EventOperationCreated, voter)
	_ = c.subscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)

	return voter, nil
}

func (c *circleVoterService) CircleVoterLeaveCircle(
	ctx context.Context,
	circleId int64,
) error {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return err
	}

	voter, err := c.storage.CircleVoterByCircleId(circleId, authClaims.Subject)

	if err != nil {
		return fmt.Errorf("cannot leave as voter from cirlce")
	}

	err = c.storage.DeleteCircleVoter(voter.ID)

	if err != nil {
		c.log.Errorf(
			"error removing voter %s from circle id %d: %s",
			authClaims.Subject,
			circleId,
			err,
		)
		return fmt.Errorf("leaving as voter from cirlce failed")
	}

	voterEvent := createVoterChangedEvent(model.EventOperationDeleted, voter)
	_ = c.subscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)

	return nil
}

func createVoterChangedEvent(
	operation model.EventOperation,
	voter *model.CircleVoter,
) *model.CircleVoterChangedEvent {
	return &model.CircleVoterChangedEvent{
		Operation: operation,
		Voter: &model.CircleVoterResponse{
			ID:         voter.ID,
			Voter:      voter.Voter,
			VotedFor:   voter.VotedFor,
			Commitment: voter.Commitment,
			CreatedAt:  voter.CreatedAt,
			UpdatedAt:  voter.UpdatedAt,
		},
	}
}

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

type CircleCandidateService interface {
	CircleCandidatesFiltered(
		ctx context.Context,
		circleId int64,
		filterBy *model.CircleCandidatesFilterBy,
	) ([]*model.CircleCandidate, *model.CircleCandidate, error)
	CircleCandidateCommitment(
		ctx context.Context,
		circleId int64,
		commitment model.Commitment,
	) (*model.Commitment, error)
	CircleCandidateJoinCircle(
		ctx context.Context,
		circleId int64,
	) (*model.CircleCandidate, error)
	CircleCandidateLeaveCircle(
		ctx context.Context,
		circleId int64,
	) error
}

type CircleCandidateRepository interface {
	CircleCandidatesFiltered(
		circleId int64,
		filterBy *model.CircleCandidatesFilterBy,
	) ([]*model.CircleCandidate, error)
	CreateNewCircleCandidate(voter *model.CircleCandidate) (*model.CircleCandidate, error)
	CircleCandidateByCircleId(circleId int64, voterId string) (*model.CircleCandidate, error)
	UpdateCircleCandidate(voter *model.CircleCandidate) (*model.CircleCandidate, error)
	DeleteCircleCandidate(candidateId int64) error
	IsCandidateInCircle(userIdentityId string, circleId int64) (bool, error)
	CircleById(id int64) (*model.Circle, error)
}

type CircleCandidateSubscription interface {
	CircleCandidateChangedEvent(
		ctx context.Context,
		circleId int64,
		event *model.CircleCandidateChangedEvent,
	) error
}

type circleCandidateService struct {
	storage      CircleCandidateRepository
	subscription CircleCandidateSubscription
	config       *config.Config
	log          logger.Logger
}

func NewCircleCandidateService(
	circleCandidateRepo CircleCandidateRepository,
	subscription CircleCandidateSubscription,
	config *config.Config,
	log logger.Logger,
) CircleCandidateService {
	return &circleCandidateService{
		storage:      circleCandidateRepo,
		subscription: subscription,
		config:       config,
		log:          log,
	}
}

func (c *circleCandidateService) CircleCandidatesFiltered(
	ctx context.Context,
	circleId int64,
	filterBy *model.CircleCandidatesFilterBy,
) ([]*model.CircleCandidate, *model.CircleCandidate, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, nil, err
	}

	candidates, err := c.storage.CircleCandidatesFiltered(circleId, filterBy)

	if err != nil {
		return nil, nil, err
	}

	candidate, err := c.storage.CircleCandidateByCircleId(circleId, authClaims.Subject)

	if database.RecordNotFound(err) {
		return candidates, nil, nil
	}

	if err != nil && !database.RecordNotFound(err) {
		return nil, nil, err
	}

	return candidates, candidate, nil
}

// CircleCandidateCommitment updates the commitment of a circle candidate.
// It takes the following parameters:
// - ctx: the context.Context
// - circleId: the ID of the circle
// - commitment: the new commitment value
// It returns the updated commitment and an error if any.
func (c *circleCandidateService) CircleCandidateCommitment(
	ctx context.Context,
	circleId int64,
	commitment model.Commitment,
) (*model.Commitment, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	candidate, err := c.storage.CircleCandidateByCircleId(circleId, authClaims.Subject)

	if err != nil {
		c.log.Errorf("error candidate id %s not in circle: %s", authClaims.Subject, err)
		return nil, err
	}

	candidate.Commitment = commitment
	_, err = c.storage.UpdateCircleCandidate(candidate)

	if err != nil {
		return nil, err
	}

	candidateEvent := createCandidateChangedEvent(model.EventOperationUpdated, candidate)
	_ = c.subscription.CircleCandidateChangedEvent(ctx, circleId, candidateEvent)

	return &candidate.Commitment, nil
}

func (c *circleCandidateService) CircleCandidateJoinCircle(
	ctx context.Context,
	circleId int64,
) (*model.CircleCandidate, error) {
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

	IsCandidateInCircle, err := c.storage.IsCandidateInCircle(authClaims.Subject, circleId)

	if err != nil {
		return nil, err
	}

	if IsCandidateInCircle {
		err = fmt.Errorf("user is already as candidate in the circle")
		return nil, err
	}

	circleCandidate := &model.CircleCandidate{
		Candidate:   authClaims.Subject,
		Circle:      circle,
		CircleRefer: &circle.ID,
		Commitment:  model.CommitmentCommitted,
	}

	candidate, err := c.storage.CreateNewCircleCandidate(circleCandidate)

	if err != nil {
		c.log.Errorf("error adding candidate to circle id %d: %s", circleId, err)
		return nil, err
	}

	candidateEvent := createCandidateChangedEvent(model.EventOperationCreated, candidate)
	_ = c.subscription.CircleCandidateChangedEvent(ctx, circleId, candidateEvent)

	return candidate, nil
}

func (c *circleCandidateService) CircleCandidateLeaveCircle(
	ctx context.Context,
	circleId int64,
) error {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return err
	}

	candidate, err := c.storage.CircleCandidateByCircleId(circleId, authClaims.Subject)

	if err != nil {
		return fmt.Errorf("cannot leave as candidate from cirlce")
	}

	err = c.storage.DeleteCircleCandidate(candidate.ID)

	if err != nil {
		c.log.Errorf(
			"error removing candidate %s from circle id %d: %s",
			authClaims.Subject,
			circleId,
			err,
		)
		return fmt.Errorf("leaving as candidate from cirlce failed")
	}

	candidateEvent := createCandidateChangedEvent(model.EventOperationDeleted, candidate)
	_ = c.subscription.CircleCandidateChangedEvent(ctx, circleId, candidateEvent)

	return nil
}

func createCandidateChangedEvent(
	operation model.EventOperation,
	candidate *model.CircleCandidate,
) *model.CircleCandidateChangedEvent {
	return &model.CircleCandidateChangedEvent{
		Operation: operation,
		Candidate: &model.CircleCandidateResponse{
			ID:         candidate.ID,
			Candidate:  candidate.Candidate,
			Commitment: candidate.Commitment,
			CreatedAt:  candidate.CreatedAt,
			UpdatedAt:  candidate.UpdatedAt,
		},
	}
}

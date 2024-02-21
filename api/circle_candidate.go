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
	CircleCandidateAddToCircle(
		ctx context.Context,
		circleId int64,
		circleCandidateInput *model.CircleCandidateRequest,
	) (*model.CircleCandidate, error)
	CircleCandidateRemoveFromCircle(
		ctx context.Context,
		circleId int64,
		circleCandidateInput *model.CircleCandidateRequest,
	) (*model.CircleCandidate, error)
}

type CircleCandidateRepository interface {
	CircleCandidatesFiltered(
		circleId int64,
		filterBy *model.CircleCandidatesFilterBy,
	) ([]*model.CircleCandidate, error)
	CreateNewCircleCandidate(voter *model.CircleCandidate) (*model.CircleCandidate, error)
	CircleCandidateByCircleId(circleId int64, userIdentityId string) (*model.CircleCandidate, error)
	CircleCandidateCountByCircleId(
		circleId int64,
	) (int64, error)
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

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return nil, err
	}

	if !circle.Active {
		c.log.Infof(
			"tried to commit as candidate for an inactive circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return nil, fmt.Errorf("circle inactive")
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

	candidateEvent := CreateCandidateChangedEvent(model.EventOperationUpdated, candidate)
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

	if !circle.Active {
		c.log.Infof(
			"tried to join as candidate for an inactive circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return nil, fmt.Errorf("circle inactive")
	}

	candidatesCount, err := c.storage.CircleCandidateCountByCircleId(circleId)

	if err != nil {
		return nil, err
	}

	if candidatesCount > int64(c.config.Circle.MaxCandidates) {
		err = fmt.Errorf("circle has more than %d allowed candidates", c.config.Circle.MaxCandidates)
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

	candidateEvent := CreateCandidateChangedEvent(model.EventOperationCreated, candidate)
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

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return err
	}

	if !circle.Active {
		c.log.Infof(
			"tried to leave as candidate for an inactive circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return fmt.Errorf("circle inactive")
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

	candidateEvent := CreateCandidateChangedEvent(model.EventOperationDeleted, candidate)
	_ = c.subscription.CircleCandidateChangedEvent(ctx, circleId, candidateEvent)

	return nil
}

func (c *circleCandidateService) CircleCandidateAddToCircle(
	ctx context.Context,
	circleId int64,
	circleCandidateInput *model.CircleCandidateRequest,
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

	// checks whether user is eligible to add candidate to this circle
	if circle.CreatedFrom != authClaims.Subject {
		c.log.Infof(
			"user is not eligible to add candidate to circle: user %s, circle ID %d",
			authClaims.Subject,
			circle.ID,
		)
		err = fmt.Errorf("user is not eligible add candidate to circle")
		return nil, err
	}

	if !circle.Active {
		c.log.Infof(
			"tried to add candidate for an inactive circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return nil, fmt.Errorf("circle inactive")
	}

	IsCandidateInCircle, err := c.storage.IsCandidateInCircle(circleCandidateInput.Candidate, circleId)

	if err != nil {
		return nil, err
	}

	if IsCandidateInCircle {
		err = fmt.Errorf("user is already as candidate in the circle")
		return nil, err
	}

	circleCandidate := &model.CircleCandidate{
		Candidate:   circleCandidateInput.Candidate,
		Circle:      circle,
		CircleRefer: &circle.ID,
	}

	newCandidate, err := c.storage.CreateNewCircleCandidate(circleCandidate)

	if err != nil {
		c.log.Errorf("error adding candidate to circle id %d: %s", circleId, err)
		return nil, err
	}

	candidateEvent := CreateCandidateChangedEvent(model.EventOperationCreated, newCandidate)
	_ = c.subscription.CircleCandidateChangedEvent(ctx, circleId, candidateEvent)

	return newCandidate, nil
}

func (c *circleCandidateService) CircleCandidateRemoveFromCircle(
	ctx context.Context,
	circleId int64,
	circleCandidateInput *model.CircleCandidateRequest,
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

	// checks whether user is eligible to remove candidate from this circle
	if circle.CreatedFrom != authClaims.Subject {
		c.log.Infof(
			"user is not eligible to remove candidate from circle: user %s, circle ID %d",
			authClaims.Subject,
			circle.ID,
		)
		err = fmt.Errorf("user is not eligible to remove candidate from circle")
		return nil, err
	}

	if !circle.Active {
		c.log.Infof(
			"tried to add candidate for an inactive circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return nil, fmt.Errorf("circle inactive")
	}

	candidate, err := c.storage.CircleCandidateByCircleId(circleId, circleCandidateInput.Candidate)

	if err != nil && !database.RecordNotFound(err) {
		return nil, err
	}

	if !database.RecordNotFound(err) && candidate.Commitment == model.CommitmentRejected {
		c.log.Infof(
			"user has rejected candidacy for this circle: user %s, circle ID %d",
			authClaims.Subject,
			circle.ID,
		)
		err = fmt.Errorf("user cannot be removed as candidate from this circle")
		return nil, err
	}

	circleCandidate := &model.CircleCandidate{
		Candidate:   circleCandidateInput.Candidate,
		Circle:      circle,
		CircleRefer: &circle.ID,
	}

	newCandidate, err := c.storage.CreateNewCircleCandidate(circleCandidate)

	if err != nil {
		c.log.Errorf("error adding candidate to circle id %d: %s", circleId, err)
		return nil, err
	}

	candidateEvent := CreateCandidateChangedEvent(model.EventOperationCreated, newCandidate)
	_ = c.subscription.CircleCandidateChangedEvent(ctx, circleId, candidateEvent)

	return newCandidate, nil
}

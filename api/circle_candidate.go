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
}

type CircleCandidateRepository interface {
	CircleCandidatesFiltered(
		circleId int64,
		userIdentityId string,
		filterBy *model.CircleCandidatesFilterBy,
	) ([]*model.CircleCandidate, error)
	CreateNewCircleCandidate(voter *model.CircleCandidate) (*model.CircleCandidate, error)
	CircleCandidateByCircleId(circleId int64, voterId string) (*model.CircleCandidate, error)
	UpdateCircleCandidate(voter *model.CircleCandidate) (*model.CircleCandidate, error)
	IsCandidateInCircle(userIdentityId string, circleId int64) (bool, error)
	CircleById(id int64) (*model.Circle, error)
}

type circleCandidateService struct {
	storage CircleCandidateRepository
	config  *config.Config
	log     logger.Logger
}

func NewCircleCandidateService(
	circleCandidateRepo CircleCandidateRepository,
	config *config.Config,
	log logger.Logger,
) CircleCandidateService {
	return &circleCandidateService{
		storage: circleCandidateRepo,
		config:  config,
		log:     log,
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

	candidates, err := c.storage.CircleCandidatesFiltered(circleId, authClaims.Subject, filterBy)

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

	circleVoter := &model.CircleCandidate{
		Candidate:   authClaims.Subject,
		Circle:      circle,
		CircleRefer: &circle.ID,
		Commitment:  model.CommitmentCommitted,
	}
	voter, err := c.storage.CreateNewCircleCandidate(circleVoter)

	if err != nil {
		c.log.Errorf("error adding candidate to circle id %d: %s", circleId, err)
		return nil, err
	}

	return voter, nil
}

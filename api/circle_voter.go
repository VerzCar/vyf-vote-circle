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
	CircleVotersAddToCircle(
		ctx context.Context,
		circleId int64,
		circleVotersInput []*model.CircleVoterRequest,
	) ([]*model.CircleVoter, error)
	CircleVoterRemoveFromCircle(
		ctx context.Context,
		circleId int64,
		circleVoterInput *model.CircleVoterRequest,
	) error
}

type CircleVoterRepository interface {
	CircleVotersFiltered(
		circleId int64,
		filterBy *model.CircleVotersFilterBy,
	) ([]*model.CircleVoter, error)
	CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	CircleVoterByCircleId(circleId int64, userIdentityId string) (*model.CircleVoter, error)
	CircleVoterCountByCircleId(
		circleId int64,
	) (int64, error)
	IsVoterInCircle(userIdentityId string, circleId int64) (bool, error)
	CircleById(id int64) (*model.Circle, error)
	DeleteCircleVoter(voterId int64) error
	HasVoterVotedForCircle(
		circleId int64,
		voterId int64,
	) (bool, error)
}

type CircleVoterSubscription interface {
	CircleVoterChangedEvent(
		ctx context.Context,
		circleId int64,
		event *model.CircleVoterChangedEvent,
	) error
}

type CircleVoterOptionService interface {
	UserOption(
		ctx context.Context,
	) (*model.UserOptionResponse, error)
}

type circleVoterService struct {
	storage           CircleVoterRepository
	subscription      CircleVoterSubscription
	userOptionService CircleCandidateOptionService
	config            *config.Config
	log               logger.Logger
}

func NewCircleVoterService(
	circleVoterRepo CircleVoterRepository,
	subscription CircleVoterSubscription,
	userOptionService CircleCandidateOptionService,
	config *config.Config,
	log logger.Logger,
) CircleVoterService {
	return &circleVoterService{
		storage:           circleVoterRepo,
		subscription:      subscription,
		userOptionService: userOptionService,
		config:            config,
		log:               log,
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

	if !circle.IsEditable() {
		c.log.Infof(
			"tried to join as voter for an ineditable circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return nil, fmt.Errorf("circle is not editable")
	}

	votersCount, err := c.storage.CircleVoterCountByCircleId(circleId)

	if err != nil {
		return nil, err
	}

	userOption, _ := c.userOptionService.UserOption(ctx)

	if votersCount >= int64(userOption.MaxVoters) {
		err = fmt.Errorf("circle has more than %d allowed voters", userOption.MaxVoters)
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

	voterEvent := CreateVoterChangedEvent(model.EventOperationCreated, voter)
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

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return err
	}

	if !circle.IsEditable() {
		c.log.Infof(
			"tried to leave as voter for an ineditable circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return fmt.Errorf("circle is not editable")
	}

	voter, err := c.storage.CircleVoterByCircleId(circleId, authClaims.Subject)

	if err != nil {
		return fmt.Errorf("cannot leave as voter from cirlce")
	}

	hasVoted, err := c.storage.HasVoterVotedForCircle(circleId, voter.ID)

	if err != nil && !database.RecordNotFound(err) {
		return fmt.Errorf("cannot leave as voter from cirlce")
	}
	if err == nil && hasVoted {
		return fmt.Errorf("voter has voted")
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

	voterEvent := CreateVoterChangedEvent(model.EventOperationDeleted, voter)
	_ = c.subscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)

	return nil
}

func (c *circleVoterService) CircleVotersAddToCircle(
	ctx context.Context,
	circleId int64,
	circleVotersInput []*model.CircleVoterRequest,
) ([]*model.CircleVoter, error) {
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
			"user is not eligible to add voter to circle: user %s, circle ID %d",
			authClaims.Subject,
			circle.ID,
		)
		err = fmt.Errorf("user is not eligible to add voter to circle")
		return nil, err
	}

	if !circle.IsEditable() {
		c.log.Infof(
			"tried to add voter for an ineditable circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return nil, fmt.Errorf("circle is not editable")
	}

	updatedVoters := make([]*model.CircleVoter, 0)

	for _, voter := range circleVotersInput {
		newVoter, err := c.addVoterToCircle(ctx, circle, voter.Voter)

		if err != nil {
			return nil, err
		}

		updatedVoters = append(updatedVoters, newVoter)

		voterEvent := CreateVoterChangedEvent(model.EventOperationCreated, newVoter)
		_ = c.subscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)
	}

	return updatedVoters, nil
}

func (c *circleVoterService) CircleVoterRemoveFromCircle(
	ctx context.Context,
	circleId int64,
	circleVoterInput *model.CircleVoterRequest,
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

	// checks whether user is eligible to add candidate to this circle
	if circle.CreatedFrom != authClaims.Subject {
		c.log.Infof(
			"user is not eligible to add voter to circle: user %s, circle ID %d",
			authClaims.Subject,
			circle.ID,
		)
		err = fmt.Errorf("user is not eligible to add voter to circle")
		return err
	}

	if !circle.IsEditable() {
		c.log.Infof(
			"tried to add voter for an ineditable circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return fmt.Errorf("circle is not editable")
	}

	voter, err := c.storage.CircleVoterByCircleId(circleId, circleVoterInput.Voter)

	if err != nil {
		return fmt.Errorf("cannot remove voter from circle")
	}

	hasVoted, err := c.storage.HasVoterVotedForCircle(circleId, voter.ID)

	if err != nil && !database.RecordNotFound(err) {
		return fmt.Errorf("cannot leave as voter from circle")
	}
	if err == nil && hasVoted {
		return fmt.Errorf("voter has voted")
	}

	err = c.storage.DeleteCircleVoter(voter.ID)

	if err != nil {
		c.log.Errorf(
			"error removing voter %s from circle id %d: %s",
			authClaims.Subject,
			circleId,
			err,
		)
		return fmt.Errorf("removing voter from cirlce failed")
	}

	voterEvent := CreateVoterChangedEvent(model.EventOperationDeleted, voter)
	_ = c.subscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)

	return nil
}

func (c *circleVoterService) addVoterToCircle(
	ctx context.Context,
	circle *model.Circle,
	voterIdentId string,
) (*model.CircleVoter, error) {
	IsCandidateInCircle, err := c.storage.IsVoterInCircle(voterIdentId, circle.ID)

	if err != nil {
		return nil, err
	}

	if IsCandidateInCircle {
		err = fmt.Errorf("user is already as voter in the circle")
		return nil, err
	}

	candidatesCount, err := c.storage.CircleVoterCountByCircleId(circle.ID)

	if err != nil {
		return nil, fmt.Errorf("count of voter failure")
	}

	userOption, _ := c.userOptionService.UserOption(ctx)

	if candidatesCount >= int64(userOption.MaxCandidates) {
		err = fmt.Errorf("circle has more than %d allowed voters", userOption.MaxVoters)
		return nil, err
	}

	circleVoter := &model.CircleVoter{
		Voter:       voterIdentId,
		Circle:      circle,
		CircleRefer: &circle.ID,
	}

	newVoter, err := c.storage.CreateNewCircleVoter(circleVoter)

	if err != nil {
		c.log.Errorf("error adding voter to circle id %d: %s", circle.ID, err)
		return nil, err
	}

	return newVoter, nil
}

package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
	"time"
)

type CircleService interface {
	Circle(
		ctx context.Context,
		circleId int64,
	) (*model.Circle, error)
	Circles(
		ctx context.Context,
	) ([]*model.Circle, error)
	UpdateCircle(
		ctx context.Context,
		circleUpdateRequest *model.CircleUpdateRequest,
	) (*model.Circle, error)
	CreateCircle(
		ctx context.Context,
		circleCreateRequest *model.CircleCreateRequest,
	) (*model.Circle, error)
	EligibleToBeInCircle(
		ctx context.Context,
		circleId int64,
	) (bool, error)
	AddToGlobalCircle(
		ctx context.Context,
	) error
}

type CircleRepository interface {
	CircleById(id int64) (*model.Circle, error)
	Circles(userIdentityId string) ([]*model.Circle, error)
	UpdateCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	IsVoterInCircle(userIdentityId string, circle *model.Circle) (bool, error)
	CountCirclesOfUser(userIdentityId string) (int64, error)
}

type circleService struct {
	storage CircleRepository
	config  *config.Config
	log     logger.Logger
}

func NewCircleService(
	circleRepo CircleRepository,
	config *config.Config,
	log logger.Logger,
) CircleService {
	return &circleService{
		storage: circleRepo,
		config:  config,
		log:     log,
	}
}

func (c *circleService) Circle(
	ctx context.Context,
	circleId int64,
) (*model.Circle, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return nil, err
	}

	eligibleToBeInCircle, err := c.eligibleToBeInCircle(authClaims.Subject, circle)

	if err != nil {
		return nil, err
	}

	if !eligibleToBeInCircle {
		c.log.Infof("user is not eligible to be in circle: user %s, circle ID %d", authClaims.Subject, circle.ID)
		err = fmt.Errorf("user is not eligible to be in circle")
		return nil, err
	}

	return circle, nil
}

// Circles will determine all the circles the authenticated
// user has and returns the circles as a list.
// If the user hasn't any circles the return value will be empty.
func (c *circleService) Circles(
	ctx context.Context,
) ([]*model.Circle, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circles, err := c.storage.Circles(authClaims.Subject)

	switch {
	case err != nil && !database.RecordNotFound(err):
		{
			return nil, err
		}
	case database.RecordNotFound(err):
		{
			return nil, nil
		}
	}

	return circles, nil
}

func (c *circleService) UpdateCircle(
	ctx context.Context,
	circleUpdateRequest *model.CircleUpdateRequest,
) (*model.Circle, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circle, err := c.storage.CircleById(circleUpdateRequest.ID)

	if err != nil {
		return nil, err
	}

	userId := authClaims.Subject

	// checks whether user is eligible to update this circle
	if circle.CreatedFrom != userId {
		c.log.Infof("user is not eligible to update circle: user %s, circle ID %d", userId, circle.ID)
		err = fmt.Errorf("user is not eligible to update circle")
		return nil, err
	}

	// if circle is not active anymore, it can't be updated
	if !circle.Active {
		c.log.Infof("user try to update inactive circle: user %s, circle ID %d", userId, circle.ID)
		err = fmt.Errorf("circle is not active")
		return nil, err
	}

	// if circle should be deleted, deactivated it and return deactivated circle
	if circleUpdateRequest.Delete != nil {
		if *circleUpdateRequest.Delete {
			err := c.inactivateCircle(circle)

			if err != nil {
				c.log.Warnf("could not deactivate circle, error: circle ID %d, error %s", circle.ID, err)
				return nil, err
			}

			return circle, nil
		}
	}

	// check if new valid until time is given and is in the future from now on
	// otherwise check if current valid until time has expired
	if circleUpdateRequest.ValidUntil != nil {
		currentTime := time.Now()
		if currentTime.After(*circleUpdateRequest.ValidUntil) {
			err = fmt.Errorf("valid until time must be in the future from now")
			return nil, err
		}
		circle.ValidUntil = circleUpdateRequest.ValidUntil
	}

	if circleUpdateRequest.Name != nil {
		circle.Name = *circleUpdateRequest.Name
	}

	if circleUpdateRequest.Private != nil {
		circle.Private = *circleUpdateRequest.Private
	}

	if circleUpdateRequest.ImageSrc != nil {
		circle.ImageSrc = *circleUpdateRequest.ImageSrc
	}

	if circleUpdateRequest.Voters != nil {
		var circleVoters []*model.CircleVoter
		for _, voter := range circleUpdateRequest.Voters {
			circleVoter := &model.CircleVoter{
				Voter:       voter.Voter,
				Circle:      circle,
				CircleRefer: &circle.ID,
			}
			circleVoters = append(circleVoters, circleVoter)
		}
		circle.Voters = circleVoters
	}

	circle, err = c.storage.UpdateCircle(circle)

	if err != nil {
		return nil, fmt.Errorf("error updating circle: %s", err)
	}

	return circle, nil
}

func (c *circleService) CreateCircle(
	ctx context.Context,
	circleCreateRequest *model.CircleCreateRequest,
) (*model.Circle, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circlesCount, err := c.storage.CountCirclesOfUser(authClaims.Subject)

	if err != nil {
		return nil, err
	}

	if circlesCount > c.config.Circle.MaxAmountPerUser {
		err = fmt.Errorf("user has more than %d allowed circles", c.config.Circle.MaxAmountPerUser)
		return nil, err
	}

	newCircle := &model.Circle{
		Name:        circleCreateRequest.Name,
		CreatedFrom: authClaims.Subject,
	}

	if circleCreateRequest.Private != nil {
		newCircle.Private = *circleCreateRequest.Private
	}

	// check if new valid until time is given and is in the future from now on
	if circleCreateRequest.ValidUntil != nil {
		currentTime := time.Now()
		if currentTime.After(*circleCreateRequest.ValidUntil) {
			err = fmt.Errorf("valid until time must be in the future from now")
			return nil, err
		}
		newCircle.ValidUntil = circleCreateRequest.ValidUntil
	}

	if len(circleCreateRequest.Voters) <= 0 {
		err = fmt.Errorf("voters for circle are not given")
		return nil, err
	}

	if len(circleCreateRequest.Voters) > c.config.Circle.MaxVoters {
		err = fmt.Errorf("circle has more than %d allowed voters", c.config.Circle.MaxVoters)
		return nil, err
	}

	var circleVoters = c.createCircleVoterList(authClaims.Subject, circleCreateRequest.Voters)
	newCircle.Voters = circleVoters

	circle, err := c.storage.CreateNewCircle(newCircle)

	if err != nil {
		return nil, fmt.Errorf("error creating circle: %s", err)
	}

	return circle, nil
}

// EligibleToBeInCircle checks whether the user is allowed to be in the circle.
// Either, if the user itself has created the circle or if it is one of the voters.
func (c *circleService) EligibleToBeInCircle(
	ctx context.Context,
	circleId int64,
) (bool, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return false, err
	}

	userIdentityId := authClaims.Subject

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return false, err
	}

	return c.eligibleToBeInCircle(userIdentityId, circle)
}

// AddToGlobalCircle adds the user to the global circle.
// It checks whether the user exists in the voters list and add it to
// the global list if not. This has to be done if a new user
// creates an account and therefore must be inserted in the global list
// for the first time.
func (c *circleService) AddToGlobalCircle(
	ctx context.Context,
) error {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return err
	}

	circle, err := c.storage.CircleById(1)

	if err != nil {
		return err
	}

	isVoterInCircle, err := c.storage.IsVoterInCircle(authClaims.Subject, circle)

	if err != nil {
		return err
	}

	if isVoterInCircle {
		return nil
	}

	circleVoter := &model.CircleVoter{
		Voter:       authClaims.Subject,
		Circle:      circle,
		CircleRefer: &circle.ID,
		Commitment:  model.CommitmentCommitted,
	}
	circleVoter, err = c.storage.CreateNewCircleVoter(circleVoter)

	if err != nil {
		c.log.Errorf("error adding voter to global circle: %s", err)
		return err
	}

	return nil
}

func (c *circleService) eligibleToBeInCircle(
	userIdentityId string,
	circle *model.Circle,
) (bool, error) {
	if !circle.Private {
		return true, nil
	}

	if userIdentityId == circle.CreatedFrom {
		return true, nil
	}

	return c.storage.IsVoterInCircle(userIdentityId, circle)
}

// inactivateCircle of the given circle and updates it in the database to an
// inactive circle.
func (c *circleService) inactivateCircle(
	circle *model.Circle,
) error {
	circle.Active = false
	circle, err := c.storage.UpdateCircle(circle)

	if err != nil {
		return err
	}

	return nil
}

// createCircleVoterList based on the given createdFrom and the
// circleVoterInputs. It removes all the duplicates from the
// circleVoterInputs and add the createdFrom id to the list.
func (c *circleService) createCircleVoterList(
	createdFrom string,
	circleVoterInputs []*model.CircleVoterRequest,
) []*model.CircleVoter {
	var voterIdList []string

	voterIdList = append(voterIdList, createdFrom)
	for _, voter := range circleVoterInputs {
		voterIdList = append(voterIdList, voter.Voter)
	}

	voterIdList = removeDuplicateStr(voterIdList)

	var circleVoters []*model.CircleVoter
	// add the given voters to the circle voters
	for _, voter := range voterIdList {
		circleVoter := &model.CircleVoter{
			Voter: voter,
		}
		circleVoters = append(circleVoters, circleVoter)
	}

	return circleVoters
}

func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	var list []string
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

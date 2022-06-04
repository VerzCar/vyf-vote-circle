package api

import (
	"context"
	"fmt"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	routerContext "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router/ctx"
	"time"
)

type CircleService interface {
	Circle(
		ctx context.Context,
		circleId int64,
	) (*model.Circle, error)
	UpdateCircle(
		ctx context.Context,
		circleId int64,
		circleUpdateInput *model.CircleUpdateInput,
	) (*model.Circle, error)
	CreateCircle(
		ctx context.Context,
		circleCreateInput *model.CircleCreateInput,
	) (*model.Circle, error)
}

type CircleRepository interface {
	CircleById(id int64) (*model.Circle, error)
	UpdateCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircle(circle *model.Circle) (*model.Circle, error)
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

	if !c.eligibleToBeInCircle(
		authClaims.Subject,
		circle,
	) {
		c.log.Infof("user is not eligible to be in circle: user %s, circle ID %d", authClaims.Subject, circle.ID)
		err = fmt.Errorf("user is not eligible to be in circle")
		return nil, err
	}

	if circle.Active {
		if c.hasValidationTimeExpired(circle) {
			if err := c.inactivateCircle(circle); err != nil {
				c.log.Warnf("circle has validateValidationTime error: circle ID %d, error %s", circle.ID, err)
			}
		}
	}

	return circle, nil
}

func (c *circleService) UpdateCircle(
	ctx context.Context,
	circleId int64,
	circleUpdateInput *model.CircleUpdateInput,
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
	if circleUpdateInput.Delete != nil {
		err := c.inactivateCircle(circle)

		if err != nil {
			c.log.Warnf("circle has validateValidationTime error: circle ID %d, error %s", circle.ID, err)
			return nil, err
		}

		return circle, nil
	}

	// check if new valid until time is given and is in the future from now on
	// otherwise check if current valid until time has expired
	if circleUpdateInput.ValidUntil != nil {
		currentTime := time.Now()
		if currentTime.After(*circleUpdateInput.ValidUntil) {
			err = fmt.Errorf("valid until time must be in the future from now")
			return nil, err
		}
		circle.ValidUntil = circleUpdateInput.ValidUntil
	} else if c.hasValidationTimeExpired(circle) {
		circle.Active = false
	}

	if circleUpdateInput.Name != nil {
		circle.Name = *circleUpdateInput.Name
	}

	if circleUpdateInput.Private != nil {
		circle.Private = *circleUpdateInput.Private
	}

	if circleUpdateInput.Voters != nil {
		var circleVoters []*model.CircleVoter
		for _, voter := range circleUpdateInput.Voters {
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
	circleCreateInput *model.CircleCreateInput,
) (*model.Circle, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	newCircle := &model.Circle{
		Name:        circleCreateInput.Name,
		CreatedFrom: authClaims.Subject,
	}

	if circleCreateInput.Private != nil {
		newCircle.Private = *circleCreateInput.Private
	}

	// check if new valid until time is given and is in the future from now on
	if circleCreateInput.ValidUntil != nil {
		currentTime := time.Now()
		if currentTime.After(*circleCreateInput.ValidUntil) {
			err = fmt.Errorf("valid until time must be in the future from now")
			return nil, err
		}
		newCircle.ValidUntil = circleCreateInput.ValidUntil
	}

	if len(circleCreateInput.Voters) <= 0 {
		err = fmt.Errorf("voters for circle are not given")
		return nil, err
	}

	var circleVoters []*model.CircleVoter
	for _, voter := range circleCreateInput.Voters {
		circleVoter := &model.CircleVoter{
			Voter: voter.Voter,
		}
		circleVoters = append(circleVoters, circleVoter)
	}
	newCircle.Voters = circleVoters

	circle, err := c.storage.CreateNewCircle(newCircle)

	if err != nil {
		return nil, fmt.Errorf("error creating circle: %s", err)
	}

	return circle, nil
}

// eligibleToBeInCircle checks whether the user is allowed to be in the circle.
// Either, if the user itself has created the circle or if it is one of the voters.
func (c *circleService) eligibleToBeInCircle(
	userIdentityId string,
	circle *model.Circle,
) bool {
	if userIdentityId == circle.CreatedFrom {
		return true
	}

	for _, voter := range circle.Voters {
		if voter.Voter == userIdentityId {
			return true
		}
	}

	return false
}

// hasValidationTimeExpired checks if the validationUntil time has expired.
// If an validUntil time is set, it will be compared to the current time
// and validated if it has expired.
// Returns false if either no validationUntil time is set
// or the validation time has not expired, otherwise true.
func (c *circleService) hasValidationTimeExpired(
	circle *model.Circle,
) bool {

	if circle.ValidUntil == nil {
		return false
	}

	currentTime := time.Now()
	validUntilTime := *circle.ValidUntil

	if currentTime.After(validUntilTime) {
		return true
	}

	return false
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

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
}

type CircleRepository interface {
	CircleById(id int64) (*model.Circle, error)
	UpdateCircle(circle *model.Circle) (*model.Circle, error)
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
	return nil, nil
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

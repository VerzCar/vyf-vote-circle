package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
	"github.com/VerzCar/vyf-vote-circle/utils"
	"strings"
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
	CirclesOpenCommitments(
		ctx context.Context,
	) ([]*model.CirclePaginated, error)
	CirclesFiltered(
		ctx context.Context,
		name *string,
	) ([]*model.CirclePaginated, error)
	CirclesOfInterest(
		ctx context.Context,
	) ([]*model.CirclePaginated, error)
	UpdateCircle(
		ctx context.Context,
		circleId int64,
		circleUpdateRequest *model.CircleUpdateRequest,
	) (*model.Circle, error)
	CreateCircle(
		ctx context.Context,
		circleCreateRequest *model.CircleCreateRequest,
	) (*model.Circle, error)
	DeleteCircle(
		ctx context.Context,
		circleId int64,
	) error
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
	CirclesByIds(circleIds []int64) ([]*model.CirclePaginated, error)
	Circles(userIdentityId string) ([]*model.Circle, error)
	CirclesFiltered(name string) ([]*model.CirclePaginated, error)
	CirclesOfInterest(userIdentityId string) ([]*model.CirclePaginated, error)
	UpdateCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircle(circle *model.Circle) (*model.Circle, error)
	CreateNewCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	IsVoterInCircle(userIdentityId string, circleId int64) (bool, error)
	IsCandidateInCircle(
		userIdentityId string,
		circleId int64,
	) (bool, error)
	CountCirclesOfUser(userIdentityId string) (int64, error)
	CircleCandidatesOpenCommitments(
		userIdentityId string,
	) ([]*model.CircleCandidate, error)
	ExistVoteByCircleId(
		circleId int64,
	) (bool, error)
}

type CircleUserOptionService interface {
	UserOption(
		ctx context.Context,
	) (*model.UserOptionResponse, error)
}

type circleService struct {
	storage           CircleRepository
	userOptionService CircleUserOptionService
	config            *config.Config
	log               logger.Logger
}

func NewCircleService(
	circleRepo CircleRepository,
	userOptionService CircleUserOptionService,
	config *config.Config,
	log logger.Logger,
) CircleService {
	return &circleService{
		storage:           circleRepo,
		userOptionService: userOptionService,
		config:            config,
		log:               log,
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

func (c *circleService) CirclesOpenCommitments(
	ctx context.Context,
) ([]*model.CirclePaginated, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circleCandidates, err := c.storage.CircleCandidatesOpenCommitments(authClaims.Subject)

	switch {
	case err != nil && !database.RecordNotFound(err):
		{
			return nil, err
		}
	case database.RecordNotFound(err) || len(circleCandidates) <= 0:
		{
			return nil, nil
		}
	}

	var circleIds []int64

	for _, candidate := range circleCandidates {
		circleIds = append(circleIds, candidate.CircleID)
	}

	circles, err := c.storage.CirclesByIds(circleIds)

	if err != nil {
		return nil, err
	}

	return circles, nil
}

// CirclesFiltered takes a name parameter and returns a list of circles that
// match the given name, filtered from the authenticated user's circles. If there
// are no matching circles, the return value will be empty.
// Parameters:
// - ctx: The context.Context object for the request.
// - name: A pointer to a string representing the name to filter the circles by.
// Returns:
// - []*model.CirclePaginated: A list of circles that match the given name.
// - error: An error if any occurred during the execution.
func (c *circleService) CirclesFiltered(
	ctx context.Context,
	name *string,
) ([]*model.CirclePaginated, error) {
	circles, err := c.storage.CirclesFiltered(*name)

	if err != nil {
		return nil, err
	}

	return circles, nil
}

// CirclesOfInterest determines all the circles of interest for the authenticated user and returns them as a list.
// If the user doesn't have any circles of interest, the return value will be empty.
func (c *circleService) CirclesOfInterest(
	ctx context.Context,
) ([]*model.CirclePaginated, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return nil, err
	}

	circles, err := c.storage.CirclesOfInterest(authClaims.Subject)

	if err != nil {
		return nil, err
	}

	return circles, nil
}

func (c *circleService) UpdateCircle(
	ctx context.Context,
	circleId int64,
	circleUpdateRequest *model.CircleUpdateRequest,
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

	// if circle is not editable anymore, it can't be updated
	if !circle.IsEditable() {
		c.log.Infof("user try to update inactive or closed circle: user %s, circle ID %d", userId, circle.ID)
		err = fmt.Errorf("circle is not editable")
		return nil, err
	}

	currentTime := currentTruncatedTime()
	// check if new valid from time is given and is in the future from now on
	// otherwise check if current valid from time has expired
	if circleUpdateRequest.ValidFrom != nil {

		if circle.Stage == model.CircleStageHot {
			hasVotes, err := c.storage.ExistVoteByCircleId(circle.ID)

			if hasVotes || err != nil {
				return nil, fmt.Errorf("circle is in hot stage and cannot be updated in time range")
			}
		}

		validFrom, err := extractValidFrom(currentTime, *circleUpdateRequest.ValidFrom)

		if err != nil {
			return nil, err
		}

		circle.ValidFrom = *validFrom
	}

	// check if new valid until time is given and is in the future from now on
	// otherwise check if current valid until time has expired
	if circleUpdateRequest.ValidUntil != nil {
		validUntil, err := extractValidUntil(currentTime, circle.ValidFrom, *circleUpdateRequest.ValidUntil)

		if err != nil {
			return nil, err
		}

		circle.ValidUntil = validUntil
	} else {
		circle.ValidUntil = nil
	}

	if circleUpdateRequest.Name != nil {
		circle.Name = strings.TrimSpace(*circleUpdateRequest.Name)
	}

	if circleUpdateRequest.ImageSrc != nil {
		circle.ImageSrc = *circleUpdateRequest.ImageSrc
	}

	if circleUpdateRequest.Description != nil {
		circle.Description = strings.TrimSpace(*circleUpdateRequest.Description)
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

	userOption, _ := c.userOptionService.UserOption(ctx)

	if circlesCount >= int64(userOption.MaxCircles) {
		err = fmt.Errorf("user has more than %d allowed circles", userOption.MaxCircles)
		return nil, err
	}

	newCircle := &model.Circle{
		Name:        strings.TrimSpace(circleCreateRequest.Name),
		CreatedFrom: authClaims.Subject,
	}

	if circleCreateRequest.Private != nil {
		newCircle.Private = *circleCreateRequest.Private
	}

	if newCircle.Private && len(circleCreateRequest.Voters) <= 0 {
		err = fmt.Errorf("circle must contain at least one voter if private")
		return nil, err
	}

	if newCircle.Private && len(circleCreateRequest.Voters) >= userOption.PrivateOption.MaxVoters {
		err = fmt.Errorf(
			"circle has %d more than %d allowed voters",
			len(circleCreateRequest.Voters),
			userOption.PrivateOption.MaxVoters,
		)
		return nil, err
	}

	if newCircle.Private && len(circleCreateRequest.Candidates) <= 0 {
		err = fmt.Errorf("circle must contain at least one candidate if private")
		return nil, err
	}

	if newCircle.Private && len(circleCreateRequest.Candidates) >= userOption.PrivateOption.MaxCandidates {
		err = fmt.Errorf(
			"circle has %d more than %d allowed candidates",
			len(circleCreateRequest.Candidates),
			userOption.PrivateOption.MaxCandidates,
		)
		return nil, err
	}

	if len(circleCreateRequest.Voters) > 0 {
		if !newCircle.Private {
			err = fmt.Errorf("circle must be private to add voters")
			return nil, err
		}

		circleVoters := c.createCircleVoterList(circleCreateRequest.Voters)
		newCircle.Voters = circleVoters
	}

	if len(circleCreateRequest.Candidates) > 0 {
		if !newCircle.Private {
			err = fmt.Errorf("circle must be private to add candidates")
			return nil, err
		}

		circleCandidates := c.createCircleCandidateList(circleCreateRequest.Candidates)
		newCircle.Candidates = circleCandidates
	}

	if circleCreateRequest.Description != nil {
		newCircle.Description = strings.TrimSpace(*circleCreateRequest.Description)
	}

	currentTime := currentTruncatedTime()

	// check if new valid from time is given and is in the future from now on
	if circleCreateRequest.ValidFrom != nil {
		validFrom, err := extractValidFrom(currentTime, *circleCreateRequest.ValidFrom)

		if err != nil {
			return nil, err
		}

		newCircle.ValidFrom = *validFrom
	} else {
		newCircle.ValidFrom = currentTime
	}

	// check if new valid until time is given and is in the future from now on
	if circleCreateRequest.ValidUntil != nil {
		validUntil, err := extractValidUntil(currentTime, newCircle.ValidFrom, *circleCreateRequest.ValidUntil)

		if err != nil {
			return nil, err
		}

		newCircle.ValidUntil = validUntil
	}

	circle, err := c.storage.CreateNewCircle(newCircle)

	if err != nil {
		return nil, fmt.Errorf("error creating circle: %s", err)
	}

	return circle, nil
}

// DeleteCircle deletes a circle identified by the given circleId.
// It first validates the user's authentication claims to ensure they have the necessary permissions.
// If the user is not eligible to delete the circle, it returns an error with an appropriate message.
// If the circle is not active anymore, it returns an error indicating that the circle is not active.
// If there are no errors, it updates the circle's active field to false, indicating that it is deleted.
// It then updates the circle in the storage. If any error occurs during the update, it returns the error.
// If the update is successful, it returns nil.
func (c *circleService) DeleteCircle(
	ctx context.Context,
	circleId int64,
) error {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return err
	}

	userId := authClaims.Subject

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return err
	}

	// checks whether user is eligible to delete this circle
	if circle.CreatedFrom != userId {
		c.log.Infof("user is not eligible to delete circle: user %s, circle ID %d", userId, circle.ID)
		err = fmt.Errorf("user is not eligible to delete circle")
		return err
	}

	// if circle is not active anymore, it can't be updated
	if !circle.Active {
		c.log.Infof("user try to delete inactive circle: user %s, circle ID %d", userId, circle.ID)
		err = fmt.Errorf("circle is not active")
		return err
	}

	err = c.inactivateCircle(circle)

	if err != nil {
		c.log.Warnf("could not deactivate circle, error: circle ID %d, error %s", circle.ID, err)
		return fmt.Errorf("could not delete circle")
	}

	return nil
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

	isVoterInCircle, err := c.storage.IsVoterInCircle(authClaims.Subject, circle.ID)

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
	_, err = c.storage.CreateNewCircleVoter(circleVoter)

	if err != nil {
		c.log.Errorf("error adding voter to global circle: %s", err)
		return err
	}

	return nil
}

// determines, when the circle is private, if the user is eligible to
// interact with the circle
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

	ok, err := c.storage.IsVoterInCircle(userIdentityId, circle.ID)

	if err != nil && !database.RecordNotFound(err) {
		return false, err
	}

	if ok {
		return true, nil
	}

	return c.storage.IsCandidateInCircle(userIdentityId, circle.ID)
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

// based on the
// circleVoterInputs. It removes all the duplicates from the
// circleVoterInputs list.
func (c *circleService) createCircleVoterList(
	circleVoterInputs []*model.CircleVoterRequest,
) []*model.CircleVoter {
	var voterIdList []string

	for _, voter := range circleVoterInputs {
		voterIdList = append(voterIdList, voter.Voter)
	}

	voterIdList = utils.RemoveDuplicateStr(voterIdList)

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

// based on the
// circleCandidatesInputs. It removes all the duplicates from the
// circleCandidatesInputs list.
func (c *circleService) createCircleCandidateList(
	circleCandidatesInputs []*model.CircleCandidateRequest,
) []*model.CircleCandidate {
	var candidateIdList []string

	for _, candidate := range circleCandidatesInputs {
		candidateIdList = append(candidateIdList, candidate.Candidate)
	}

	candidateIdList = utils.RemoveDuplicateStr(candidateIdList)

	var circleCandidates []*model.CircleCandidate
	// add the given voters to the circle voters
	for _, candidate := range candidateIdList {
		circleCandidate := &model.CircleCandidate{
			Candidate: candidate,
		}
		circleCandidates = append(circleCandidates, circleCandidate)
	}

	return circleCandidates
}

func currentTruncatedTime() time.Time {
	return time.Now().UTC().Truncate(60 * time.Second)
}

// Function to check if a time is in the future from now
func isTimeInFuture(currentTime, futureTime time.Time) error {
	if currentTime.After(futureTime) {
		return fmt.Errorf("time must be in the future from now")
	}
	return nil
}

// Function to check if validUntil is after validFrom
func isValidUntilAfterValidFrom(validFrom, validUntil time.Time) error {
	if validUntil.Before(validFrom) {
		return fmt.Errorf("valid until time must not be before valid from time")
	}
	return nil
}

// Function to handle validFrom time validation
func extractValidFrom(
	currentTime time.Time,
	validFrom time.Time,
) (*time.Time, error) {
	validFromTime := validFrom.UTC().Truncate(60 * time.Second)
	if err := isTimeInFuture(currentTime, validFromTime); err != nil {
		return nil, err
	}
	return &validFromTime, nil
}

// Function to handle validUntil time validation
func extractValidUntil(
	currentTime time.Time,
	validFrom time.Time,
	validUntil time.Time,
) (*time.Time, error) {
	validUntilTime := validUntil.UTC().Truncate(60 * time.Second)
	if err := isTimeInFuture(currentTime, validUntilTime); err != nil {
		return nil, err
	}

	if err := isValidUntilAfterValidFrom(validFrom, validUntilTime); err != nil {
		return nil, err
	}

	return &validUntilTime, nil
}

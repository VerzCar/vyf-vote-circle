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

type VoteService interface {
	CreateVote(
		ctx context.Context,
		circleId int64,
		voteInput *model.VoteCreateInput,
	) (bool, error)
}

type VoteRepository interface {
	CircleById(id int64) (*model.Circle, error)
	CircleVoterByCircleId(circleId int64, voterId string) (*model.CircleVoter, error)
	CreateNewVote(
		voterId int64,
		electedId int64,
		circleId int64,
	) (*model.Vote, error)
	ElectedVoterCountsByCircleId(circleId int64, electedId int64) (int64, error)
	VoterElectedByCircleId(
		circleId int64,
		voterId int64,
		electedId int64,
	) (*model.Vote, error)
}

type VoteCache interface {
	UpdateRanking(
		ctx context.Context,
		circleId int64,
		identityId string,
		votes int64,
	) error
}

type VoteSubscription interface {
	RankingChangedEvent(circleId int64)
}

type voteService struct {
	storage      VoteRepository
	cache        VoteCache
	subscription VoteSubscription
	config       *config.Config
	log          logger.Logger
}

func NewVoteService(
	circleRepo VoteRepository,
	cache VoteCache,
	subscription VoteSubscription,
	config *config.Config,
	log logger.Logger,
) VoteService {
	return &voteService{
		storage:      circleRepo,
		cache:        cache,
		subscription: subscription,
		config:       config,
		log:          log,
	}
}

func (c *voteService) CreateVote(
	ctx context.Context,
	circleId int64,
	voteInput *model.VoteCreateInput,
) (bool, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return false, err
	}

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return false, err
	}

	if !circle.Active {
		c.log.Infof(
			"tried to vote for an inactive circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return false, fmt.Errorf("circle inactive")
	}

	voterId := authClaims.Subject

	voter, err := c.storage.CircleVoterByCircleId(circleId, voterId)

	if err != nil {
		c.log.Errorf("error voter id %s not in circle: %s", voterId, err)
		return false, err
	}

	elected, err := c.storage.CircleVoterByCircleId(circleId, voteInput.Elected)

	if err != nil {
		c.log.Errorf("error elected id %s not in circle: %s", voteInput.Elected, err)
		return false, err
	}

	// validate if voter already elected once - if so throw an error
	_, err = c.storage.VoterElectedByCircleId(circleId, voter.ID, elected.ID)

	if err != nil && !database.RecordNotFound(err) {
		c.log.Errorf("error getting voter %d for elected %d not in circle: %s", voter.ID, elected.ID, err)
		return false, err
	}
	if err == nil {
		c.log.Errorf(
			"failure voter %s for elected %s already voted in circle: %d",
			voter.Voter,
			elected.Voter,
			circleId,
		)
		return false, fmt.Errorf("already voted in circle")
	}

	// TODO put this write block in transaction as the update ranking in the cache could fail
	_, err = c.storage.CreateNewVote(voter.ID, elected.ID, circleId)

	if err != nil {
		return false, err
	}

	voteCount, err := c.storage.ElectedVoterCountsByCircleId(circleId, elected.ID)

	if err != nil {
		return false, err
	}

	err = c.cache.UpdateRanking(ctx, circleId, elected.Voter, voteCount)

	if err != nil {
		return false, err
	}

	c.subscription.RankingChangedEvent(circleId)

	return true, nil
}

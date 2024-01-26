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
		voteReq *model.VoteCreateRequest,
	) (bool, error)
}

type VoteRepository interface {
	CircleById(id int64) (*model.Circle, error)
	CircleVoterByCircleId(circleId int64, voterId string) (*model.CircleVoter, error)
	CircleCandidateByCircleId(
		circleId int64,
		candidateId string,
	) (*model.CircleCandidate, error)
	UpdateCircleVoter(voter *model.CircleVoter) (*model.CircleVoter, error)
	CreateNewVote(
		voterId int64,
		candidateId int64,
		circleId int64,
	) (*model.Vote, error)
	CountsVotesOfCandidateByCircleId(circleId int64, candidateId int64) (int64, error)
	VoterCandidateByCircleId(
		circleId int64,
		voterId int64,
		electedId int64,
	) (*model.Vote, error)
	CreateNewRanking(ranking *model.Ranking) (*model.Ranking, error)
	UpdateRanking(ranking *model.Ranking) (*model.Ranking, error)
	RankingByCircleId(circleId int64, identityId string) (*model.Ranking, error)
}

type VoteCache interface {
	UpsertRanking(
		ctx context.Context,
		circleId int64,
		candidate *model.CircleCandidate,
		ranking *model.Ranking,
		votes int64,
	) (*model.RankingResponse, error)
}

type VoteSubscription interface {
	RankingChangedEvent(
		ctx context.Context,
		circleId int64,
		ranking *model.RankingResponse,
	) error
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
	voteReq *model.VoteCreateRequest,
) (bool, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("error getting auth claims: %s", err)
		return false, err
	}

	voterId := authClaims.Subject

	if voterId == voteReq.CandidateID {
		c.log.Errorf("error voter id %s is equal candidate id: %s", voterId, voteReq.CandidateID)
		return false, fmt.Errorf("cannot vote for yourself")
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

	voter, err := c.storage.CircleVoterByCircleId(circleId, voterId)

	if err != nil {
		c.log.Errorf("error voter id %s not in circle: %s", voterId, err)
		return false, err
	}

	candidate, err := c.storage.CircleCandidateByCircleId(circleId, voteReq.CandidateID)

	if err != nil {
		c.log.Errorf("error candidate id %s not in circle: %s", voteReq.CandidateID, err)
		return false, err
	}

	if candidate.Commitment != model.CommitmentCommitted {
		c.log.Infof(
			"tried to vote for an uncommitted candidate with circle id %d and candidate id %d",
			circleId,
			candidate.ID,
		)
		return false, fmt.Errorf("candidate uncommitted")
	}

	// validate if voter already elected once - if so throw an error
	_, err = c.storage.VoterCandidateByCircleId(circleId, voter.ID, candidate.ID)

	if err != nil && !database.RecordNotFound(err) {
		c.log.Errorf("error getting voter %d for candidate %d not in circle: %s", voter.ID, candidate.ID, err)
		return false, err
	}
	if err == nil {
		c.log.Errorf(
			"voter %s for candidate %s already voted in circle: %d",
			voter.Voter,
			candidate.Candidate,
			circleId,
		)
		return false, fmt.Errorf("already voted in circle")
	}

	// TODO: put this write block in transaction as the update ranking in the cache could fail
	_, err = c.storage.CreateNewVote(voter.ID, candidate.ID, circleId)

	if err != nil {
		return false, err
	}

	voteCount, err := c.storage.CountsVotesOfCandidateByCircleId(circleId, candidate.ID)

	if err != nil {
		return false, err
	}

	// update the ranking meta information
	var ranking *model.Ranking

	ranking, err = c.storage.RankingByCircleId(circleId, candidate.Candidate)

	switch {
	case err != nil && !database.RecordNotFound(err):
		return false, err
	case database.RecordNotFound(err):
		newRanking := &model.Ranking{
			IdentityID: candidate.Candidate,
			Number:     0,
			Votes:      voteCount,
			CircleID:   circleId,
		}

		ranking, err = c.storage.CreateNewRanking(newRanking)

		if err != nil {
			return false, err
		}
		break
	default:
		ranking, err = c.storage.UpdateRanking(&model.Ranking{ID: ranking.ID, Votes: voteCount})

		if err != nil {
			return false, err
		}
	}

	// update the voters meta information
	voter.VotedFor = &candidate.Candidate
	_, err = c.storage.UpdateCircleVoter(voter)

	if err != nil {
		return false, err
	}

	updatedRanking, err := c.cache.UpsertRanking(ctx, circleId, candidate, ranking, voteCount)

	if err != nil {
		return false, err
	}

	_ = c.subscription.RankingChangedEvent(ctx, circleId, updatedRanking)

	return true, nil
}

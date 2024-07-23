package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/cache"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	routerContext "github.com/VerzCar/vyf-vote-circle/app/router/ctx"
	"time"
)

type VoteService interface {
	CreateVote(
		ctx context.Context,
		circleId int64,
		voteReq *model.VoteCreateRequest,
	) (bool, error)
	RevokeVote(
		ctx context.Context,
		circleId int64,
	) (bool, error)
}

type VoteRepository interface {
	CircleById(id int64) (*model.Circle, error)
	CircleVoterByCircleId(circleId int64, userIdentityId string) (*model.CircleVoter, error)
	CircleCandidateByCircleId(
		circleId int64,
		userIdentityId string,
	) (*model.CircleCandidate, error)
	CreateNewVote(
		ctx context.Context,
		circleId int64,
		voter *model.CircleVoter,
		candidate *model.CircleCandidate,
		upsertRankingCache cache.UpsertRankingCacheCallback,
	) (*model.RankingResponse, int64, error)
	VoteByCircleId(
		circleId int64,
		voterId int64,
	) (*model.Vote, error)
	DeleteVote(
		ctx context.Context,
		circleId int64,
		vote *model.Vote,
		voter *model.CircleVoter,
		upsertRankingCache cache.UpsertRankingCacheCallback,
		removeRankingCache cache.RemoveRankingCacheCallback,
	) (*model.RankingResponse, int64, error)
	HasVoterVotedForCircle(
		circleId int64,
		voterId int64,
	) (bool, error)
	UpdateRanking(ranking *model.Ranking) (*model.Ranking, error)
}

type VoteCache interface {
	UpsertRanking(
		ctx context.Context,
		circleId int64,
		candidate *model.CircleCandidate,
		ranking *model.Ranking,
		votes int64,
	) (*model.RankingResponse, error)
	RemoveRanking(
		ctx context.Context,
		circleId int64,
		candidate *model.CircleCandidate,
	) error
	RankingList(
		ctx context.Context,
		circleId int64,
		fromRanking *model.RankingResponse,
	) ([]*model.RankingResponse, error)
}

type VoteRankingSubscription interface {
	RankingChangedEvent(
		ctx context.Context,
		circleId int64,
		events []*model.RankingChangedEvent,
	) error
}

type VoteCircleVoterSubscription interface {
	CircleVoterChangedEvent(
		ctx context.Context,
		circleId int64,
		event *model.CircleVoterChangedEvent,
	) error
}

type VoteCircleCandidateSubscription interface {
	CircleCandidateChangedEvent(
		ctx context.Context,
		circleId int64,
		event *model.CircleCandidateChangedEvent,
	) error
}

type voteService struct {
	storage                     VoteRepository
	cache                       VoteCache
	rankingSubscription         VoteRankingSubscription
	circleVoterSubscription     VoteCircleVoterSubscription
	circleCandidateSubscription VoteCircleCandidateSubscription
	config                      *config.Config
	log                         logger.Logger
}

func NewVoteService(
	circleRepo VoteRepository,
	cache VoteCache,
	rankingSubscription VoteRankingSubscription,
	circleVoterSubscription VoteCircleVoterSubscription,
	circleCandidateSubscription VoteCircleCandidateSubscription,
	config *config.Config,
	log logger.Logger,
) VoteService {
	return &voteService{
		storage:                     circleRepo,
		cache:                       cache,
		rankingSubscription:         rankingSubscription,
		circleVoterSubscription:     circleVoterSubscription,
		circleCandidateSubscription: circleCandidateSubscription,
		config:                      config,
		log:                         log,
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

	if !circle.IsEditable() {
		c.log.Infof(
			"tried to vote for an ineditable circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return false, fmt.Errorf("circle is not editable")
	}

	if circle.Stage == model.CircleStageCold {
		c.log.Infof(
			"tried to vote for an cold circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return false, fmt.Errorf("circle is cold")
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
	hasVoted, err := c.storage.HasVoterVotedForCircle(circleId, voter.ID)

	if err != nil && !database.RecordNotFound(err) {
		return false, fmt.Errorf("already voted in circle")
	}
	if err == nil && hasVoted {
		c.log.Errorf(
			"voter %s for candidate %s already voted in circle: %d",
			voter.Voter,
			candidate.Candidate,
			circleId,
		)
		return false, fmt.Errorf("already voted in circle")
	}

	cachedRanking, voteCount, err := c.storage.CreateNewVote(ctx, circleId, voter, candidate, c.cache.UpsertRanking)

	if err != nil {
		return false, err
	}

	events := make([]*model.RankingChangedEvent, 0)

	if voteCount > 1 {
		event := CreateRankingChangedEvent(model.EventOperationUpdated, cachedRanking)
		events = append(events, event)
	} else {
		event := CreateRankingChangedEvent(model.EventOperationCreated, cachedRanking)
		events = append(events, event)
	}

	// TODO: update only if the number and index has not changed from the cachedRanking
	changedRankings, err := c.changedRankings(ctx, circleId, cachedRanking)

	if err != nil {
		return false, err
	}

	for _, changedRanking := range changedRankings {
		event := CreateRankingChangedEvent(model.EventOperationUpdated, changedRanking)
		events = append(events, event)
	}

	_ = c.rankingSubscription.RankingChangedEvent(ctx, circleId, events)

	voterEvent := CreateVoterChangedEvent(model.EventOperationUpdated, voter)
	_ = c.circleVoterSubscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)

	//TODO: do not only send events also update rankings table, or do it async
	// in the background from time to time, as votes are already persisted.

	return true, nil
}

func (c *voteService) RevokeVote(
	ctx context.Context,
	circleId int64,
) (bool, error) {
	authClaims, err := routerContext.ContextToAuthClaims(ctx)

	if err != nil {
		c.log.Errorf("getting auth claims: %s", err)
		return false, err
	}

	circle, err := c.storage.CircleById(circleId)

	if err != nil {
		return false, err
	}

	if !circle.Active {
		c.log.Infof(
			"tried to revoke vote for an inactive circle with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return false, fmt.Errorf("circle inactive")
	}

	currentTime := currentTruncatedTime()

	if circle.ValidUntil != nil && currentTime.After(circle.ValidUntil.UTC().Truncate(60*time.Second)) {
		c.log.Infof(
			"tried to revoke vote for an circle that is not valid anymore with circle id %d and subject %s",
			circleId,
			authClaims.Subject,
		)
		return false, fmt.Errorf("circle closed")
	}

	voter, err := c.storage.CircleVoterByCircleId(circleId, authClaims.Subject)

	if err != nil {
		c.log.Errorf("voter userIdentity id %s not in circle: %s", authClaims.Subject, err)
		return false, err
	}

	vote, err := c.storage.VoteByCircleId(circleId, voter.ID)

	if err != nil && !database.RecordNotFound(err) {
		c.log.Errorf("getting vote for voter %d for circle id %d: %s", voter.ID, circleId, err)
		return false, err
	}

	if database.RecordNotFound(err) {
		c.log.Errorf("user has not voted for circle id %d", circleId)
		return false, fmt.Errorf("no voting exists")
	}

	cachedRanking, voteCount, err := c.storage.DeleteVote(
		ctx,
		circleId,
		vote,
		voter,
		c.cache.UpsertRanking,
		c.cache.RemoveRanking,
	)

	if err != nil {
		return false, err
	}

	if voteCount > 0 {
		events := make([]*model.RankingChangedEvent, 0)

		event := CreateRankingChangedEvent(model.EventOperationUpdated, cachedRanking)
		events = append(events, event)

		// TODO: update only if the number and index has not changed from the cachedRanking
		changedRankings, err := c.changedRankings(ctx, circleId, cachedRanking)

		if err != nil {
			return false, err
		}

		for _, changedRanking := range changedRankings {
			event := CreateRankingChangedEvent(model.EventOperationUpdated, changedRanking)
			events = append(events, event)
		}

		_ = c.rankingSubscription.RankingChangedEvent(ctx, circleId, events)

		voterEvent := CreateVoterChangedEvent(model.EventOperationUpdated, voter)
		_ = c.circleVoterSubscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)

		return true, nil
	}

	events := make([]*model.RankingChangedEvent, 0)

	event := CreateRankingChangedEvent(model.EventOperationDeleted, cachedRanking)
	events = append(events, event)

	// TODO: update only if the number and index has not changed from the cachedRanking
	changedRankings, err := c.changedRankings(ctx, circleId, nil)

	if err != nil {
		return false, err
	}

	for _, changedRanking := range changedRankings {
		event := CreateRankingChangedEvent(model.EventOperationUpdated, changedRanking)
		events = append(events, event)
	}

	_ = c.rankingSubscription.RankingChangedEvent(ctx, circleId, events)

	voterEvent := CreateVoterChangedEvent(model.EventOperationUpdated, voter)
	_ = c.circleVoterSubscription.CircleVoterChangedEvent(ctx, circleId, voterEvent)

	candidateEvent := CreateCandidateChangedEvent(model.EventOperationRepositioned, vote.Candidate)
	_ = c.circleCandidateSubscription.CircleCandidateChangedEvent(ctx, circleId, candidateEvent)

	return true, nil
}

func (c *voteService) changedRankings(
	ctx context.Context,
	circleId int64,
	updatedRanking *model.RankingResponse,
) ([]*model.RankingResponse, error) {
	return c.cache.RankingList(ctx, circleId, updatedRanking)
}

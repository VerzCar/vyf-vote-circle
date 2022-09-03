package api

import (
	"context"
	"github.com/google/uuid"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	"sync"
)

type RankingSubscriptionService interface {
	Rankings(
		ctx context.Context,
		circleId int64,
	) (<-chan []*model.Ranking, error)
	RankingChangedEvent(circleId int64)
}

type rankingSubscriptionService struct {
	storage          RankingRepository
	cache            RankingCache
	rankingObservers model.RankingObservers
	mu               sync.Mutex
	rankingService   RankingService
	config           *config.Config
	log              logger.Logger
}

func NewRankingSubscriptionService(
	circleRepo RankingRepository,
	cache RankingCache,
	rankingService RankingService,
	config *config.Config,
	log logger.Logger,
) RankingSubscriptionService {
	rankingObservers := make(model.RankingObservers)

	return &rankingSubscriptionService{
		storage:          circleRepo,
		cache:            cache,
		rankingObservers: rankingObservers,
		rankingService:   rankingService,
		config:           config,
		log:              log,
	}
}

// Rankings as a channel with
func (s *rankingSubscriptionService) Rankings(
	ctx context.Context,
	circleId int64,
) (<-chan []*model.Ranking, error) {
	s.initRankingListObservableMap(circleId)
	observerId := uuid.New().String()
	rankings := make(model.RankingListObservable, 1)

	// Start a goroutine to allow for cleaning up subscriptions that are disconnected.
	// This go routine will only get past Done() when a client terminates the subscription. This allows us
	// to only then remove the reference from the list of ChatObservers since it is no longer needed.
	go func() {
		<-ctx.Done()
		s.mu.Lock()
		delete(s.rankingObservers[circleId], observerId)
		s.mu.Unlock()
	}()

	s.mu.Lock()
	// Keep a reference of the channel so that we can push changes into it when new messages are posted.
	s.rankingObservers[circleId][observerId] = rankings
	s.mu.Unlock()
	// This is optional, and this allows newly subscribed clients to get a list of all the rankings that have been
	// listed so far. Upon subscribing the client will be pushed the rankings once, further changes are handled
	// in the CreateVote mutation.
	rankingList, err := s.rankingService.Rankings(ctx, circleId)

	if err == nil {
		s.rankingObservers[circleId][observerId] <- rankingList
	}

	return rankings, nil
}

// RankingChangedEvent update the changed rankings for the circle.
// Pushes the new list to the observables.
func (s *rankingSubscriptionService) RankingChangedEvent(
	circleId int64,
) {
	ctx := context.Background()
	rankingList, err := s.rankingService.Rankings(ctx, circleId)

	if err == nil {
		s.mu.Lock()
		for _, observer := range s.rankingObservers[circleId] {
			observer <- rankingList
		}
		s.mu.Unlock()
	}
}

func (s *rankingSubscriptionService) initRankingListObservableMap(
	circleId int64,
) {
	if s.rankingObservers[circleId] == nil {
		s.rankingObservers[circleId] = make(model.RankingListObservableMap)
	}
}

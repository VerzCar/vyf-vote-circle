package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/ably/ably-go/ably"
)

type CircleVoterSubscriptionService interface {
	CircleVoterChangedEvent(
		ctx context.Context,
		circleId int64,
		event *model.CircleVoterChangedEvent,
	) error
}

type circleVoterSubscriptionService struct {
	pubSubService *ably.Realtime
	log           logger.Logger
}

func NewCircleVoterSubscriptionService(
	pubSubService *ably.Realtime,
	log logger.Logger,
) CircleVoterSubscriptionService {
	return &circleVoterSubscriptionService{
		pubSubService: pubSubService,
		log:           log,
	}
}

func (s *circleVoterSubscriptionService) CircleVoterChangedEvent(
	ctx context.Context,
	circleId int64,
	event *model.CircleVoterChangedEvent,
) error {
	channelName := fmt.Sprintf("circle-%d:voter", circleId)
	msgName := "circle-voter-changed"

	channel := s.pubSubService.Channels.Get(channelName)

	err := channel.Publish(ctx, msgName, event)

	if err != nil {
		s.log.Errorf(
			"could not publish message to channel: %s with message name: %s cause: %s",
			channelName,
			msgName,
			err,
		)
		return err
	}

	return nil
}

func CreateVoterChangedEvent(
	operation model.EventOperation,
	voter *model.CircleVoter,
) *model.CircleVoterChangedEvent {
	return &model.CircleVoterChangedEvent{
		Operation: operation,
		Voter: &model.CircleVoterResponse{
			ID:         voter.ID,
			Voter:      voter.Voter,
			VotedFor:   voter.VotedFor,
			Commitment: voter.Commitment,
			CreatedAt:  voter.CreatedAt,
			UpdatedAt:  voter.UpdatedAt,
		},
	}
}

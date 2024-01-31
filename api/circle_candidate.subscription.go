package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/ably/ably-go/ably"
)

type CircleCandidateSubscriptionService interface {
	CircleCandidateChangedEvent(
		ctx context.Context,
		circleId int64,
		event *model.CircleCandidateChangedEvent,
	) error
}

type circleCandidateSubscriptionService struct {
	pubSubService *ably.Realtime
	log           logger.Logger
}

func NewCircleCandidateSubscriptionService(
	pubSubService *ably.Realtime,
	log logger.Logger,
) CircleCandidateSubscriptionService {
	return &circleCandidateSubscriptionService{
		pubSubService: pubSubService,
		log:           log,
	}
}

func (s *circleCandidateSubscriptionService) CircleCandidateChangedEvent(
	ctx context.Context,
	circleId int64,
	event *model.CircleCandidateChangedEvent,
) error {
	channelName := fmt.Sprintf("circle-%d:candidate", circleId)
	msgName := "circle-candidate-changed"

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

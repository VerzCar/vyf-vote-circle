package api

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/ably/ably-go/ably"
)

type RankingSubscriptionService interface {
	RankingChangedEvent(
		ctx context.Context,
		circleId int64,
		events []*model.RankingChangedEvent,
	) error
}

type rankingSubscriptionService struct {
	pubSubService *ably.Realtime
	log           logger.Logger
}

func NewRankingSubscriptionService(
	pubSubService *ably.Realtime,
	log logger.Logger,
) RankingSubscriptionService {
	return &rankingSubscriptionService{
		pubSubService: pubSubService,
		log:           log,
	}
}

// Will notify all clients of changed ranking of certain circle.
// It will open the channel for the circle and send the updated ranking as
// message to all subscribed clients.
func (s *rankingSubscriptionService) RankingChangedEvent(
	ctx context.Context,
	circleId int64,
	events []*model.RankingChangedEvent,
) error {
	channelName := fmt.Sprintf("circle-%d:rankings", circleId)
	msgName := "ranking-changed"

	messages := make([]*ably.Message, 0)

	for _, event := range events {
		message := &ably.Message{
			Name: msgName,
			Data: event,
		}
		messages = append(messages, message)
	}

	channel := s.pubSubService.Channels.Get(channelName)
	err := channel.PublishMultiple(ctx, messages)

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

func CreateRankingChangedEvent(
	operation model.EventOperation,
	ranking *model.RankingResponse,
) *model.RankingChangedEvent {
	return &model.RankingChangedEvent{
		Operation: operation,
		Ranking:   ranking,
	}
}

package pubsub

import (
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/ably/ably-go/ably"
)

// Connect the pub sub service
func Connect(log logger.Logger, conf *config.Config) *ably.Realtime {
	log.Infof("Connect to ably service")

	client, err := ably.NewRealtime(
		ably.WithKey(conf.Ably.Apikey),
		ably.WithClientID(conf.Ably.ClientId),
	)

	if err != nil {
		log.Fatalf("Connect to ably service failed. cause: %s", err)
	}

	return client
}

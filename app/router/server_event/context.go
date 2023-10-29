package server_event

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func ContextToServerEventChan[T any](ctx *gin.Context) (ServerEventChan[T], error) {
	v, ok := ctx.Get(serverEventChanContextKey)
	if !ok {
		err := fmt.Errorf("could not retrieve server event channel")
		return nil, err
	}

	serverEventChan, ok := v.(ServerEventChan[T])

	if !ok {
		err := fmt.Errorf("could not retrieve server event channel")
		return nil, err
	}

	return serverEventChan, nil
}

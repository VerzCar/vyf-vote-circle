package server_event

import "github.com/gin-gonic/gin"

type ServerEventService[T any] interface {
	ServeHTTP() gin.HandlerFunc
	Publish(msg T)
}

func NewServerEventService[T any]() ServerEventService[T] {
	event := &ServerEvent[T]{
		Message:       make(chan T),
		NewClients:    make(chan chan T),
		ClosedClients: make(chan chan T),
		TotalClients:  make(map[chan T]bool),
	}

	go event.listen()

	return event
}

package app

import (
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/gin-gonic/gin"
	"log"
)

type ServerEventService interface {
	ServeHTTP() gin.HandlerFunc
	MessageSub(sub <-chan []*model.Ranking)
}

// ServerEvent keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type ServerEvent struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message <-chan []*model.Ranking

	// New client connections
	NewClients chan chan []*model.Ranking

	// Closed client connections
	ClosedClients chan chan []*model.Ranking

	// Total client connections
	TotalClients map[chan []*model.Ranking]bool
}

// New event messages are broadcast to all registered client connection channels
type ClientChan chan []*model.Ranking

func NewServerEventService() ServerEventService {
	event := &ServerEvent{
		Message:       make(chan []*model.Ranking),
		NewClients:    make(chan chan []*model.Ranking),
		ClosedClients: make(chan chan []*model.Ranking),
		TotalClients:  make(map[chan []*model.Ranking]bool),
	}

	go event.listen()

	return event
}

func (stream *ServerEvent) ServeHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Initialize client channel
		clientChan := make(ClientChan)

		// Send new connection to event server
		stream.NewClients <- clientChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- clientChan
		}()

		c.Set("clientChan", clientChan)

		c.Next()
	}
}

func (stream *ServerEvent) MessageSub(sub <-chan []*model.Ranking) {
	stream.Message = sub
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *ServerEvent) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true
			log.Printf("Client added. %d registered clients", len(stream.TotalClients))

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)
			log.Printf("Removed client. %d registered clients", len(stream.TotalClients))

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

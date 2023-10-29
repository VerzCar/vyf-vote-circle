package server_event

// ServerEvent keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type ServerEvent[T any] struct {
	// Events are pushed to this channel by the main events-gathering routine
	Message chan T

	// New client connections
	NewClients chan chan T

	// Closed client connections
	ClosedClients chan chan T

	// Total client connections
	TotalClients map[chan T]bool
}

// New event messages are broadcast to all registered client connection channels
type ServerEventChan[T any] chan T

// Publish message to clients.
func (stream *ServerEvent[T]) Publish(msg T) {
	stream.Message <- msg
}

// listen to all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (stream *ServerEvent[T]) listen() {
	for {
		select {
		// Add new available client
		case client := <-stream.NewClients:
			stream.TotalClients[client] = true

		// Remove closed client
		case client := <-stream.ClosedClients:
			delete(stream.TotalClients, client)
			close(client)

		// Broadcast message to client
		case eventMsg := <-stream.Message:
			for clientMessageChan := range stream.TotalClients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

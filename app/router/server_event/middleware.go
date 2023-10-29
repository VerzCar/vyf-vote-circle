package server_event

import "github.com/gin-gonic/gin"

const serverEventChanContextKey = "ServerEventChanContextKey"

func (stream *ServerEvent[T]) ServeHTTP() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		serverEventChan := make(ServerEventChan[T])

		// Send new connection to event server
		stream.NewClients <- serverEventChan

		defer func() {
			// Send closed connection to event server
			stream.ClosedClients <- serverEventChan
		}()

		ctx.Set(serverEventChanContextKey, serverEventChan)

		ctx.Next()
	}
}

package app

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/graph/generated"
)

// GQL defines the Graphql handler
func gqlHandler(resolver *Resolver) gin.HandlerFunc {
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

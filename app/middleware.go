package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gitlab.vecomentman.com/libs/awsx"
	routerContext "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router/ctx"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router/header"
	"net/http"
)

// ginContextToContext creates a gin middleware to add its context
// to the context.Context
func (s *Server) ginContextToContext() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		routerContext.SetGinContext(ctx)
		ctx.Next()
	}
}

// authGuard verifies the Authorization token against the SSO service.
// If the authentication fails the request will be aborted.
// Otherwise, the given subject of the token will be saved in the context and
// the next request served.
func (s *Server) authGuard(authService awsx.AuthService) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		accessToken, err := header.Authorization(ctx, "Bearer")

		if err != nil {
			ctx.String(http.StatusUnauthorized, fmt.Sprintf("error: %s", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token, err := authService.DecodeAccessToken(ctx, accessToken)

		if err != nil {
			ctx.String(http.StatusUnauthorized, fmt.Sprintf("error decoding token"))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		routerContext.SetAuthClaimsContext(ctx, token)
		ctx.Next()
	}
}

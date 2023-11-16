package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (s *Server) TokenAbly() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenRequest, err := s.tokenService.GenerateAblyToken(ctx.Request.Context())

		if err != nil {
			s.log.Errorf("service error: %v", err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, tokenRequest)
	}
}

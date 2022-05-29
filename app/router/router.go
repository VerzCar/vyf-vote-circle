package router

import (
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
)

// Setup the router
func Setup(environment string) *gin.Engine {
	if environment == config.EnvironmentProd {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()
	router := gin.Default()

	addMiddleware(router)

	return router
}

// AddMiddleware adds the available middlewares to the router
func addMiddleware(r *gin.Engine) {
	r.Use(cors.Default())
}

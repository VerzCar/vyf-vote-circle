package router

import (
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
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

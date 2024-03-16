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
	corsOptions := cors.Options{
		AllowedOrigins: []string{"https://vyf-web-app-c3f1d65ba31f.herokuapp.com"},
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://vyf-web-app-c3f1d65ba31f.herokuapp.com"
		},
		AllowedMethods:   []string{"GET", "POST", "OPTION", "DELETE", "PUT", "PATCH"},
		AllowedHeaders:   []string{"Origin"},
		ExposedHeaders:   []string{"Content-Length"},
		MaxAge:           10800, // 3 hours
		AllowCredentials: true,
	}
	r.Use(cors.New(corsOptions))
}

package router

import (
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"net/http"
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
		AllowedOrigins: []string{
			"https://vyf-web-app-c3f1d65ba31f.herokuapp.com",
			"https://vyf-web-app-c3f1d65ba31f.herokuapp.com/",
		},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"Origin"},
		ExposedHeaders:   []string{"Content-Length"},
		MaxAge:           10800, // 3 hours
		AllowCredentials: true,
	}
	r.Use(cors.New(corsOptions))
}

package router

import (
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"net/http"
)

// Setup the router
func Setup(envConfig *config.Config) *gin.Engine {
	if envConfig.Environment == config.EnvironmentProd {
		gin.SetMode(gin.ReleaseMode)
	}
	gin.ForceConsoleColor()
	router := gin.Default()

	addMiddleware(router, envConfig)

	return router
}

// AddMiddleware adds the available middlewares to the router
func addMiddleware(r *gin.Engine, envConfig *config.Config) {
	corsOptions := cors.Options{
		AllowedOrigins: envConfig.Security.Cors.Origins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Accept-Encoding",
			"X-Requested-With",
			"Content-Type",
			"Accept",
			"x-client-key",
			"x-client-token",
			"x-client-secret",
			"X-CSRF-Token",
			"Cache-Control",
			"Authorization",
		},
		ExposedHeaders:   []string{"Content-Length"},
		MaxAge:           10800, // 3 hours
		AllowCredentials: true,
	}
	r.Use(cors.New(corsOptions))
}

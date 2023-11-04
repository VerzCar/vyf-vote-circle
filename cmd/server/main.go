package main

import (
	"fmt"
	"github.com/VerzCar/vyf-lib-awsx"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api"
	"github.com/VerzCar/vyf-vote-circle/app"
	"github.com/VerzCar/vyf-vote-circle/app/cache"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/app/database"
	"github.com/VerzCar/vyf-vote-circle/app/pubsub"
	"github.com/VerzCar/vyf-vote-circle/app/router"
	"github.com/VerzCar/vyf-vote-circle/repository"
	"github.com/VerzCar/vyf-vote-circle/utils"
	"github.com/go-playground/validator/v10"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "startup error: %s\\n", err)
		os.Exit(1)
	}
}

var validate *validator.Validate

func run() error {
	configPath := utils.FromBase("app/config/")

	envConfig := config.NewConfig(configPath)
	log := logger.NewLogger(configPath)

	log.Infof("Startup service...")
	log.Infof("Configuration loaded.")

	db := database.Connect(log, envConfig)

	storage := repository.NewStorage(db, envConfig, log)

	sqlDb, _ := db.DB()
	err := storage.RunMigrationsUp(sqlDb)

	if err != nil {
		return err
	}

	redisCache := cache.Connect(log, envConfig)
	redis := cache.NewRedisCache(redisCache, envConfig, log)

	// initialize auth service
	authService, err := awsx.NewAuthService(
		awsx.AppClientId(envConfig.Aws.Auth.ClientId),
		awsx.ClientSecret(envConfig.Aws.Auth.ClientSecret),
		awsx.AwsDefaultRegion(envConfig.Aws.Auth.AwsDefaultRegion),
		awsx.UserPoolId(envConfig.Aws.Auth.UserPoolId),
	)

	if err != nil {
		return err
	}

	// initialize pub sub service
	pubSubService := pubsub.Connect(log, envConfig)

	// initialize api services
	circleService := api.NewCircleService(storage, envConfig, log)
	rankingService := api.NewRankingService(storage, redis, envConfig, log)
	rankingSubscriptionService := api.NewRankingSubscriptionService(
		storage,
		redis,
		rankingService,
		pubSubService,
		envConfig,
		log,
	)
	voteService := api.NewVoteService(storage, redis, rankingSubscriptionService, envConfig, log)

	validate = validator.New()

	r := router.Setup(envConfig.Environment)
	server := app.NewServer(
		r,
		authService,
		circleService,
		rankingService,
		rankingSubscriptionService,
		voteService,
		validate,
		envConfig,
		log,
	)

	err = server.Run()

	if err != nil {
		return err
	}

	return nil
}

package cache

import (
	"context"
	"fmt"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/VerzCar/vyf-vote-circle/utils"
	"github.com/go-redis/redis/v8"
)

// Connect the cache database
func Connect(log logger.Logger, conf *config.Config) *redis.Client {
	opt, err := redis.ParseURL(redisUrl(conf))

	if err != nil {
		log.Fatalf("Options parsing for cache failed. cause: %s", err)
	}

	opt.ReadTimeout = utils.FormatDuration(uint(conf.Redis.Timeout))

	rdb := redis.NewClient(opt)

	ctx := context.Background()

	val, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("Connect to redis failed. cause: %s", err)
	}

	if val != "PONG" {
		log.Fatalf("Connect to redis failed. cause: ping failed.")
	}

	log.Infof("Connected successfully to redis via: %s", redisUrl(conf))

	return rdb
}

func redisUrl(conf *config.Config) string {
	redisConn := "redis"
	if conf.Environment == config.EnvironmentProd {
		redisConn = "rediss"
	}

	return fmt.Sprintf(
		"%s://%s:%s@%s:%d/%d",
		redisConn,
		conf.Redis.Username,
		conf.Redis.Password,
		conf.Redis.Host,
		conf.Redis.Port,
		conf.Redis.Db,
	)
}

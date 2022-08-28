package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/config"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils"
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
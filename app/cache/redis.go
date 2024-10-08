package cache

import (
	"context"
	"encoding/json"
	"errors"
	logger "github.com/VerzCar/vyf-lib-logger"
	"github.com/VerzCar/vyf-vote-circle/api/model"
	"github.com/VerzCar/vyf-vote-circle/app/config"
	"github.com/go-redis/redis/v8"
	"time"
)

type UpsertRankingCacheCallback func(
	context.Context,
	int64,
	*model.CircleCandidate,
	*model.Ranking,
	int64,
) (*model.RankingResponse, error)

type RemoveRankingCacheCallback func(
	context.Context,
	int64,
	*model.CircleCandidate,
) error

type RedisCache interface {
	UpsertRanking(
		ctx context.Context,
		circleId int64,
		candidate *model.CircleCandidate,
		ranking *model.Ranking,
		votes int64,
	) (*model.RankingResponse, error)
	RemoveRanking(
		ctx context.Context,
		circleId int64,
		candidate *model.CircleCandidate,
	) error
	RankingList(
		ctx context.Context,
		circleId int64,
		fromRanking *model.RankingResponse,
	) ([]*model.RankingResponse, error)
	ExistsRankingListForCircle(
		ctx context.Context,
		circleId int64,
	) (bool, error)
	BuildRankingList(
		ctx context.Context,
		circleId int64,
		rankingCacheItems []*model.RankingCacheItem,
	) error
}

type redisCache struct {
	redis  Client
	config *config.Config
	log    logger.Logger
}

func NewRedisCache(
	redis Client,
	config *config.Config,
	log logger.Logger,
) RedisCache {
	return &redisCache{
		redis:  redis,
		config: config,
		log:    log,
	}
}

type Entry struct {
	Val    string
	Exists bool
}

type EntryNumber struct {
	Val    int64
	Exists bool
}

type HEntry struct {
	Exists bool
}

// setJson converts the given value as JSON into the cache
// with the given key.
func (c *redisCache) setJson(ctx context.Context, key string, value interface{}, t time.Duration) error {
	encodedData, err := json.Marshal(value)

	if err != nil {
		return err
	}

	return c.redis.Set(ctx, key, encodedData, t).Err()
}

// getJson gets the entry from the given key
// as JSON format and Unmarshal it to the given destination
// interface structure type
func (c *redisCache) getJson(ctx context.Context, key string, dest interface{}) error {
	entry := c.redis.Get(ctx, key)

	switch {
	case errors.Is(entry.Err(), redis.Nil):
		return entry.Err()
	case entry.Err() != nil:
		return entry.Err()
	default:
		err := json.Unmarshal([]byte(entry.Val()), dest)
		return err
	}
}

// get the entry from the given key
// as Entry. If the entry does not exist, the Entry
// has the flag Exists set false and no error will be returned.
// Or if an error happens, an error will be returned.
// Otherwise, the Value and the Exists flag will be returned, without
// any error.
func (c *redisCache) get(ctx context.Context, key string) (Entry, error) {
	entry := c.redis.Get(ctx, key)

	sReturn := Entry{Exists: false}

	switch {
	case errors.Is(entry.Err(), redis.Nil):
		return sReturn, nil
	case entry.Err() != nil:
		return sReturn, entry.Err()
	default:
		sReturn.Exists = true
		sReturn.Val = entry.Val()
		return sReturn, nil
	}
}

func (c *redisCache) set(ctx context.Context, key string, value interface{}, t time.Duration) error {
	err := c.redis.Set(ctx, key, value, t).Err()
	return err
}

// FlushAll the cache and flush the db
func (c *redisCache) FlushAll() error {
	ctx := context.Background()
	return c.redis.FlushDB(ctx).Err()
}

package cache

import (
	"testing"
)

/*type cacheTestData struct {
	Name  string
	Email string
	Id    int
}

type cacheTestEntry struct {
	key      string
	data     cacheTestData
	duration time.Duration
}

var (
	rdb              *redis.Client
	log              logger.Logger
	conf             *config.Config
	ctx              context.Context
	defaultCacheData *cacheTestEntry
)*/

/*func TestMain(m *testing.M) {
	configPath := utils.FromBase("config/")

	conf = config.Load(configPath)
	log = logger.NewLogger(configPath)

	rdb = Connect(log, conf)

	ctx = context.Background()

	// set default values for test data
	defaultCacheData = &cacheTestEntry{
		key: "redis_test_key#1",
		data: cacheTestData{
			Name:  "TestName",
			Email: "example.me@mail.com",
			Id:    12,
		},
		duration: utils.FormatDuration(120),
	}

	os.Exit(m.Run())
}*/

func TestConnect(t *testing.T) {
	t.Skip("Test needs to be implemented")
}

/*func TestToJSON(t *testing.T) {
	assert := assertPkg.New(t)

	err := ToJSON(&ctx, rdb, defaultCacheData.key, defaultCacheData.data, defaultCacheData.duration)

	if err != nil {
		assert.Failf("to json failed: %s", err.Error())
	}

	expectedEntryData := cacheTestData{}
	entry := rdb.Get(ctx, defaultCacheData.key)

	switch {
	case entry.Err() == redis.Nil:
		assert.Failf("key does not exist in cache: %s", entry.Err().Error())
	case entry.Err() != nil:
		assert.Failf("error reading from cache: %s", entry.Err().Error())
	default:
		err := json.Unmarshal([]byte(entry.Val()), &expectedEntryData)
		if err != nil {
			assert.Failf("error converting entry: %s", err.Error())
		}
	}

	assert.Equal(expectedEntryData, defaultCacheData.data)

}

func TestJSONTo(t *testing.T) {
	assert := assertPkg.New(t)

	// set the default test data first
	err := ToJSON(&ctx, rdb, defaultCacheData.key, defaultCacheData.data, defaultCacheData.duration)

	if err != nil {
		assert.Failf("to json failed: %s", err.Error())
	}

	expectedEntryData := &cacheTestData{}

	err = JSONTo(&ctx, rdb, defaultCacheData.key, expectedEntryData)

	if err != nil {
		assert.Failf("json to struct failed: %s", err.Error())
	}

	assert.Equal(expectedEntryData, &defaultCacheData.data)
}

func TestGet(t *testing.T) {
	assert := assertPkg.New(t)

	// set the default test data first
	err := rdb.Set(ctx, defaultCacheData.key, defaultCacheData.data.Email, defaultCacheData.duration).Err()

	if err != nil {
		assert.Failf("set default value failed: %s", err.Error())
	}

	entry, err := Get(&ctx, rdb, defaultCacheData.key)

	assert.True(entry.Exists)
	assert.Equal(entry.Val, defaultCacheData.data.Email)

}

func TestUniqueKey(t *testing.T) {
	assert := assertPkg.New(t)

	uniqueKey := UniqueKey()

	assert.NotEmpty(uniqueKey)
}
*/

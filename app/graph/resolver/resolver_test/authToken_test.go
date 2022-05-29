package resolver_test

/*
import (
	assertPkg "github.com/stretchr/testify/assert"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/cache/testcache"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database/testdb"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils/testing/client"
	"testing"
	"time"
)

const (
	authToken = `mutation authToken($credentials: Credentials!) {
		authToken(credentials: $credentials) {
			token
			type
			expiresAt
		}
		}`
)

func TestAuthToken(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    factory.Martin.Email,
				Password: factory.Martin.Password,
			}))

	assert.NoError(err)

	assert.Equal(resp.AuthToken.Type, resolver.Config.Token.Type)
	assert.NotEqual(resp.AuthToken.Token, "")
	expectedTime := time.Now()
	timeDelta := utils.FormatDuration(resolver.Config.Ttl.Token.Default)
	expectedTime.Add(timeDelta)
	assert.Greater(resp.AuthToken.ExpiresAt.String(), expectedTime.String())

}

func TestAuthToken_Assert_Wrong_Email(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    "not.exist@email.com",
				Password: factory.Martin.Password,
			}))

	assert.Equal(err.Error(), `[{"message":"authentication failed","path":["authToken"]}]`)

}

func TestAuthToken_Assert_Wrong_Pwd(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    factory.Martin.Email,
				Password: "wrongPassword123",
			}))

	assert.Equal(err.Error(), `[{"message":"authentication failed","path":["authToken"]}]`)

}

func TestAuthToken_Assert_Inactive_User(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&resp,
		client.Var("credentials",
			model.Credentials{
				Email:    factory.Albert.Email,
				Password: factory.Albert.Password,
			}))

	assert.Equal(err.Error(), `[{"message":"authentication failed","path":["authToken"]}]`)

}
*/

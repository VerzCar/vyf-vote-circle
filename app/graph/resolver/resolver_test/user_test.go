package resolver_test

/*import (
	assertPkg "github.com/stretchr/testify/assert"
	model2 "gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/cache/testcache"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database/testdb"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils/testing/client"
	"testing"
)

const (
	user = `query user {
		  user {
			id
			email
		  }
		}`
)

func TestUser(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var respToken struct {
		AuthToken model2.Token
	}

	err := c.Post(
		authToken,
		&respToken,
		client.Var(
			"credentials",
			model2.Credentials{
				Email:    factory.Martin.Email,
				Password: factory.Martin.Password,
			},
		),
	)

	assert.NoError(err)

	var resp struct {
		User model2.User
	}

	err = c.Post(
		user,
		&resp,
		client.AddHeader("Authorization", resolver.Config.Token.Type+" "+respToken.AuthToken.Token),
	)

	assert.NoError(err)

	assert.Equal(resp.User.ID, factory.Martin.ID)
	assert.Equal(resp.User.Email, factory.Martin.Email)

}

func TestUser_Assert_Anonymous_User(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		User model2.User
	}

	err := c.Post(
		user,
		&resp,
		client.AddHeader(
			"Authorization",
			resolver.Config.Token.Type+" eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIxNDZjODJjYy1mYThmLTRkMjktYjAxMi1lOTRiYWEyZjAxMTUiLCJleHAiOjE2sdhjwewedyM30.f1Ilvh-eFYqnJiCVe09YYoXgwHtcjzNjgRmQc2lPc0M",
		),
	)

	assert.Equal(err.Error(), `[{"message":"authentication failed","path":["user"]}]`)

}
*/

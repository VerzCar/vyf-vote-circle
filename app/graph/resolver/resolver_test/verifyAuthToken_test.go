package resolver_test

/*import (
	assertPkg "github.com/stretchr/testify/assert"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/cache/testcache"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database/testdb"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils/testing/client"
	"testing"
)

const (
	verifyAuthToken = `query verifyAuthToken($authToken: String!) {
		verifyAuthToken(authToken: $authToken)
		}`
)

func TestVerifyAuthToken(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	// get valid auth token first
	var respUser struct {
		AuthToken model.Token
	}

	err := c.Post(
		authToken,
		&respUser,
		client.Var("credentials",
			model.Credentials{
				Email:    factory.Martin.Email,
				Password: factory.Martin.Password,
			}))

	assert.NoError(err)

	var resp struct {
		VerifyAuthToken bool
	}

	err = c.Post(
		verifyAuthToken,
		&resp,
		client.Var("authToken",
			respUser.AuthToken.Token))

	assert.NoError(err)

	assert.Equal(resp.VerifyAuthToken, true)

}

func TestVerifyAuthToken_Assert_Invalid_AuthToken(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		VerifyAuthToken bool
	}

	err := c.Post(
		verifyAuthToken,
		&resp,
		client.Var("authToken",
			"invalidAuthToken23324#$23"))

	assert.Equal(err.Error(), `[{"message":"verification failed","path":["verifyAuthToken"]}]`)
	assert.Equal(resp.VerifyAuthToken, false)
}
*/

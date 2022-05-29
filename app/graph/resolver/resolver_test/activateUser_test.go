package resolver_test

/*
import (
	assertPkg "github.com/stretchr/testify/assert"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/cache/testcache"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database/testdb"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils/testing/client"
	"testing"
)

const (
	activateUser = `mutation activateUser($payload: UserActivate!) {
  activateUser(payload: $payload) {
    email
    firstName
  }
}`
)

// TODO create tests for activate user with activation token
// from previous create post

func TestActivateUser_Assert_ActivationToken(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		User model.User
	}

	err := c.Post(
		activateUser,
		&resp,
		client.Var(
			"payload",
			model.UserActivate{
				ActivationToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIzNTViYzZlNi1iYTg4LTRkNzMtODk2My0zZTUyNjRhYmQ3OGMiLCJleHAiOjE2Mjg0MjIxMDgsImlhdCI6MTYyODQyMTgwOH0.S_S3u31AkcnHl5yYZWd9ibOVKaGXpiOg7Qm5LPa8208",
			},
		),
	)

	assert.Equal(err.Error(), `[{"message":"user cannot be activated","path":["activateUser"]}]`)

}
*/

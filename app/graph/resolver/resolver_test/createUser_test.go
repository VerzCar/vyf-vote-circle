package resolver_test

/*
import (
	assertPkg "github.com/stretchr/testify/assert"
	model2 "gitlab.vecomentman.com/vote-your-face/service/vote_circle/api/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/cache/testcache"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database/testdb"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils/testing/client"
	"testing"
	"time"
)

const (
	createUser = `mutation createUser($user: UserCreate!, $resendActivationEmail: Boolean) {
  createUser(user: $user, resendActivationEmail: $resendActivationEmail) {
    emailSend
    emailSendTo
    emailVerificationExpiration
  }
}`
)

func TestCreateUser(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		CreateUser model2.UserCreated
	}

	err := c.Post(
		createUser,
		&resp,
		client.Var(
			"user",
			model2.UserCreate{
				Email:     factory.ElonNew.Email,
				Password:  factory.ElonNew.Password,
				FirstName: factory.ElonNew.FirstName,
				LastName:  nil,
				Gender:    &factory.ElonNew.Gender,
				Locale:    &model2.LocaleInput{ID: factory.ElonNew.Locale.ID},
				Address:   nil,
			},
		), client.Var(
			"resendActivationEmail",
			false,
		),
	)

	assert.NoError(err)

	assert.Equal(resp.CreateUser.EmailSendTo, factory.ElonNew.Email)
	assert.Equal(resp.CreateUser.EmailSend, true)
	expectedTime := time.Now()
	timeDelta := utils.FormatDuration(resolver.Config.Ttl.Token.Account.Activation)
	expectedTime.Add(timeDelta)
	assert.Greater(resp.CreateUser.EmailVerificationExpiration.String(), expectedTime.String())

	// resend email without given user values
	err = c.Post(
		createUser,
		&resp,
		client.Var(
			"user",
			model2.UserCreate{
				Email:     factory.ElonNew.Email,
				Password:  "",
				FirstName: "",
				LastName:  nil,
				Gender:    nil,
				Locale:    nil,
				Address:   nil,
			},
		), client.Var(
			"resendActivationEmail",
			true,
		),
	)

	assert.NoError(err)

	assert.Equal(resp.CreateUser.EmailSendTo, factory.ElonNew.Email)
	assert.Equal(resp.CreateUser.EmailSend, true)
	expectedTime = time.Now()
	timeDelta = utils.FormatDuration(resolver.Config.Ttl.Token.Account.Activation)
	expectedTime.Add(timeDelta)
	assert.Greater(resp.CreateUser.EmailVerificationExpiration.String(), expectedTime.String())

	// assertion email already sent
	err = c.Post(
		createUser,
		&resp,
		client.Var(
			"user",
			model2.UserCreate{
				Email:     factory.ElonNew.Email,
				Password:  factory.ElonNew.Password,
				FirstName: factory.ElonNew.FirstName,
				LastName:  nil,
				Gender:    &factory.ElonNew.Gender,
				Locale:    &model2.LocaleInput{ID: factory.ElonNew.Locale.ID},
				Address:   nil,
			},
		), client.Var(
			"resendActivationEmail",
			false,
		),
	)

	assert.Equal(err.Error(), `[{"message":"email activation already sent","path":["createUser"]}]`)

}

func TestCreateUser_Assert_Email_Already_Exists(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	factory.Setup(resolver.DB)

	var resp struct {
		CreateUser model2.UserCreated
	}

	err := c.Post(
		createUser,
		&resp,
		client.Var(
			"user",
			model2.UserCreate{
				Email:     factory.Martin.Email,
				Password:  factory.Martin.Password,
				FirstName: factory.Martin.FirstName,
				LastName:  nil,
				Gender:    &factory.Martin.Gender,
				Locale:    &model2.LocaleInput{ID: factory.ElonNew.Locale.ID},
				Address:   nil,
			},
		), client.Var(
			"resendActivationEmail",
			false,
		),
	)

	assert.Equal(err.Error(), `[{"message":"email address already exists","path":["createUser"]}]`)

}

func TestCreateUser_Assert_PWD_Complexity_Not_Sufficient(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		CreateUser model2.UserCreated
	}

	err := c.Post(
		createUser,
		&resp,
		client.Var(
			"user",
			model2.UserCreate{
				Email:     factory.ElonNew.Email,
				Password:  "123",
				FirstName: factory.ElonNew.FirstName,
				LastName:  nil,
				Gender:    &factory.ElonNew.Gender,
				Locale:    &model2.LocaleInput{ID: factory.ElonNew.Locale.ID},
				Address:   nil,
			},
		), client.Var(
			"resendActivationEmail",
			false,
		),
	)

	assert.Equal(err.Error(), `[{"message":"password complexity not sufficient","path":["createUser"]}]`)

}

func TestCreateUser_Assert_Invalid_Contact(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		CreateUser model2.UserCreated
	}

	err := c.Post(
		createUser,
		&resp,
		client.Var(
			"user",
			model2.UserCreate{
				Email:     factory.ElonNew.Email,
				Password:  factory.ElonNew.Password,
				FirstName: factory.ElonNew.FirstName,
				LastName:  nil,
				Gender:    &factory.ElonNew.Gender,
				Locale:    &model2.LocaleInput{ID: factory.ElonNew.Locale.ID},
				Address:   nil,
				Contact: &model2.ContactInput{
					Email:                        "",
					PhoneNumber:                  "017683464774",
					PhoneNumberCountryAlphaCode:  &model2.CountryInput{ID: 65},
					PhoneNumber2:                 nil,
					PhoneNumber2CountryAlphaCode: nil,
					Web:                          nil,
				},
			},
		), client.Var(
			"resendActivationEmail",
			false,
		),
	)

	assert.Equal(err.Error(), `[{"message":"contact input invalid","path":["createUser"]}]`)

}

func TestCreateUser_Assert_InvalidAddress(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	var resp struct {
		CreateUser model2.UserCreated
	}

	err := c.Post(
		createUser,
		&resp,
		client.Var(
			"user",
			model2.UserCreate{
				Email:     factory.ElonNew.Email,
				Password:  factory.ElonNew.Password,
				FirstName: factory.ElonNew.FirstName,
				LastName:  nil,
				Gender:    &factory.ElonNew.Gender,
				Locale:    &model2.LocaleInput{ID: factory.ElonNew.Locale.ID},
				Address: &model2.AddressInput{
					Address:    "Street 1",
					City:       "",
					PostalCode: "83410",
					Country:    &model2.CountryInput{ID: 65},
				},
			},
		), client.Var(
			"resendActivationEmail",
			false,
		),
	)

	assert.Equal(err.Error(), `[{"message":"address input invalid","path":["createUser"]}]`)

}
*/

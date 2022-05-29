package resolver_test

/*import (
	"fmt"
	assertPkg "github.com/stretchr/testify/assert"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/cache/testcache"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/database/testdb"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/graph/client"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/model"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/users/test/factory"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils/testing/dockerpool"
	"testing"
)

const (
	createCompany = `mutation createCompany($company: CompanyCreate!) {
  createCompany(company: $company) {
    id
    name
    owner {
      id
    }
    users {
      id
    }
    address {
      city
    }
  }
}`
)

func TestCreateCompany(t *testing.T) {
	assert := assertPkg.New(t)

	testdb.Reset(resolver.DB, resolver.Log)
	testcache.Reset(resolver.Rdb)
	factory.Setup(resolver.DB)

	// create user svc as a container to test it against payment service
	// this is necessary as the payment service needs to query the user svc
	pool := dockerpool.Pool()
	userSvc := dockerpool.UserSvc(pool, resolver.Config)

	if err := pool.Retry(
		func() error {

			if !userSvc.Container.State.Running {
				return fmt.Errorf("container is not running")
			}
			t.Logf("Container Port %s", userSvc.GetPort("8080/tcp"))
			t.Logf("Container running: %t", userSvc.Container.State.Running)

			return nil
		},
	); err != nil {
		assert.Failf("Could not connect to docker: %s", err.Error())
	}

	paymentId := "3749f28c-65f9-4310-8d75-199b1a3334d6"

	gqlClient := client.New()

	var respToken struct {
		AuthToken model.Token
	}

	err := gqlClient.Post(
		authToken,
		&respToken,
		client.Path("http://localhost:"+userSvc.GetPort("8080/tcp")),
		client.Var(
			"credentials",
			model.Credentials{
				Email:    factory.Martin.Email,
				Password: factory.Martin.Password,
			},
		),
	)

	if err != nil {
		assert.Fail(err.Error())
	}

	var resp struct {
		Company model.Company
	}

	err = gqlClient.Post(
		createCompany,
		&resp,
		client.Path("http://localhost:"+userSvc.GetPort("8080/tcp")),
		client.AddHeader(
			"Authorization",
			fmt.Sprintf("%s %s", resolver.Config.Token.Type, respToken.AuthToken.Token),
		),
		client.Var(
			"company",
			model.CompanyCreate{
				Name: factory.NewCompany.Name,
				Address: &model.AddressInput{
					Address:    factory.NewCompany.Address.Address,
					City:       factory.NewCompany.Address.City,
					PostalCode: factory.NewCompany.Address.PostalCode,
					Country:    factory.NewCompany.Address.Country,
				},
				Contact: &model.ContactInput{
					Email:          factory.NewCompany.Contact.Email,
					PhoneNumber:    factory.NewCompany.Contact.PhoneNumber,
					PhoneNumberCc:  factory.NewCompany.Contact.PhoneNumberCc,
					PhoneNumber2:   &factory.NewCompany.Contact.PhoneNumber2,
					PhoneNumber2cc: &factory.NewCompany.Contact.PhoneNumber2cc,
					Web:            &factory.NewCompany.Contact.Web,
				},
				TaxID:     factory.NewCompany.TaxID,
				PaymentId: paymentId,
			},
		),
	)

	if err != nil {
		assert.Failf("Mutation failed: %s", err.Error())
	}

	assert.Equal(resp.Company.Name, factory.NewCompany.Name)

}
*/

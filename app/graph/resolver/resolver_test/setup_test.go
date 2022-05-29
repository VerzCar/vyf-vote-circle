package resolver_test

/*import (
	"github.com/gin-gonic/gin"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/cache"
	testdb2 "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/database/testdb"
	gqlResolver "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/graph/resolver"
	mainRouter "gitlab.vecomentman.com/vote-your-face/service/vote_circle/app/router"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/config"
	"gitlab.vecomentman.com/libs/logger"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils"
	"gitlab.vecomentman.com/vote-your-face/service/vote_circle/utils/testing/client"
	"os"
	"testing"
)

var (
	resolver gqlResolver.Resolver
	router   *gin.Engine
	c        *client.Client
)

// Setup test env case
func TestMain(m *testing.M) {
	configPath := utils.FromBase("config/")

	resolver.Config = config.Load(configPath)
	resolver.Log = logger.NewLogger(configPath)

	resolver.Rdb = cache.Connect(resolver.Log, resolver.Config)

	resolver.DB = testdb2.Connect(resolver.Log, resolver.Config)
	testdb2.Setup(resolver.DB, resolver.Log, resolver.Config)

	router = mainRouter.Setup(&resolver)

	c = client.New(router)

	code := m.Run()

	os.Exit(code)
}*/

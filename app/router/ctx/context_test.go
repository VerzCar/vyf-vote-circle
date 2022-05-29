package ctx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestSetGinContext(t *testing.T) {
	req := &http.Request{}
	parentGinContext := &gin.Context{
		Request: req,
	}
	SetGinContext(parentGinContext)

	ginContext := parentGinContext.Request.Context().Value(ginContextKey)
	require.NotNil(t, ginContext)

	_, ok := ginContext.(*gin.Context)

	require.True(t, ok)
}

func TestContextToGinContext(t *testing.T) {
	ctx := context.Background()
	parentGinContext := &gin.Context{}
	testContext := context.WithValue(ctx, ginContextKey, parentGinContext)
	ginCtx, err := ContextToGinContext(testContext)

	require.Nil(t, err)
	require.Empty(t, ginCtx.Params.ByName("example"))
}

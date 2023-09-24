package ctx

import (
	"context"
	"fmt"
	"github.com/VerzCar/vyf-lib-awsx"
	"github.com/gin-gonic/gin"
)

const ginContextKey = "GinCtxKey"
const authClaimsContextKey = "AuthClaimsContextKey"

func SetGinContext(ctx *gin.Context) {
	c := context.WithValue(ctx.Request.Context(), ginContextKey, ctx)
	ctx.Request = ctx.Request.WithContext(c)
}

func ContextToGinContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)

	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)

	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}

	return gc, nil
}

func SetAuthClaimsContext(ctx *gin.Context, val interface{}) {
	c := context.WithValue(ctx.Request.Context(), authClaimsContextKey, val)
	ctx.Request = ctx.Request.WithContext(c)
}

func ContextToAuthClaims(ctx context.Context) (*awsx.JWTToken, error) {
	authClaimsValue := ctx.Value(authClaimsContextKey)

	if authClaimsValue == nil {
		err := fmt.Errorf("could not retrieve auth claims")
		return nil, err
	}

	authClaims, ok := authClaimsValue.(*awsx.JWTToken)

	if !ok {
		err := fmt.Errorf("auth claims has wrong type")
		return nil, err
	}

	return authClaims, nil
}

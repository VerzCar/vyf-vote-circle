package ctx

import (
	"context"
	"fmt"
	"github.com/VerzCar/vyf-lib-awsx"
	"github.com/gin-gonic/gin"
)

const authClaimsContextKey = "AuthClaimsContextKey"

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

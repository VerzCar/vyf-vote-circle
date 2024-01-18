package ctx

import (
	"context"
	"github.com/VerzCar/vyf-lib-awsx"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestSetAuthClaimsContext(t *testing.T) {
	tests := []struct {
		name string
		val  interface{}
	}{
		{
			name: "Set valid auth claims",
			val:  &awsx.JWTToken{Issuer: "testIssuer"},
		},
		{
			name: "Set nil auth claims",
			val:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx := &gin.Context{
					Request: &http.Request{},
				}
				SetAuthClaimsContext(ctx, tt.val)

				got := ctx.Request.Context().Value(authClaimsContextKey)
				assert.Equal(t, tt.val, got)
			},
		)
	}
}

func TestContextToAuthClaims(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		wantErr bool
	}{
		{
			name:    "Get valid auth claims",
			val:     &awsx.JWTToken{Issuer: "testIssuer"},
			wantErr: false,
		},
		{
			name:    "Get nil auth claims",
			val:     nil,
			wantErr: true,
		},
		{
			name:    "Get invalid auth claims",
			val:     "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				ctx := context.WithValue(context.Background(), authClaimsContextKey, tt.val)

				got, err := ContextToAuthClaims(ctx)
				if tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.Equal(t, tt.val, got)
				}
			},
		)
	}
}

package header

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// Authorization gets the value in the HTTP: Authorization header from the gin context
// as string.
func Authorization(c *gin.Context, headerType string) (string, error) {
	var err error
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		err = fmt.Errorf("authorization header is empty")
		return "", err
	}

	// support both cases bearer and Bearer prefix
	prefixHeaderLow := strings.ToLower(headerType + " ")
	prefixHeaderTitle := strings.Title(headerType + " ")

	authToken := ""

	switch {
	case strings.HasPrefix(authHeader, prefixHeaderTitle):
		authToken = strings.TrimPrefix(authHeader, prefixHeaderTitle)
	case strings.HasPrefix(authHeader, prefixHeaderLow):
		authToken = strings.TrimPrefix(authHeader, prefixHeaderLow)
	default:
		err = fmt.Errorf("wrong header [Authorization] type - not a " + headerType + " token")
		return "", err
	}

	return authToken, nil
}

// BearerToken prepares the given access token to with the Bearer prefix.
// Returns the formatted access token with the bearer prefix.
func BearerToken(accessToken string) string {
	return fmt.Sprintf("%s %s", "Bearer", accessToken)
}

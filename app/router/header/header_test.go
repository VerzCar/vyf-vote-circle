package header_test

import (
	"fmt"
	"github.com/VerzCar/vyf-vote-circle/app/router/header"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var headerType string

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST(
		"/ping", func(c *gin.Context) {
			token, err := header.Authorization(c, headerType)

			if err != nil {
				c.AbortWithError(http.StatusBadRequest, err)
				return
			}

			c.String(http.StatusOK, token)
		},
	)
	return r
}

func TestAuthorization(t *testing.T) {
	router := setupRouter()

	expectedToken := "ey234fft34r0434frfgtgb5t"

	req, _ := http.NewRequest(http.MethodPost, "/ping", nil)

	tests := []struct {
		name          string
		w             *httptest.ResponseRecorder
		req           *http.Request
		expectedToken string
		headerType    string
		deleteHeader  bool
		want          int
	}{
		{
			name:          "should extract token from header successfully",
			w:             httptest.NewRecorder(),
			req:           req,
			expectedToken: expectedToken,
			headerType:    "Bearer",
			want:          http.StatusOK,
		},
		{
			name:          "should extract token from header with another bearer prefix successfully",
			w:             httptest.NewRecorder(),
			req:           req,
			expectedToken: expectedToken,
			headerType:    "bearer",
			want:          http.StatusOK,
		},
		{
			name:          "should fail because bearer prefix does not exist",
			w:             httptest.NewRecorder(),
			req:           req,
			expectedToken: expectedToken,
			headerType:    "",
			want:          http.StatusBadRequest,
		},
		{
			name:          "should fail because auth header does not exist",
			w:             httptest.NewRecorder(),
			req:           req,
			expectedToken: expectedToken,
			headerType:    "Bearer",
			deleteHeader:  true,
			want:          http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {

				headerType = test.headerType

				if !test.deleteHeader {
					req.Header.Set("Authorization", constructToken(headerType, expectedToken))
					if test.headerType == "" {
						req.Header.Set("Authorization", expectedToken)
					}
				} else {
					req.Header.Del("Authorization")
				}

				router.ServeHTTP(test.w, test.req)

				if !reflect.DeepEqual(test.w.Code, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, test.w.Code, test.want)
				}

				if test.w.Code == http.StatusOK && !reflect.DeepEqual(test.w.Body.String(), test.expectedToken) {
					t.Errorf(
						"test: %v failed. \ntoken: %v \nwanted: %v",
						test.name,
						test.w.Body.String(),
						test.expectedToken,
					)
				}
			},
		)
	}
}

func TestBearerToken(t *testing.T) {

	expectedBearerToken := "Bearer ey234fft34r0434frfgtgb5t"

	tests := []struct {
		name  string
		token string
		want  string
	}{
		{
			name:  "should extract token from header successfully",
			token: "ey234fft34r0434frfgtgb5t",
			want:  expectedBearerToken,
		},
	}

	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				bearerToken := header.BearerToken(test.token)

				if !reflect.DeepEqual(bearerToken, test.want) {
					t.Errorf("test: %v failed. \ngot: %v \nwanted: %v", test.name, bearerToken, test.want)
				}
			},
		)
	}
}

func constructToken(headerType string, token string) string {
	return fmt.Sprintf("%s %s", headerType, token)
}

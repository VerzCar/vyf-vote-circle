package client

import (
	"net/http"
	netUrl "net/url"
)

// Var adds a variable into the outgoing request
func Var(name string, value interface{}) Option {
	return func(bd *Request) {
		if bd.Variables == nil {
			bd.Variables = map[string]interface{}{}
		}

		bd.Variables[name] = value
	}
}

// Operation sets the operation name for the outgoing request
func Operation(name string) Option {
	return func(bd *Request) {
		bd.OperationName = name
	}
}

// Extensions sets the extensions to be sent with the outgoing request
func Extensions(extensions map[string]interface{}) Option {
	return func(bd *Request) {
		bd.Extensions = extensions
	}
}

// Path sets the url that this request will be made against.
func Path(url string) Option {
	return func(bd *Request) {
		u, _ := netUrl.Parse(url)
		bd.HTTP.URL = u
	}
}

// AddHeader adds a header to the outgoing request.
func AddHeader(key string, value string) Option {
	return func(bd *Request) {
		bd.HTTP.Header.Add(key, value)
	}
}

// BasicAuth authenticates the request using http basic auth.
func BasicAuth(username, password string) Option {
	return func(bd *Request) {
		bd.HTTP.SetBasicAuth(username, password)
	}
}

// AddCookie adds a cookie to the outgoing request
func AddCookie(cookie *http.Cookie) Option {
	return func(bd *Request) {
		bd.HTTP.AddCookie(cookie)
	}
}

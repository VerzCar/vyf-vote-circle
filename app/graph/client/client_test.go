package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/require"
)

// RoundTripFunc .
type RoundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

//NewTestClient returns *http.Client with Transport replaced to avoid making real calls
func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func TestNew(t *testing.T) {
	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`OK`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)

	c := New(tClient, Path("https://example.com"))
	require.NotEmpty(t, c.opts)
	require.NotEmpty(t, c.httpClient)
}

func TestPost(t *testing.T) {
	var resp struct {
		Name string
	}

	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			require.Equal(t, `{"query":"user(id:$id){name}","variables":{"id":1}}`, string(b))

			w := new(bytes.Buffer)
			err = json.NewEncoder(w).Encode(
				map[string]interface{}{
					"data": map[string]interface{}{
						"name": "bob",
					},
				},
			)

			if err != nil {
				panic(err)
			}
			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(w),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)

	c := New(tClient)

	err := c.Post("user(id:$id){name}", &resp, Var("id", 1))

	require.Nil(t, err)
	require.Equal(t, "bob", resp.Name)
}

func TestClientMultipartFormData(t *testing.T) {
	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			require.NoError(t, err)
			require.Contains(t, string(b), `Content-Disposition: form-data; name="operations"`)
			require.Contains(t, string(b), `{"query":"mutation ($input: Input!) {}","variables":{"file":{}}`)
			require.Contains(t, string(b), `Content-Disposition: form-data; name="map"`)
			require.Contains(t, string(b), `{"0":["variables.file"]}`)
			require.Contains(t, string(b), `Content-Disposition: form-data; name="0"; filename="example.txt"`)
			require.Contains(t, string(b), `Content-Type: text/plain`)
			require.Contains(t, string(b), `Hello World`)

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)
	var resp struct{}

	c := New(tClient)

	c.Post(
		"{ id }",
		&resp,
		func(bd *Request) {
			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			bodyWriter.WriteField("operations", `{"query":"mutation ($input: Input!) {}","variables":{"file":{}}`)
			bodyWriter.WriteField("map", `{"0":["variables.file"]}`)

			h := make(textproto.MIMEHeader)
			h.Set("Content-Disposition", `form-data; name="0"; filename="example.txt"`)
			h.Set("Content-Type", "text/plain")
			ff, _ := bodyWriter.CreatePart(h)
			ff.Write([]byte("Hello World"))
			bodyWriter.Close()

			bd.HTTP.Body = ioutil.NopCloser(bodyBuf)
			bd.HTTP.Header.Set("Content-Type", bodyWriter.FormDataContentType())
		},
	)
}

func TestAddHeader(t *testing.T) {
	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			require.Equal(t, "ASDF", req.Header.Get("Test-Key"))

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)

	c := New(tClient)

	var resp struct{}
	c.Post(
		"{ id }",
		&resp,
		AddHeader("Test-Key", "ASDF"),
	)
}

func TestAddClientHeader(t *testing.T) {
	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			require.Equal(t, "ASDF", req.Header.Get("Test-Key"))

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)

	c := New(tClient, AddHeader("Test-Key", "ASDF"))

	var resp struct{}
	c.Post(
		"{ id }",
		&resp,
	)
}

func TestBasicAuth(t *testing.T) {
	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			user, pass, ok := req.BasicAuth()
			require.True(t, ok)
			require.Equal(t, "user", user)
			require.Equal(t, "pass", pass)

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)

	c := New(tClient)

	var resp struct{}
	c.Post(
		"{ id }",
		&resp,
		BasicAuth("user", "pass"),
	)
}

func TestAddCookie(t *testing.T) {
	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			c, err := req.Cookie("foo")
			require.NoError(t, err)
			require.Equal(t, "value", c.Value)

			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(bytes.NewBufferString(`{}`)),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)

	c := New(tClient)

	var resp struct{}
	c.Post(
		"{ id }",
		&resp,
		AddCookie(&http.Cookie{Name: "foo", Value: "value"}),
	)
}

func TestAddExtensions(t *testing.T) {
	tClient := NewTestClient(
		func(req *http.Request) *http.Response {
			b, err := ioutil.ReadAll(req.Body)
			if err != nil {
				panic(err)
			}
			require.Equal(
				t,
				`{"query":"user(id:1){name}","extensions":{"persistedQuery":{"sha256Hash":"ceec2897e2da519612279e63f24658c3e91194cbb2974744fa9007a7e1e9f9e7","version":1}}}`,
				string(b),
			)

			w := new(bytes.Buffer)
			err = json.NewEncoder(w).Encode(
				map[string]interface{}{
					"data": map[string]interface{}{
						"name": "bob",
					},
				},
			)

			if err != nil {
				panic(err)
			}
			return &http.Response{
				StatusCode: 200,
				// Send response to be tested
				Body: ioutil.NopCloser(w),
				// Must be set to non-nil value or it panics
				Header: make(http.Header),
			}
		},
	)

	c := New(tClient)

	var resp struct {
		Name string
	}
	c.Post(
		"user(id:1){name}",
		&resp,
		Extensions(
			map[string]interface{}{
				"persistedQuery": map[string]interface{}{
					"version":    1,
					"sha256Hash": "ceec2897e2da519612279e63f24658c3e91194cbb2974744fa9007a7e1e9f9e7",
				},
			},
		),
	)
}

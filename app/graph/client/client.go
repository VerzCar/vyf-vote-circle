package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"io"
	"net/http"
	"reflect"
	"regexp"
	"time"
)

type (
	// Client used GraphQL servers.
	Client struct {
		httpClient *http.Client
		opts       []Option
	}

	// Option implements a visitor that mutates an outgoing GraphQL request
	//
	// This is the Option pattern - https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
	Option func(bd *Request)

	// Request represents an outgoing GraphQL request
	Request struct {
		Query         string                 `json:"query"`
		Variables     map[string]interface{} `json:"variables,omitempty"`
		OperationName string                 `json:"operationName,omitempty"`
		Extensions    map[string]interface{} `json:"extensions,omitempty"`
		HTTP          *http.Request          `json:"-"`
	}

	// Response is a GraphQL layer response from a handler.
	Response struct {
		Data       interface{}
		Errors     json.RawMessage
		Extensions map[string]interface{}
	}
)

// New creates a graphql client
// Options can be set that should be applied to all requests made with this client
func New(httpClient *http.Client, opts ...Option) *Client {
	p := &Client{
		httpClient: httpClient,
		opts:       opts,
	}

	return p
}

// Post sends a http POST request to the graphql endpoint with the given query then unpacks
// the response into the given object.
func (p *Client) Post(query string, response interface{}, options ...Option) error {
	respDataRaw, err := p.RawPost(query, options...)
	if err != nil {
		return err
	}

	unpackErr := unpack(respDataRaw.Data, response)

	if respDataRaw.Errors != nil {
		return RawJsonError{respDataRaw.Errors}
	}
	return unpackErr
}

// RawPost is similar to Post, except it skips decoding the raw json response
// unpacked onto Response. This is used to test extension keys which are not
// available when using Post.
func (p *Client) RawPost(query string, options ...Option) (*Response, error) {
	r, err := p.newRequest(query, options...)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %s", err.Error())
	}

	resp, err := p.httpClient.Do(r)

	if err != nil {
		return nil, fmt.Errorf("post failed: %s", err.Error())
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("could not read response body: %s", err.Error())
	}

	// decode it into map string first, let mapstructure do the final decode
	// because it can be much stricter about unknown fields.
	respDataRaw := &Response{}
	err = json.Unmarshal(bodyBytes, &respDataRaw)
	if err != nil {
		return nil, fmt.Errorf("could not decode response body: %s", err.Error())
	}

	return respDataRaw, nil
}

func (p *Client) newRequest(query string, options ...Option) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, "/", nil)

	if err != nil {
		return nil, fmt.Errorf("create request failed: %s", err)
	}

	bd := &Request{
		Query: query,
		HTTP:  req,
	}

	bd.HTTP.Header.Set("Content-Type", "application/json")

	// per client options from client.New apply first
	for _, option := range p.opts {
		option(bd)
	}
	// per request options
	for _, option := range options {
		option(bd)
	}

	contentType := bd.HTTP.Header.Get("Content-Type")
	switch {
	case regexp.MustCompile(`multipart/form-data; ?boundary=.*`).MatchString(contentType):
		break
	case "application/json" == contentType:
		requestBody, err := json.Marshal(bd)
		if err != nil {
			return nil, fmt.Errorf("encode: %s", err.Error())
		}
		bd.HTTP.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	default:
		panic("unsupported encoding" + bd.HTTP.Header.Get("Content-Type"))
	}

	return bd.HTTP, nil
}

func unpack(data interface{}, into interface{}) error {
	d, err := mapstructure.NewDecoder(
		&mapstructure.DecoderConfig{
			Result:      into,
			TagName:     "json",
			ErrorUnused: true,
			ZeroFields:  true,
			Squash:      true,
			DecodeHook: mapstructure.ComposeDecodeHookFunc(
				toTimeHookFunc(),
			),
		},
	)
	if err != nil {
		return fmt.Errorf("mapstructure: %s", err.Error())
	}

	return d.Decode(data)
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{},
	) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

// +build !exported

package dairyclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	exampleID       = 666
	exampleURL      = `http://www.dairycart.com`
	exampleUsername = `username`
	examplePassword = `password` // lol not really
	exampleSKU      = `sku`
	exampleBadJSON  = `{"invalid lol}`
)

func handlerGenerator(handlers map[string]func(res http.ResponseWriter, req *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, handlerFunc := range handlers {
			if r.URL.Path == path {
				handlerFunc(w, r)
				return
			}
		}
	})
}

func createInternalClient(t *testing.T, ts *httptest.Server) *V1Client {
	u, err := url.Parse(ts.URL)
	assert.Nil(t, err, "no error should be returned when parsing a test server's URL")

	c := &V1Client{
		Client: ts.Client(),
		AuthCookie: &http.Cookie{
			Name: "dairycart",
		},
		URL: u,
	}
	return c
}

func TestExecuteRequestAddsCookieToRequests(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			cookies := req.Cookies()
			if len(cookies) == 0 {
				assert.FailNow(t, "no cookies attached to the request")
			}

			cookieFound := false
			for _, c := range cookies {
				if c.Name == "dairycart" {
					cookieFound = true
				}
			}
			assert.True(t, cookieFound)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", ts.URL, exampleEndpoint), nil)
	assert.Nil(t, err, "no error should be returned when creating a new request")

	c.executeRequest(req)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestUnexportedBuildURL(t *testing.T) {
	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()
	c := createInternalClient(t, ts)

	testCases := []struct {
		query    map[string]string
		parts    []string
		expected string
	}{
		{
			query:    nil,
			parts:    []string{""},
			expected: fmt.Sprintf("%s/v1/", ts.URL),
		},
		{
			query:    nil,
			parts:    []string{"things", "and", "stuff"},
			expected: fmt.Sprintf("%s/v1/things/and/stuff", ts.URL),
		},
		{
			query:    map[string]string{"param": "value"},
			parts:    []string{"example"},
			expected: fmt.Sprintf("%s/v1/example?param=value", ts.URL),
		},
	}

	for _, tc := range testCases {
		actual := c.buildURL(tc.query, tc.parts...)
		assert.Equal(t, tc.expected, actual, "expected and actual built URLs don't match")
	}
}

func TestExists(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodHead, "exists should be making HEAD requests")
			res.WriteHeader(http.StatusOK)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	actual, err := c.exists(c.buildURL(nil, "whatever"))
	assert.Nil(t, err)
	assert.True(t, actual, "exists should return false when the status code is %d", http.StatusOK)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestExistsReturnsFalseWhen404IsReturned(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodHead, "exists should be making HEAD requests")
			res.WriteHeader(http.StatusNotFound)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	actual, err := c.exists(c.buildURL(nil, "whatever"))
	assert.Nil(t, err)
	assert.False(t, actual, "exists should return false when the status code is %d", http.StatusNotFound)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestExistsReturnsFalseAndErrorWhenFailingToExecuteRequest(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	c := createInternalClient(t, ts)
	ts.Close()

	actual, err := c.exists(c.buildURL(nil, "whatever"))
	assert.NotNil(t, err)
	assert.False(t, actual, "exists should return false when the status code is %d", http.StatusOK)
}

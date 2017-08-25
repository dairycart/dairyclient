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

func TestGet(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodGet, "get should be making GET requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	err := c.get(c.buildURL(nil, "whatever"), &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "actual struct should equal expected struct")
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestGetReturnsErrorWhenPassedNilOrNonPointer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()
	c := createInternalClient(t, ts)

	actual := struct {
		Things string `json:"things"`
	}{}

	ptrErr := c.get(c.buildURL(nil, "whatever"), actual)
	assert.NotNil(t, ptrErr)

	nilErr := c.get(c.buildURL(nil, "whatever"), nil)
	assert.NotNil(t, nilErr)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodDelete, "delete should be making DELETE requests")
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	err := c.delete(c.buildURL(nil, "whatever"))
	assert.Nil(t, err)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestDeleteReturnsErroWhenFailingToExecuteRequest(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	c := createInternalClient(t, ts)
	ts.Close()

	err := c.delete(c.buildURL(nil, "whatever"))
	assert.NotNil(t, err)
}

func TestDeleteReturnsErrorWhenStatusCodeIsNot200(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodDelete, "delete should be making DELETE requests")
			res.WriteHeader(http.StatusInternalServerError)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	err := c.delete(c.buildURL(nil, "whatever"))
	assert.NotNil(t, err)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestMakeDataRequest(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/post/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodPost, "makeDataRequest should only be making PUT or POST requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "post", "whatever"), expected, &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "actual struct should equal expected struct")
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestMakeDataRequestReturnsErrorWhenPassedNilOrNonPointer(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	c := createInternalClient(t, ts)
	ts.Close()

	ptrErr := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), struct{}{}, struct{}{})
	assert.NotNil(t, ptrErr, "makeDataRequest should return an error when passed a non-pointer output param")

	nilErr := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), struct{}{}, nil)
	assert.NotNil(t, nilErr, "makeDataRequest should return an error when passed a nil output param")
}

func TestMakeDataRequestReturnsErrorWhenPassedAnInvalidInputStruct(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	c := createInternalClient(t, ts)
	ts.Close()

	f := &testBreakableStruct{Thing: "dongs"}
	err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), f, &struct{}{})
	assert.NotNil(t, err, "makeDataRequest should return an error when passed an invalid input struct")
}

func TestMakeDataRequestReturnsErrorWhenFailingToExecuteRequest(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	c := createInternalClient(t, ts)
	ts.Close()

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "post", "whatever"), expected, &actual)
	assert.NotNil(t, err, "makeDataRequest should return an error when failing to execute request")
}

func TestMakeDataRequestReturnsErrorWhenFailingToUnmarshalBody(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			fmt.Fprintf(res, exampleBadJSON)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), expected, &actual)
	assert.NotNil(t, err)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestPost(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodPost, "post should only be making POST requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	err := c.post(c.buildURL(nil, "whatever"), expected, &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "actual struct should equal expected struct")
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestPatch(t *testing.T) {
	t.Parallel()
	var endpointCalled bool
	exampleEndpoint := "/v1/whatever"

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		exampleEndpoint: func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.Equal(t, req.Method, http.MethodPatch, "patch should only be making PATCH requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	expected := struct {
		Things string `json:"things"`
	}{
		Things: "stuff",
	}

	actual := struct {
		Things string `json:"things"`
	}{}

	err := c.patch(c.buildURL(nil, "whatever"), expected, &actual)
	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "actual struct should equal expected struct")
	assert.True(t, endpointCalled, "endpoint should have been called")
}

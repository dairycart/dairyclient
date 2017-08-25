// +build !exported

package dairyclient

import (
	"fmt"
	"io/ioutil"
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

type subtest struct {
	Message string
	Test    func(t *testing.T)
}

////////////////////////////////////////////////////////
//                                                    //
//                 Helper Functions                   //
//                                                    //
////////////////////////////////////////////////////////

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

func runSubtestSuite(t *testing.T, tests []subtest) {
	testPassed := true
	for _, test := range tests {
		if !testPassed {
			t.FailNow()
		}
		testPassed = t.Run(test.Message, test.Test)
	}
}

////////////////////////////////////////////////////////
//                                                    //
//                   Actual Tests                     //
//                                                    //
////////////////////////////////////////////////////////

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

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", ts.URL, exampleEndpoint), nil)
	assert.Nil(t, err, "no error should be returned when creating a new request")

	c.executeRequest(req)
	assert.True(t, endpointCalled, "endpoint should have been called")
}

func TestUnexportedBuildURL(t *testing.T) {
	ts := httptest.NewTLSServer(http.NotFoundHandler())
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
	var normalEndpointCalled bool
	var fourOhFourEndpointCalled bool

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/normal": func(res http.ResponseWriter, req *http.Request) {
			normalEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodHead, "exists should be making HEAD requests")
			res.WriteHeader(http.StatusOK)
		},
		"/v1/four_oh_four": func(res http.ResponseWriter, req *http.Request) {
			fourOhFourEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodHead, "exists should be making HEAD requests")
			res.WriteHeader(http.StatusNotFound)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	c := createInternalClient(t, ts)

	normalUse := func(t *testing.T) {
		actual, err := c.exists(c.buildURL(nil, "normal"))
		assert.Nil(t, err)
		assert.True(t, actual, "exists should return false when the status code is %d", http.StatusOK)
		assert.True(t, normalEndpointCalled, "endpoint should have been called")
	}

	notFound := func(t *testing.T) {
		actual, err := c.exists(c.buildURL(nil, "four_oh_four"))
		assert.Nil(t, err)
		assert.False(t, actual, "exists should return false when the status code is %d", http.StatusNotFound)
		assert.True(t, fourOhFourEndpointCalled, "endpoint should have been called")
	}

	failsToRequest := func(t *testing.T) {
		ts.Close()
		actual, err := c.exists(c.buildURL(nil, "whatever"))
		assert.NotNil(t, err)
		assert.False(t, actual, "exists should return false when the status code is %d", http.StatusOK)
	}

	subtests := []subtest{
		{
			Message: "normal use",
			Test:    normalUse,
		},
		{
			Message: "not found",
			Test:    notFound,
		},
		{
			Message: "failed request",
			Test:    failsToRequest,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestGet(t *testing.T) {
	t.Parallel()
	var normalEndpointCalled bool
	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/normal": func(res http.ResponseWriter, req *http.Request) {
			normalEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodGet, "get should be making GET requests")
			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := createInternalClient(t, ts)

	normalUse := func(t *testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		err := c.get(c.buildURL(nil, "normal"), &actual)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "actual struct should equal expected struct")
		assert.True(t, normalEndpointCalled, "endpoint should have been called")
	}

	nilInput := func(t *testing.T) {
		nilErr := c.get(c.buildURL(nil, "whatever"), nil)
		assert.NotNil(t, nilErr)
	}

	nonPointerInput := func(t *testing.T) {
		actual := struct {
			Things string `json:"things"`
		}{}

		ptrErr := c.get(c.buildURL(nil, "whatever"), actual)
		assert.NotNil(t, ptrErr)
	}

	subtests := []subtest{
		{
			Message: "normal use",
			Test:    normalUse,
		},
		{
			Message: "nil input",
			Test:    nilInput,
		},
		{
			Message: "non-pointer input",
			Test:    nonPointerInput,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	var normalEndpointCalled bool
	var fiveHundredEndpointCalled bool

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/normal": func(res http.ResponseWriter, req *http.Request) {
			normalEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodDelete, "delete should be making DELETE requests")
		},
		"/v1/five_hundred": func(res http.ResponseWriter, req *http.Request) {
			fiveHundredEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodDelete, "delete should be making DELETE requests")
			res.WriteHeader(http.StatusInternalServerError)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	c := createInternalClient(t, ts)

	normalUse := func(t *testing.T) {
		err := c.delete(c.buildURL(nil, "normal"))
		assert.Nil(t, err)
		assert.True(t, normalEndpointCalled, "endpoint should have been called")
	}

	badStatusCode := func(t *testing.T) {
		err := c.delete(c.buildURL(nil, "five_hundred"))
		assert.NotNil(t, err)
		assert.True(t, fiveHundredEndpointCalled, "endpoint should have been called")
	}

	failedRequest := func(t *testing.T) {
		ts.Close()
		err := c.delete(c.buildURL(nil, "whatever"))
		assert.NotNil(t, err)
	}

	subtests := []subtest{
		{
			Message: "normal use",
			Test:    normalUse,
		},
		{
			Message: "bad status code",
			Test:    badStatusCode,
		},
		{
			Message: "failed request",
			Test:    failedRequest,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestMakeDataRequest(t *testing.T) {
	t.Parallel()
	var normalEndpointCalled bool
	var badJSONEndpointCalled bool

	handlers := map[string]func(res http.ResponseWriter, req *http.Request){
		"/v1/whatever": func(res http.ResponseWriter, req *http.Request) {
			normalEndpointCalled = true
			assert.Equal(t, req.Method, http.MethodPost, "makeDataRequest should only be making PUT or POST requests")

			bodyBytes, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			requestBody := string(bodyBytes)
			assert.Equal(t, requestBody, `{"things":"stuff"}`, "makeDataRequest should attach the correct JSON to the request body")

			exampleResponse := `
				{
					"things": "stuff"
				}
			`
			fmt.Fprintf(res, exampleResponse)
		},
		"/v1/bad_json": func(res http.ResponseWriter, req *http.Request) {
			badJSONEndpointCalled = true
			fmt.Fprintf(res, exampleBadJSON)
		},
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	c := createInternalClient(t, ts)

	normalUse := func(t *testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), expected, &actual)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "actual struct should equal expected struct")
		assert.True(t, normalEndpointCalled, "endpoint should have been called")
	}

	nilArgument := func(t *testing.T) {
		ptrErr := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), struct{}{}, struct{}{})
		assert.NotNil(t, ptrErr, "makeDataRequest should return an error when passed a non-pointer output param")
	}

	nonPtrArgument := func(t *testing.T) {
		nilErr := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), struct{}{}, nil)
		assert.NotNil(t, nilErr, "makeDataRequest should return an error when passed a nil output param")
	}

	invalidStructArgument := func(t *testing.T) {
		f := &testBreakableStruct{Thing: "dongs"}
		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), f, &struct{}{})
		assert.NotNil(t, err, "makeDataRequest should return an error when passed an invalid input struct")
	}

	unmarshalFailure := func(t *testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "bad_json"), expected, &actual)
		assert.NotNil(t, err)
		assert.True(t, badJSONEndpointCalled, "endpoint should have been called")
	}

	failedRequest := func(t *testing.T) {
		expected := struct {
			Things string `json:"things"`
		}{
			Things: "stuff",
		}

		actual := struct {
			Things string `json:"things"`
		}{}

		ts.Close()
		err := c.makeDataRequest(http.MethodPost, c.buildURL(nil, "whatever"), expected, &actual)
		assert.NotNil(t, err, "makeDataRequest should return an error when failing to execute request")
	}

	subtests := []subtest{
		{
			Message: "normal use",
			Test:    normalUse,
		},
		{
			Message: "nil argument",
			Test:    nilArgument,
		},
		{
			Message: "non-pointer argument",
			Test:    nonPtrArgument,
		},
		{
			Message: "invalid struct",
			Test:    invalidStructArgument,
		},
		{
			Message: "unmarshal failure",
			Test:    unmarshalFailure,
		},
		{
			Message: "failed request",
			Test:    failedRequest,
		},
	}
	runSubtestSuite(t, subtests)
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

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
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

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
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

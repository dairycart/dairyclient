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

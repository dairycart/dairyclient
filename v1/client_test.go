package dairyclient_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dairycart/dairyclient/v1"
)

////////////////////////////////////////////////////////
//                                                    //
//                 Constructor Tests                  //
//                                                    //
////////////////////////////////////////////////////////

func TestV1ClientConstructor(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
	defer ts.Close()

	c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

	assert.NotNil(t, c.AuthCookie)
	assert.Nil(t, err)
}

func TestV1ClientConstructorReturnsErrorWithInvalidURL(t *testing.T) {
	t.Parallel()
	c, err := dairyclient.NewV1Client(":", exampleUsername, examplePassword, nil)

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestV1ClientConstructorWithFailureToLogin(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
	ts.Close()

	c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestV1ClientConstructorWhereLoginCookieIsNotReturned(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(obligatoryLoginHandler(false))
	defer ts.Close()

	c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestDairyclientAddsCookieToRequest(t *testing.T) {
	t.Parallel()
	var endpointCalled bool

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/product/%s", exampleSKU): func(res http.ResponseWriter, req *http.Request) {
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
	c := buildTestClient(t, ts)

	res, err := c.ProductExists(exampleSKU)

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.True(t, endpointCalled, "endpoint should be called")
}

func TestBuildURL(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()
	c := buildTestClient(t, ts)

	expected := fmt.Sprintf("%s/v1/things/stuff?query=params", ts.URL)
	exampleParams := map[string]string{
		"query": "params",
	}
	actual, err := c.BuildURL(exampleParams, "things", "stuff")

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "BuildURL doesn't return the correct result. Expected `%s`, got `%s`", expected, actual)
}

func TestBuildURLReturnsErrorWhenFailingToParseURLParts(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()
	c := buildTestClient(t, ts)

	actual, err := c.BuildURL(nil, `%gh&%ij`)

	assert.NotNil(t, err)
	assert.Empty(t, actual)
}

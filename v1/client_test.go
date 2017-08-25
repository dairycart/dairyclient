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

	normalUse := func(t *testing.T) {
		ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
		defer ts.Close()

		c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

		assert.NotNil(t, c.AuthCookie)
		assert.Nil(t, err)
	}

	invalidURL := func(t *testing.T) {
		c, err := dairyclient.NewV1Client(":", exampleUsername, examplePassword, nil)
		assert.Nil(t, c)
		assert.NotNil(t, err)
	}

	loginFailure := func(t *testing.T) {
		ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
		ts.Close()
		c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

		assert.Nil(t, c)
		assert.NotNil(t, err)
	}

	sansCookie := func(t *testing.T) {
		ts := httptest.NewTLSServer(obligatoryLoginHandler(false))
		defer ts.Close()

		c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

		assert.Nil(t, c)
		assert.NotNil(t, err)
	}

	subtests := []subtest{
		{
			Message: "normal use",
			Test:    normalUse,
		},
		{
			Message: "invalid url",
			Test:    invalidURL,
		},
		{
			Message: "login failure",
			Test:    loginFailure,
		},
		{
			Message: "no cookie returned",
			Test:    sansCookie,
		},
	}
	runSubtestSuite(t, subtests)
}

func TestBuildURL(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.NotFoundHandler())
	defer ts.Close()
	c := buildTestClient(t, ts)

	normalUse := func(t *testing.T) {

		expected := fmt.Sprintf("%s/v1/things/stuff?query=params", ts.URL)
		exampleParams := map[string]string{
			"query": "params",
		}
		actual, err := c.BuildURL(exampleParams, "things", "stuff")

		assert.Nil(t, err)
		assert.Equal(t, expected, actual, "BuildURL doesn't return the correct result. Expected `%s`, got `%s`", expected, actual)
	}

	invalidURL := func(t *testing.T) {
		actual, err := c.BuildURL(nil, `%gh&%ij`)

		assert.NotNil(t, err)
		assert.Empty(t, actual)
	}

	subtests := []subtest{
		{
			Message: "normal use",
			Test:    normalUse,
		},
		{
			Message: "invalid url",
			Test:    invalidURL,
		},
	}
	runSubtestSuite(t, subtests)
}

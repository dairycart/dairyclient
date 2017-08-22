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

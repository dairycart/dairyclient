package dairyclient_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/dairycart/dairyclient"
)

const (
	exampleURL      = `http://www.dairycart.com`
	exampleUsername = `username`
	examplePassword = `password` // lol not really
)

////////////////////////////////////////////////////////
//                                                    //
//               Test Helper Functions                //
//                                                    //
////////////////////////////////////////////////////////

func buildTestCookie() *http.Cookie {
	c := &http.Cookie{Name: "dairycart"}
	return c
}

func buildTestClient(t *testing.T, ts *httptest.Server) *dairyclient.V1Client {
	u, err := url.Parse(ts.URL)
	assert.Nil(t, err)

	c := &dairyclient.V1Client{
		URL:        u,
		Client:     ts.Client(),
		AuthCookie: buildTestCookie(),
	}

	return c
}

func obligatoryLoginHandler(addCookie bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if addCookie {
			cookie := &http.Cookie{
				Name: "dairycart",
			}
			http.SetCookie(w, cookie)
		}
	})
}

func handlerGenerator(handlers map[string]func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, handlerFunc := range handlers {
			if r.URL.Path == path {
				handlerFunc(w, r)
				return
			}
		}
	})
}

////////////////////////////////////////////////////////
//                                                    //
//                 Constructor Tests                  //
//                                                    //
////////////////////////////////////////////////////////

func TestV1ClientConstructor(t *testing.T) {
	ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
	defer ts.Close()

	c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

	assert.NotNil(t, c.AuthCookie)
	assert.Nil(t, err)
}

func TestV1ClientConstructorReturnsErrorWithInvalidURL(t *testing.T) {
	c, err := dairyclient.NewV1Client(":", exampleUsername, examplePassword, nil)

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestV1ClientConstructorWithFailureToLogin(t *testing.T) {
	ts := httptest.NewTLSServer(obligatoryLoginHandler(true))
	ts.Close()

	c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestV1ClientConstructorWhereLoginCookieIsNotReturned(t *testing.T) {
	ts := httptest.NewTLSServer(obligatoryLoginHandler(false))
	defer ts.Close()

	c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

	assert.Nil(t, c)
	assert.NotNil(t, err)
}

func TestDairyClientAddsCookieToRequest(t *testing.T) {
	var endpointCalled bool
	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/product/sku": func(res http.ResponseWriter, req *http.Request) {
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

	res, err := c.ProductExists("sku")

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.True(t, endpointCalled)
}

func TestBuildURL(t *testing.T) {
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
	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()
	c := buildTestClient(t, ts)

	actual, err := c.BuildURL(nil, `%gh&%ij`)

	assert.NotNil(t, err)
	assert.Empty(t, actual)
}

////////////////////////////////////////////////////////
//                                                    //
//               Auth Functions Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestCreateUser(t *testing.T) {
	var endpointCalled bool

	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/user": func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true

			response := `
				{
					"id": 1,
					"first_name": "First",
					"last_name": "Last",
					"email": "email@address.com",
					"is_admin": false
				}
			`

			fmt.Fprintf(res, response)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	exampleInput := dairyclient.UserCreationInput{
		FirstName: "First",
		LastName:  "Last",
		Email:     "email@address.com",
	}
	_, err := c.CreateUser(exampleInput)

	// actual, err := c.CreateUser(exampleInput)
	// expected := &dairyclient.User{
	// 	FirstName: "First",
	// 	LastName:  "Last",
	// 	Email:     "email@address.com",
	// }
	// assert.Equal(t, expected, actual, "expected response did not match actual response.")

	assert.Nil(t, err)
	assert.True(t, endpointCalled)
}

func TestCreateUserReturnsErrorWhenFailingToExecuteRequest(t *testing.T) {
	ts := httptest.NewServer(http.NotFoundHandler())
	defer ts.Close()
	c := buildTestClient(t, ts)
	c.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return errors.New("arbitrary error")
	}

	exampleInput := dairyclient.UserCreationInput{
		FirstName: "First",
		LastName:  "Last",
		Email:     "email@address.com",
	}
	_, err := c.CreateUser(exampleInput)

	// FIXME: I'm not actually testing what I should be testing. :(
	// fmt.Println(err)
	assert.NotNil(t, err)
}

func TestCreateUserReturnsErrorWhenReceivingABadResponse(t *testing.T) {
	var endpointCalled bool

	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/user": func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			badResponse := `
				{
					"id": 1,
				}
			`
			fmt.Fprintf(res, badResponse)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	exampleInput := dairyclient.UserCreationInput{
		FirstName: "First",
		LastName:  "Last",
		Email:     "email@address.com",
	}
	_, err := c.CreateUser(exampleInput)

	assert.NotNil(t, err)
	assert.True(t, endpointCalled)
}

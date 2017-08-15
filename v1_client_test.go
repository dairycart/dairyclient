package dairyclient_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
	var skuEndpointCalled bool

	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/login": func(w http.ResponseWriter, r *http.Request) {
			http.SetCookie(w, buildTestCookie())
		},
		"/v1/product/sku": func(res http.ResponseWriter, req *http.Request) {
			skuEndpointCalled = true
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

	c, err := dairyclient.NewV1Client(ts.URL, exampleUsername, examplePassword, ts.Client())

	res, err := c.CheckProductExistence("sku")

	assert.NotNil(t, res)
	assert.Nil(t, err)
	assert.True(t, skuEndpointCalled)
}

package dairyclient_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dairycart/dairyclient/v1"
)

const (
	exampleID       = 666
	exampleURL      = `http://www.dairycart.com`
	exampleUsername = `username`
	examplePassword = `password` // lol not really
	exampleSKU      = `sku`
	exampleBadJSON  = `{"invalid lol}`
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

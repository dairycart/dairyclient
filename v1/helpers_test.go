package dairyclient_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tdewolff/minify"

	jsonMinify "github.com/tdewolff/minify/json"

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

type subtest struct {
	Message string
	Test    func(t *testing.T)
}

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

func handlerGenerator(handlers map[string]http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for path, handlerFunc := range handlers {
			if r.URL.Path == path {
				handlerFunc(w, r)
				return
			}
		}
	})
}

func minifyJSON(t *testing.T, jsonBody string) string {
	jsonMinifier := minify.New()
	jsonMinifier.AddFunc("application/json", jsonMinify.Minify)
	minified, err := jsonMinifier.String("application/json", jsonBody)
	assert.Nil(t, err)
	return minified
}

func generateHandler(t *testing.T, expectedBody string, expectedMethod string, responseBody string, responseHeader int) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		actualBody, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err)
		assert.Equal(t, minifyJSON(t, expectedBody), string(actualBody), "expected and actual bodies should be equal")

		assert.True(t, req.Method == expectedMethod)
		// additionalFunc(res, req)

		res.WriteHeader(responseHeader)
		fmt.Fprintf(res, responseBody)
	}
}

func generateHeadHandler(t *testing.T, responseHeader int) http.HandlerFunc {
	handler := generateHandler(
		t,
		"",
		http.MethodHead,
		"",
		responseHeader,
	)
	return handler
}

func generateGetHandler(t *testing.T, responseBody string, responseHeader int) http.HandlerFunc {
	handler := generateHandler(
		t,
		"",
		http.MethodGet,
		responseBody,
		responseHeader,
	)
	return handler
}

func generatePostHandler(t *testing.T, expectedBody string, responseBody string, responseHeader int) http.HandlerFunc {
	handler := generateHandler(
		t,
		expectedBody,
		http.MethodPost,
		responseBody,
		responseHeader,
	)
	return handler
}

func generateDeleteHandler(t *testing.T, responseBody string, responseHeader int) http.HandlerFunc {
	handler := generateHandler(
		t,
		"",
		http.MethodDelete,
		responseBody,
		responseHeader,
	)
	return handler
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

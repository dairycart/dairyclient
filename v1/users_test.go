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
//                User Function Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestCreateUser(t *testing.T) {
	t.Parallel()
	var endpointCalled bool

	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/v1/user": func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.True(t, req.Method == http.MethodPost)

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
	actual, err := c.CreateUser(exampleInput)
	expected := &dairyclient.User{
		DBRow:     dairyclient.DBRow{ID: 1},
		FirstName: "First",
		LastName:  "Last",
		Email:     "email@address.com",
	}
	assert.Equal(t, expected, actual, "expected response did not match actual response.")

	assert.Nil(t, err)
	assert.True(t, endpointCalled)
}

func TestCreateUserReturnsErrorWhenFailingToExecuteRequest(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(http.NotFoundHandler())
	c := buildTestClient(t, ts)
	ts.Close()

	exampleInput := dairyclient.UserCreationInput{
		FirstName: "First",
		LastName:  "Last",
		Email:     "email@address.com",
	}
	_, err := c.CreateUser(exampleInput)

	assert.NotNil(t, err)
}

func TestCreateUserReturnsErrorWhenReceivingABadResponse(t *testing.T) {
	t.Parallel()
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

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	var endpointCalled bool

	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		fmt.Sprintf("/v1/user/%d", exampleID): func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.True(t, req.Method == http.MethodDelete)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	err := c.DeleteUser(exampleID)
	assert.Nil(t, err)
	assert.True(t, endpointCalled)
}

func TestDeleteUserWhenErrorEncounteredExecutingRequest(t *testing.T) {
	t.Parallel()

	ts := httptest.NewServer(http.NotFoundHandler())
	c := buildTestClient(t, ts)
	ts.Close()

	err := c.DeleteUser(exampleID)
	assert.NotNil(t, err)
}

func TestDeleteUserWhenResponseContainsError(t *testing.T) {
	t.Parallel()
	var endpointCalled bool

	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		fmt.Sprintf("/v1/user/%d", exampleID): func(res http.ResponseWriter, req *http.Request) {
			endpointCalled = true
			assert.True(t, req.Method == http.MethodDelete)
			res.WriteHeader(http.StatusNotFound)
		},
	}

	ts := httptest.NewServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	err := c.DeleteUser(exampleID)
	assert.NotNil(t, err)
	assert.True(t, endpointCalled)
}

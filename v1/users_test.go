package dairyclient_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairyclient/v1"

	"github.com/stretchr/testify/assert"
)

////////////////////////////////////////////////////////
//                                                    //
//                User Function Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestCreateUser(t *testing.T) {
	t.Parallel()

	expectedBody := `
		{
			"first_name": "First",
			"last_name": "Last",
			"username": "",
			"email": "email@address.com",
			"password": "",
			"is_admin": false
		}
	`
	responseBody := `
		{
			"id": 1,
			"first_name": "First",
			"last_name": "Last",
			"email": "email@address.com",
			"is_admin": false
		}
	`
	h := generatePostHandler(t, expectedBody, responseBody, http.StatusOK)

	handlers := map[string]http.HandlerFunc{
		"/v1/user": h,
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
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
}

func TestCreateUserReturnsErrorWhenFailingToExecuteRequest(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.NotFoundHandler())
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
	expectedBody := `
		{
			"first_name": "First",
			"last_name": "Last",
			"username": "",
			"email": "email@address.com",
			"password": "",
			"is_admin": false
		}
	`
	badResponse := `
		{
			"id": 1,
		}
	`
	handler := generatePostHandler(
		t,
		expectedBody,
		badResponse,
		http.StatusInternalServerError,
	)
	handlers := map[string]http.HandlerFunc{"/v1/user": handler}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	exampleInput := dairyclient.UserCreationInput{
		FirstName: "First",
		LastName:  "Last",
		Email:     "email@address.com",
	}
	_, err := c.CreateUser(exampleInput)

	assert.NotNil(t, err)
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/user/%d", exampleID): generateDeleteHandler(t, "", http.StatusOK),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	err := c.DeleteUser(exampleID)
	assert.Nil(t, err)
}

func TestDeleteUserWhenResponseContainsError(t *testing.T) {
	t.Parallel()

	handlers := map[string]http.HandlerFunc{
		fmt.Sprintf("/v1/user/%d", exampleID): generateDeleteHandler(t, "", http.StatusNotFound),
	}

	ts := httptest.NewTLSServer(handlerGenerator(handlers))
	defer ts.Close()
	c := buildTestClient(t, ts)

	err := c.DeleteUser(exampleID)
	assert.NotNil(t, err)
}

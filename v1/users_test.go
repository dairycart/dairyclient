package dairyclient_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dairycart/dairymodels/v1"

	"github.com/stretchr/testify/assert"
)

////////////////////////////////////////////////////////
//                                                    //
//                User Function Tests                 //
//                                                    //
////////////////////////////////////////////////////////

func TestCreateUser(t *testing.T) {

	t.Run("normal usage", func(*testing.T) {

		expectedBody := `
		{
			"email": "email@address.com",
			"created_on": "0001-01-01T00:00:00Z",
			"archived_on": null,
			"first_name": "First",
			"updated_on": null,
			"id": 0,
			"username": "",
			"password_last_changed_on": null,
			"salt": null,
			"last_name": "Last",
			"is_admin": false,
			"password": ""
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

		handlers := map[string]http.HandlerFunc{
			"/v1/user": generatePostHandler(t, expectedBody, responseBody, http.StatusOK),
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		exampleInput := models.User{
			FirstName: "First",
			LastName:  "Last",
			Email:     "email@address.com",
		}

		expected := &models.User{
			ID:        1,
			FirstName: "First",
			LastName:  "Last",
			Email:     "email@address.com",
		}

		actual, err := c.CreateUser(exampleInput)
		assert.Equal(t, expected, actual, "expected response did not match actual response.")
		assert.Nil(t, err)
	})

	t.Run("with failure to execute request", func(*testing.T) {

		ts := httptest.NewTLSServer(http.NotFoundHandler())
		c := buildTestClient(t, ts)
		ts.Close()

		exampleInput := models.User{
			FirstName: "First",
			LastName:  "Last",
			Email:     "email@address.com",
		}
		_, err := c.CreateUser(exampleInput)

		assert.NotNil(t, err)
	})

	t.Run("with bad response", func(*testing.T) {

		expectedBody := `
			{
				"email": "email@address.com",
				"created_on": "0001-01-01T00:00:00Z",
				"archived_on": null,
				"first_name": "First",
				"updated_on": null,
				"id": 0,
				"username": "",
				"password_last_changed_on": null,
				"salt": null,
				"last_name": "Last",
				"is_admin": false,
				"password": ""
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

		exampleInput := models.User{
			FirstName: "First",
			LastName:  "Last",
			Email:     "email@address.com",
		}

		_, err := c.CreateUser(exampleInput)
		assert.NotNil(t, err)
	})
}

func buildNotFoundUserResponse(userID uint64) string {
	return fmt.Sprintf(`
		{
			"status": 404,
			"message": "The user you were looking for (user ID '%d') does not exist"
		}
	`, userID)
}

func TestDeleteUser(t *testing.T) {

	t.Run("normal usage", func(*testing.T) {

		okID := uint64(1)
		exampleResponse := fmt.Sprintf(`
			{
				"id": %d,
				"first_name": "Fart",
				"last_name": "Zappa",
				"email": "frank@zappa.com",
				"is_admin": false,
				"created_on": "2017-12-10T12:55:21.211807Z",
				"updated_on": "",
				"archived_on": "2017-12-10T12:56:00.322918Z"
			}
		`, okID)

		handlers := map[string]http.HandlerFunc{
			fmt.Sprintf("/v1/user/%d", okID): generateDeleteHandler(t, exampleResponse, http.StatusOK),
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		err := c.DeleteUser(okID)
		assert.Nil(t, err)
	})

	t.Run("when response contains error", func(*testing.T) {

		badID := uint64(2)
		handlers := map[string]http.HandlerFunc{
			fmt.Sprintf("/v1/user/%d", badID): generateDeleteHandler(t, buildNotFoundUserResponse(badID), http.StatusNotFound),
		}

		ts := httptest.NewTLSServer(handlerGenerator(handlers))
		defer ts.Close()
		c := buildTestClient(t, ts)

		err := c.DeleteUser(badID)
		assert.NotNil(t, err)
	})
}

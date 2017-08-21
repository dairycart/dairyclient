package dairyclient

import (
	"fmt"
	"net/http"
)

////////////////////////////////////////////////////////
//                                                    //
//                  User Functions                    //
//                                                    //
////////////////////////////////////////////////////////

// CreateUser takes a UserCreationInput and creates the user in Dairycart
func (dc *V1Client) CreateUser(nu UserCreationInput) (*User, error) {
	u, _ := dc.BuildURL(nil, "user")
	body, _ := createBodyFromStruct(nu)

	req, _ := http.NewRequest(http.MethodPost, u, body)
	res, err := dc.executeRequest(req)
	if err != nil {
		return nil, err
	}

	ru := User{}
	err = unmarshalBody(res, &ru)
	if err != nil {
		return nil, err
	}

	return &ru, nil
}

// DeleteUser deletes a user with a given ID
func (dc *V1Client) DeleteUser(userID uint64) error {
	userIDString := convertIDToString(userID)
	u, _ := dc.BuildURL(nil, "user", userIDString)

	req, _ := http.NewRequest(http.MethodDelete, u, nil)
	res, err := dc.executeRequest(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("user couldn't be deleted, status returned: %d", res.StatusCode)
	}
	return nil
}

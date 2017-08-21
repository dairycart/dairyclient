package dairyclient

import (
	"net/http"
)

////////////////////////////////////////////////////////
//                                                    //
//                  User Functions                    //
//                                                    //
////////////////////////////////////////////////////////

// CreateUser takes a UserCreationInput and creates the user in Dairycart
func (dc *V1Client) CreateUser(nu UserCreationInput) (*User, error) {
	u := dc.buildURL(nil, "user")
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
	u := dc.buildURL(nil, "user", userIDString)
	return dc.delete(u)
}

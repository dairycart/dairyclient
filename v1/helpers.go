package dairyclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/dairycart/dairymodels/v1"
)

type ClientError struct {
	Err     error
	FromAPI *models.ErrorResponse
}

func (ce *ClientError) Error() string {
	if ce.Err != nil {
		return ce.Err.Error()
	} else if ce.FromAPI != nil {
		return ce.FromAPI.Error()
	}

	return ""
}

////////////////////////////////////////////////////////
//                                                    //
//                 Helper Functions                   //
//                                                    //
////////////////////////////////////////////////////////

func mapToQueryValues(in map[string]string) url.Values {
	out := url.Values{}
	for k, v := range in {
		out.Set(k, v)
	}
	return out
}

func interfaceArgIsNotPointerOrNil(i interface{}) error {
	if i == nil {
		return errors.New("unmarshalBody cannot accept nil values")
	}
	isNotPtr := reflect.TypeOf(i).Kind() != reflect.Ptr
	if isNotPtr {
		return errors.New("unmarshalBody can only accept pointers")
	}
	return nil
}

func unmarshalBody(res *http.Response, dest interface{}) *ClientError {
	ce := &ClientError{}

	// These paths should only ever be reached in tests, and should never be encountered by an end user.
	if err := interfaceArgIsNotPointerOrNil(dest); err != nil {
		ce.Err = err
		return ce
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ce.Err = err
		return ce
	}

	err = json.Unmarshal(bodyBytes, &dest)
	if err != nil {
		return &ClientError{Err: err}
	}

	apiErr := &models.ErrorResponse{}
	err = json.Unmarshal(bodyBytes, &apiErr)
	if err != nil {
		return &ClientError{Err: err}
	}
	if apiErr.Status != 0 {
		return &ClientError{FromAPI: apiErr}
	}

	return nil
}

func convertIDToString(id uint64) string {
	return strconv.FormatUint(id, 10)
}

func createBodyFromStruct(in interface{}) (io.Reader, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(out), nil
}

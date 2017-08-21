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
)

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

func unmarshalBody(res *http.Response, dest interface{}) error {
	// These paths should only ever be reached in tests, an should never be encountered by an end user.
	if dest == nil {
		return errors.New("unmarshalBody cannot accept nil values")
	}
	isNotPtr := reflect.TypeOf(dest).Kind() != reflect.Ptr
	if isNotPtr {
		return errors.New("unmarshalBody can only accept pointers")
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bodyBytes, &dest)
	if err != nil {
		return err
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

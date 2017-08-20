package dairyclient

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
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

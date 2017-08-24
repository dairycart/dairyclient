// +build !exported

package dairyclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToQueryValues(t *testing.T) {
	t.Parallel()
	exampleQueryParams := map[string]string{
		"param": "value",
	}

	expected := url.Values{
		"param": []string{"value"},
	}
	actual := mapToQueryValues(exampleQueryParams)

	assert.Equal(t, expected, actual, "expected and actual url values should be equal")
}

type testNormalStruct struct {
	Thing string `json:"thing"`
}

func TestUnmarshalBody(t *testing.T) {
	t.Parallel()
	exampleInput := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(`{"thing":"something"}`)),
	}

	expected := testNormalStruct{Thing: "something"}
	actual := testNormalStruct{}
	err := unmarshalBody(exampleInput, &actual)

	assert.Nil(t, err)
	assert.Equal(t, expected, actual, "expected and actual unmarshaled structs should match")
}

type testFailReader struct{}

func (ft testFailReader) Read([]byte) (int, error) {
	return 0, errors.New("pineapple on pizza")
}

func TestUnmarshalBodyFailsWhenItReceivesNil(t *testing.T) {
	t.Parallel()
	exampleFailureInput := &http.Response{
		Body: ioutil.NopCloser(testFailReader{}),
	}

	err := unmarshalBody(exampleFailureInput, nil)
	assert.NotNil(t, err)
	expected := errors.New("unmarshalBody cannot accept nil values")
	assert.Equal(t, expected, err, "expected error string %s")
}

func TestUnmarshalBodyFailsWhenItReceivesANonPointer(t *testing.T) {
	t.Parallel()
	exampleFailureInput := &http.Response{
		Body: ioutil.NopCloser(testFailReader{}),
	}

	err := unmarshalBody(exampleFailureInput, testNormalStruct{})
	assert.NotNil(t, err)
	expected := errors.New("unmarshalBody can only accept pointers")
	assert.Equal(t, expected, err, "expected error string %s")
}

func TestUnmarshalBodyReturnsReadAllError(t *testing.T) {
	t.Parallel()
	exampleFailureInput := &http.Response{
		Body: ioutil.NopCloser(testFailReader{}),
	}

	err := unmarshalBody(exampleFailureInput, &testNormalStruct{})
	assert.NotNil(t, err)
}

func TestUnmarshalBodyFailsWithInvalidStruct(t *testing.T) {
	t.Parallel()
	exampleInput := &http.Response{
		Body: ioutil.NopCloser(bytes.NewBufferString(`{"invalid_lol}`)),
	}

	actual := testNormalStruct{}
	err := unmarshalBody(exampleInput, &actual)

	assert.NotNil(t, err)
}

func TestConvertIDToString(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		input    uint64
		expected string
	}{
		{
			0,
			"0",
		},
		{
			123,
			"123",
		},
		{
			math.MaxUint64,
			"18446744073709551615",
		},
	}

	for _, tc := range testCases {
		actual := convertIDToString(tc.input)
		assert.Equal(t, tc.expected, actual, "converIDToString failed: expected %s, got %s", tc.expected, actual)
	}
}

func TestCreateBodyFromStruct(t *testing.T) {
	t.Parallel()
	in := testNormalStruct{Thing: "something"}
	_, err := createBodyFromStruct(in)
	assert.Nil(t, err)
}

type testBreakableStruct struct {
	Thing json.Number `json:"thing"`
}

func TestCreateBodyFromStructReturnsErrorWithInvalidInput(t *testing.T) {
	t.Parallel()
	f := &testBreakableStruct{Thing: "dongs"}
	_, err := createBodyFromStruct(f)
	assert.NotNil(t, err)
}

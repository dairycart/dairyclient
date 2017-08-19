package dairyclient

import (
	"encoding/json"
	"errors"
	"math"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapToQueryValues(t *testing.T) {
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
	exampleInput := strings.NewReader(`{"thing":"something"}`)

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

func TestUnmarshalBodyFailsReadingAllFromInputReader(t *testing.T) {
	exampleFailureInput := testFailReader{}
	err := unmarshalBody(exampleFailureInput, nil)
	assert.NotNil(t, err)
}

func TestUnmarshalBodyFailsWithInvalidStruct(t *testing.T) {
	exampleInput := strings.NewReader(`{"invalid_lol}`)

	actual := testNormalStruct{}
	err := unmarshalBody(exampleInput, &actual)

	assert.NotNil(t, err)
}

func TestConvertIDToString(t *testing.T) {
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
	in := testNormalStruct{Thing: "something"}
	_, err := createBodyFromStruct(in)
	assert.Nil(t, err)
}

type testBreakableStruct struct {
	Thing json.Number `json:"thing"`
}

func TestCreateBodyFromStructReturnsErrorWithInvalidInput(t *testing.T) {
	f := &testBreakableStruct{Thing: "dongs"}
	_, err := createBodyFromStruct(f)
	assert.NotNil(t, err)
}

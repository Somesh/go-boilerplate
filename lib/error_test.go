package lib

import (
	"reflect"
	"testing"

)

func TestError(t *testing.T) {
	apiError := APIError{
		MessageError: "Test Error",
	}
	apiError1 := APIError{
		Message: "Test Error",
		Code:    "200",
		Fields:  map[string]string{"1": "2"},
	}
	expectedError := "Test Error"
	expectedError1 := "API: Test Error\nCode: 200\nFields:\n\t1: 2"

	output := apiError.Error()
	if output != expectedError {
		t.Errorf("Expected %s , Got %s", expectedError, output)
	}
	output = apiError1.Error()
	if output != expectedError1 {
		t.Errorf("Expected %s , Got %s", expectedError1, output)
	}
}

func TestNewAPIError(t *testing.T) {
	expected := APIError{
		Message: "Test",
		Code:    "500",
	}
	apiError := NewAPIError("Test", "500", "", 0)
	if !reflect.DeepEqual(expected, apiError) {
		t.Errorf("Invalid Ouput")
	}
}


package apiErrors

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	err := New("TEST_ERROR", "Test error message", 400)
	assert.Equal(t, "TEST_ERROR", err.Code, "Unexpected error code")
	assert.Equal(t, "Test error message", err.Message, "Unexpected error message")
	assert.Equal(t, 400, err.Status, "Unexpected status code")
}

func TestAPIError_Error(t *testing.T) {
	err := New("TEST_ERROR", "Test error message", 400)
	expected := "TEST_ERROR: Test error message"
	assert.Equal(t, expected, err.Error(), "Error string mismatch")
}

func TestAPIError_JSON(t *testing.T) {
	err := New("TEST_ERROR", "Test error message", 400)
	jsonBytes, jsonErr := err.JSON()
	assert.NoError(t, jsonErr, "Unexpected error while marshaling to JSON")

	var unmarshaled map[string]interface{}
	unmarshalErr := json.Unmarshal(jsonBytes, &unmarshaled)
	assert.NoError(t, unmarshalErr, "Unexpected error while unmarshaling JSON")

	assert.Equal(t, "TEST_ERROR", unmarshaled["code"], "Unexpected code in JSON")
	assert.Equal(t, "Test error message", unmarshaled["message"], "Unexpected message in JSON")
	_, exists := unmarshaled["status"]
	assert.False(t, exists, "Status field should not be present in JSON")
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Resource: "user"}
	expected := "resource not found: user"
	assert.Equal(t, expected, err.Error(), "Error string should match expected value")
}

func TestIsNotFound(t *testing.T) {
	assert := assert.New(t)

	notFoundErr := &NotFoundError{Resource: "user"}
	assert.True(IsNotFound(notFoundErr), "IsNotFound should return true for NotFoundError")

	otherErr := New("OTHER_ERROR", "Some other error", 500)
	assert.False(IsNotFound(otherErr), "IsNotFound should return false for non-NotFoundError")
}

func TestNewNotFoundError(t *testing.T) {

	t.Run("Test with custom message", func(t *testing.T) {
		customErr := NewNotFoundError("Custom not found message")
		assert.Equal(t, "NOT_FOUND", customErr.Code, "Expected code NOT_FOUND")
		assert.Equal(t, 404, customErr.Status, "Expected status 404")
		assert.Equal(t, "Custom not found message", customErr.Message, "Expected message 'Custom not found message'")
	})

	t.Run("Test with empty message (default message)", func(t *testing.T) {
		defaultErr := NewNotFoundError("")
		assert.Equal(t, "NOT_FOUND", defaultErr.Code, "Expected code NOT_FOUND")
		assert.Equal(t, 404, defaultErr.Status, "Expected status 404")
		assert.Equal(t, "Resource not found", defaultErr.Message, "Expected message 'Resource not found'")
	})

	t.Run("Test Error() method", func(t *testing.T) {
		customErr := NewNotFoundError("Custom not found message")
		expectedErrString := "NOT_FOUND: Custom not found message"
		assert.Equal(t, expectedErrString, customErr.Error(), "Error string does not match expected value")
	})

	t.Run("Test JSON Marshaling", func(t *testing.T) {
		customErr := NewNotFoundError("Custom not found message")

		_, jsonErr := customErr.JSON()
		assert.NoError(t, jsonErr, "Unexpected error while marshaling to JSON")
	})

	t.Run("Test JSON Marshaling with empty message", func(t *testing.T) {
		customErr := NewNotFoundError("Custom not found message")

		jsonBytes, _ := customErr.JSON()

		var unmarshaled map[string]interface{}
		err := json.Unmarshal(jsonBytes, &unmarshaled)
		assert.NoError(t, err, "Unexpected error while unmarshaling JSON")

		assert.Equal(t, "NOT_FOUND", unmarshaled["code"], "Expected code NOT_FOUND")
		assert.Equal(t, "Custom not found message", unmarshaled["message"], "Expected message 'Custom not found message'")
		_, exists := unmarshaled["status"]
		assert.False(t, exists, "Status field should not be present in JSON")
	})

}

func TestNewGenericError(t *testing.T) {

	t.Run("Test with custom message", func(t *testing.T) {
		customErr := NewGenericError("Custom something went wrong message")
		assert.Equal(t, "INTERNAL_SERVER_ERROR", customErr.Code, "Expected code INTERNAL_SERVER_ERROR")
		assert.Equal(t, 500, customErr.Status, "Expected status 500")
		assert.Equal(t, "Custom something went wrong message", customErr.Message, "Expected message 'Custom something went wrong message'")
	})

	t.Run("Test with empty message (default message)", func(t *testing.T) {
		defaultErr := NewGenericError("")
		assert.Equal(t, "INTERNAL_SERVER_ERROR", defaultErr.Code, "Expected code INTERNAL_SERVER_ERROR")
		assert.Equal(t, 500, defaultErr.Status, "Expected status 500")
		assert.Equal(t, "Something went wrong", defaultErr.Message, "Expected message 'Something went wrong'")
	})

	t.Run("Test Error() method", func(t *testing.T) {
		customErr := NewGenericError("Custom Error")
		expectedErrString := "INTERNAL_SERVER_ERROR: Custom Error"
		assert.Equal(t, expectedErrString, customErr.Error(), "Error string does not match expected value")
	})

	t.Run("Test JSON Marshaling", func(t *testing.T) {
		customErr := NewGenericError("Custom Error")

		jsonBytes, jsonErr := customErr.JSON()
		assert.NoError(t, jsonErr, "Unexpected error while marshaling to JSON")

		var unmarshaled map[string]interface{}
		assert.NoError(t, json.Unmarshal(jsonBytes, &unmarshaled), "Unexpected error while unmarshaling JSON")

		assert.Equal(t, "INTERNAL_SERVER_ERROR", unmarshaled["code"], "Expected code INTERNAL_SERVER_ERROR")
		assert.Equal(t, "Custom Error", unmarshaled["message"], "Expected message 'Custom Error'")
		assert.NotContains(t, unmarshaled, "status", "Status field should not be present in JSON")
	})

	t.Run("Test JSON Marshaling with empty message", func(t *testing.T) {
		customErr := NewGenericError("Custom something went wrong message")

		jsonBytes, _ := customErr.JSON()

		var unmarshaled map[string]interface{}
		if unmarshalErr := json.Unmarshal(jsonBytes, &unmarshaled); unmarshalErr != nil {
			t.Errorf("Unexpected error while unmarshaling JSON: %v", unmarshalErr)
		}
		assert.Equal(t, "INTERNAL_SERVER_ERROR", unmarshaled["code"], "Expected code INTERNAL_SERVER_ERROR")
		assert.Equal(t, "Custom something went wrong message", unmarshaled["message"], "Expected message 'Custom something went wrong message'")
		assert.NotContains(t, unmarshaled, "status", "Status field should not be present in JSON")

	})

}

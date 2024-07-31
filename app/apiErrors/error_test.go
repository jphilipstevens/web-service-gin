package apiErrors

import (
	"encoding/json"
	"testing"
)

func TestNew(t *testing.T) {
	err := New("TEST_ERROR", "Test error message", 400)
	if err.Code != "TEST_ERROR" {
		t.Errorf("Expected code TEST_ERROR, got %s", err.Code)
	}
	if err.Message != "Test error message" {
		t.Errorf("Expected message 'Test error message', got %s", err.Message)
	}
	if err.Status != 400 {
		t.Errorf("Expected status 400, got %d", err.Status)
	}
}

func TestAPIError_Error(t *testing.T) {
	err := New("TEST_ERROR", "Test error message", 400)
	expected := "TEST_ERROR: Test error message"
	if err.Error() != expected {
		t.Errorf("Expected error string %s, got %s", expected, err.Error())
	}
}

func TestAPIError_JSON(t *testing.T) {
	err := New("TEST_ERROR", "Test error message", 400)
	jsonBytes, jsonErr := err.JSON()
	if jsonErr != nil {
		t.Errorf("Unexpected error while marshaling to JSON: %v", jsonErr)
	}

	var unmarshaled map[string]interface{}
	if unmarshalErr := json.Unmarshal(jsonBytes, &unmarshaled); unmarshalErr != nil {
		t.Errorf("Unexpected error while unmarshaling JSON: %v", unmarshalErr)
	}

	if unmarshaled["code"] != "TEST_ERROR" {
		t.Errorf("Expected code TEST_ERROR, got %v", unmarshaled["code"])
	}
	if unmarshaled["message"] != "Test error message" {
		t.Errorf("Expected message 'Test error message', got %v", unmarshaled["message"])
	}
	if _, exists := unmarshaled["status"]; exists {
		t.Errorf("Status field should not be present in JSON")
	}
}

func TestNotFoundError(t *testing.T) {
	err := &NotFoundError{Resource: "user"}
	expected := "resource not found: user"
	if err.Error() != expected {
		t.Errorf("Expected error string %s, got %s", expected, err.Error())
	}
}

func TestIsNotFound(t *testing.T) {
	notFoundErr := &NotFoundError{Resource: "user"}
	if !IsNotFound(notFoundErr) {
		t.Errorf("IsNotFound should return true for NotFoundError")
	}

	otherErr := New("OTHER_ERROR", "Some other error", 500)
	if IsNotFound(otherErr) {
		t.Errorf("IsNotFound should return false for non-NotFoundError")
	}
}

package db

import (
	"database/sql"
	"example/web-service-gin/app/apiErrors"
	"net/http"
)

const (
	DatabaseError            apiErrors.ErrorCode = "database_error"
	NotFoundError            apiErrors.ErrorCode = "not_found"
	ConstraintViolationError apiErrors.ErrorCode = "constraint_violation"
	ConnectionError          apiErrors.ErrorCode = "connection_error"
)

func MapDBError(err *error) *apiErrors.APIError {
	var code apiErrors.ErrorCode
	var message string
	var status int

	switch *err {
	case sql.ErrNoRows:
		code = NotFoundError
		status = http.StatusNotFound
		message = "resource not found"
	default:
		code = DatabaseError
		message = "data retrieval error"
	}
	customErr := apiErrors.New(string(code), message, status)
	return customErr
}

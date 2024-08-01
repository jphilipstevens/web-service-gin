package db

import (
	"database/sql"
	"example/web-service-gin/app/apiErrors"
	"net/http"
)

const (
	DatabaseErrorCode            apiErrors.ErrorCode = "database_error"
	NotFoundErrorCode            apiErrors.ErrorCode = "not_found"
	ConstraintViolationErrorCode apiErrors.ErrorCode = "constraint_violation"
	ConnectionErrorCode          apiErrors.ErrorCode = "connection_error"
)

// TODO: use errors made in error module to not have http status codes here
var NotFoundError = apiErrors.New(string(NotFoundErrorCode), "resource not found", http.StatusNotFound)
var DatabaseError = apiErrors.New(string(DatabaseErrorCode), "data retrieval error", http.StatusInternalServerError)
var ConstraintViolationError = apiErrors.New(string(ConstraintViolationErrorCode), "constraint violation", http.StatusBadRequest)
var ConnectionError = apiErrors.New(string(ConnectionErrorCode), "connection error", http.StatusBadRequest)

func MapDBError(err *error) *apiErrors.APIError {
	var customErr *apiErrors.APIError

	switch *err {
	case sql.ErrNoRows:
		customErr = NotFoundError
	default:
		customErr = DatabaseError
	}
	return customErr
}

package httpx

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	CodeUnauthenticated  ErrorCode = "unauthenticated"
	CodePermissionDenied ErrorCode = "permission_denied"
	CodeNotFound         ErrorCode = "not_found"
	CodeInvalidArgument  ErrorCode = "invalid_argument"
	CodeInternal         ErrorCode = "internal"
)

// APIError is the only error type that reaches the client. Internal causes are
// kept in Cause, which is logged but never serialized.
type APIError struct {
	Code    ErrorCode
	Message string
	Details map[string]string
	Cause   error
}

func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *APIError) Unwrap() error { return e.Cause }

func (e *APIError) HTTPStatus() int {
	switch e.Code {
	case CodeUnauthenticated:
		return http.StatusUnauthorized
	case CodePermissionDenied:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeInvalidArgument:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func Unauthenticated(message string) *APIError {
	return &APIError{Code: CodeUnauthenticated, Message: message}
}

func PermissionDenied(message string) *APIError {
	return &APIError{Code: CodePermissionDenied, Message: message}
}

func NotFound(message string) *APIError {
	return &APIError{Code: CodeNotFound, Message: message}
}

func InvalidArgument(message string, details map[string]string) *APIError {
	return &APIError{Code: CodeInvalidArgument, Message: message, Details: details}
}

func Internal(cause error) *APIError {
	return &APIError{Code: CodeInternal, Message: "internal server error", Cause: cause}
}

func AsAPIError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return Internal(err)
}

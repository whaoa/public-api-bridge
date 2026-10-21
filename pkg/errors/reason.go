package errors

import (
	"fmt"
	"net/http"
)

const (
	ReasonUnknownCode = 0
	ReasonUnknown     = "unknown"
)

type ReasonError struct {
	Code    int    `json:"code"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
	cause   error
}

// Error implements built-in error interface.
func (e *ReasonError) Error() string {
	return fmt.Sprintf("[%s] %s. (%d)", e.Reason, e.Message, e.Code)
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (e *ReasonError) Unwrap() error {
	return e.cause
}

// Is reports whether any error in err's chain matches target.
func (e *ReasonError) Is(err error) bool {
	if v, ok := AsType[*ReasonError](err); ok {
		return v.Code == e.Code && v.Reason == e.Reason
	}
	return false
}

func IsReason(err error) bool {
	_, ok := AsType[*ReasonError](err)
	return ok
}

// Reason returns an error object with given code and message.
func Reason(code int, reason, message string) *ReasonError {
	if reason == "" {
		reason = ReasonUnknown
	}
	return &ReasonError{
		Code:    code,
		Reason:  reason,
		Message: message,
		cause:   nil,
	}
}

// ReasonFromError try to convert an error to *reasonError.
func ReasonFromError(err error) *ReasonError {
	if v, ok := AsType[*ReasonError](err); ok {
		return v
	}
	return &ReasonError{
		Code:    ReasonUnknownCode,
		Reason:  ReasonUnknown,
		Message: err.Error(),
		cause:   err,
	}
}

// BadRequest returns new BadRequest error that is mapped to a 400 response.
func BadRequest(message string) *ReasonError {
	return Reason(http.StatusBadRequest, "bad_request", message)
}

// Forbidden returns new Forbidden error that is mapped to a 403 response.
func Forbidden(message string) *ReasonError {
	return Reason(http.StatusForbidden, "forbidden", message)
}

// NotFound returns new NotFound error that is mapped to a 404 response.
func NotFound(message string) *ReasonError {
	return Reason(http.StatusNotFound, "not_found", message)
}

// Internal returns new Internal error that is mapped to a 500 response.
func Internal(message string) *ReasonError {
	return Reason(http.StatusInternalServerError, "internal", message)
}

// ServiceUnavailable returns new ServiceUnavailable error that is mapped to a 503 response.
func ServiceUnavailable(message string) *ReasonError {
	return Reason(http.StatusServiceUnavailable, "service_unavailable", message)
}

package utils

// HandleError interface

type HandleError interface {
	error
	Code() string
	HTTPStatus() int
}

// CustomError struct

type CustomError struct {
	code       string
	httpStatus int
	message    string
}

func (e *CustomError) Error() string   { return e.message }
func (e *CustomError) Code() string    { return e.code }
func (e *CustomError) HTTPStatus() int { return e.httpStatus }

// NewCustomError creates a new CustomError with code, http status and message.
func NewCustomError(code string, httpStatus int, message string) *CustomError {
	return &CustomError{code: code, httpStatus: httpStatus, message: message}
}

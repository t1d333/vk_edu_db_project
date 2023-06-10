package errors

import (
	"errors"

	"github.com/valyala/fasthttp"
)

var InternalServerError = errors.New("Internal server error")

// DB errors
var InternalDBError = errors.New("Internal DB error")

// User errors
var (
	UserAlreadyExistsError = errors.New("User already exists")
	UserNotFoundError      = errors.New("User not found")
)

var ErrorToStatusCode = map[error]int{
	InternalServerError:    fasthttp.StatusInternalServerError,
	InternalDBError:        fasthttp.StatusInternalServerError,
	UserAlreadyExistsError: fasthttp.StatusConflict,
	UserNotFoundError:      fasthttp.StatusNotFound,
}

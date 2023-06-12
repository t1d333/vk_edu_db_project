package errors

import (
	"errors"

	"github.com/valyala/fasthttp"
)

var (
	InternalServerError = errors.New("Internal server error")
	BadRequstError      = errors.New("Bad request")
)

// DB errors
var InternalDBError = errors.New("Internal DB error")

// User errors
var (
	UserAlreadyExistsError = errors.New("User already exists")
	UserNotFoundError      = errors.New("User not found")
)

// Forum errors
var (
	ForumAlreadyExistsError = errors.New("Forum already exists")
	ForumNotFound           = errors.New("Forum not found")
)

// Thread errors

var ThreadAlreadyExists = errors.New("Thread already exists")

var ErrorToStatusCode = map[error]int{
	InternalServerError:     fasthttp.StatusInternalServerError,
	BadRequstError:          fasthttp.StatusBadRequest,
	InternalDBError:         fasthttp.StatusInternalServerError,
	UserAlreadyExistsError:  fasthttp.StatusConflict,
	UserNotFoundError:       fasthttp.StatusNotFound,
	ForumAlreadyExistsError: fasthttp.StatusConflict,
	ForumNotFound:           fasthttp.StatusNotFound,
	ThreadAlreadyExists:     fasthttp.StatusConflict,
}

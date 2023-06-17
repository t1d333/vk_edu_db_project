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
var (
	ThreadAlreadyExistsError = errors.New("Thread already exists")
	ThreadNotFoundError      = errors.New("Thread not found")
)

// Post errors
var (
	EmptyPostRequestError   = errors.New("Request does not contain any post")
	ParentPostNotFoundError = errors.New("Parent post not found")
	PostNotFoundError       = errors.New("Post not found")
)

// Vote errors
var (
	VoiceArleadyExistsError = errors.New("Voice arleady exists")
	VoiceNotFoundError      = errors.New("Voice not found")
)

var ErrorToStatusCode = map[error]int{
	InternalServerError:      fasthttp.StatusInternalServerError,
	BadRequstError:           fasthttp.StatusBadRequest,
	InternalDBError:          fasthttp.StatusInternalServerError,
	UserAlreadyExistsError:   fasthttp.StatusConflict,
	UserNotFoundError:        fasthttp.StatusNotFound,
	ForumAlreadyExistsError:  fasthttp.StatusConflict,
	ForumNotFound:            fasthttp.StatusNotFound,
	ThreadAlreadyExistsError: fasthttp.StatusConflict,
	ThreadNotFoundError:      fasthttp.StatusNotFound,
	EmptyPostRequestError:    fasthttp.StatusBadRequest,
	ParentPostNotFoundError:  fasthttp.StatusConflict,
	PostNotFoundError:        fasthttp.StatusNotFound,
	VoiceArleadyExistsError:  fasthttp.StatusConflict,
	VoiceNotFoundError:       fasthttp.StatusNotFound,
}

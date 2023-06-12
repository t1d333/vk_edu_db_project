package middleware

import (
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	"github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/valyala/fasthttp"
)

func ErrorMiddlaware(handler routing.Handler) routing.Handler {
	return func(ctx *routing.Context) error {
		ctx.SetContentType("application/json")
		err := handler(ctx)
		if err == nil {
			return nil
		}

		statusCode, ok := errors.ErrorToStatusCode[err]

		if !ok {
			err = errors.InternalServerError
			statusCode = fasthttp.StatusInternalServerError
		}

		response := models.ErrorResponse{Message: err.Error()}
		body, _ := response.MarshalJSON()
		ctx.ResetBody()
		ctx.SetStatusCode(statusCode)
		ctx.SetBody(body)

		return nil
	}
}

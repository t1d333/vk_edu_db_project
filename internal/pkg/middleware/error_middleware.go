package middleware

import (
	"fmt"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	"github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/valyala/fasthttp"
)

func ErrorMiddlaware(handler routing.Handler) routing.Handler {
	return func(ctx *routing.Context) error {
		err := handler(ctx)
		if err == nil {
			return nil
		}

		statusCode, ok := errors.ErrorToStatusCode[err]

		if !ok {
			err = errors.InternalServerError
			statusCode = fasthttp.StatusInternalServerError
		}

		fmt.Println(err, ctx.Request.URI().String())
		response := models.ErrorResponse{Message: err.Error()}
		body, _ := response.MarshalJSON()
		ctx.SetStatusCode(statusCode)
		ctx.SetBody(body)

		return nil
	}
}

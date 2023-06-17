package http

import (
	routing "github.com/qiangxue/fasthttp-routing"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/pkg/middleware"
	"github.com/t1d333/vk_edu_db_project/internal/service"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type delivery struct {
	service service.Service
	logger  *zap.Logger
}

func RegisterHandlers(router *routing.Router, logger *zap.Logger, serv service.Service) {
	del := delivery{serv, logger}
	router.Post("/api/service/clear", middleware.ErrorMiddlaware(del.Clear))
	router.Get("/api/service/status", middleware.ErrorMiddlaware(del.GetStatus))
}

func (del *delivery) Clear(ctx *routing.Context) error {
	if err := del.service.Clear(); err != nil {
		return err
	}
	return nil
}

func (del *delivery) GetStatus(ctx *routing.Context) error {
	status, err := del.service.GetStatus()
	if err != nil {
		return err
	}

	body, err := status.MarshalJSON()
	if err != nil {
		return pkgErrors.InternalServerError
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

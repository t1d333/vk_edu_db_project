package http

import (
	"strconv"
	"strings"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/pkg/middleware"
	"github.com/t1d333/vk_edu_db_project/internal/post"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type delivery struct {
	service post.Service
	logger  *zap.Logger
}

func RegisterHandlers(router *routing.Router, logger *zap.Logger, serv post.Service) {
	del := delivery{serv, logger}
	router.Get("/api/post/<id>/details", middleware.ErrorMiddlaware(del.GetPost))
	router.Post("/api/post/<id>/details", middleware.ErrorMiddlaware(del.UpdatePost))
}

func (del *delivery) GetPost(ctx *routing.Context) error {
	idStr := ctx.Param("id")
	related := string(ctx.QueryArgs().Peek("related"))
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return pkgErrors.BadRequstError
	}

	response, err := del.service.GetPost(id, strings.Split(related, ","))
	if err != nil {
		return err
	}

	body, err := response.MarshalJSON()
	if err != nil {
		return pkgErrors.InternalServerError
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

func (del *delivery) UpdatePost(ctx *routing.Context) error {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return pkgErrors.BadRequstError
	}

	post := models.Post{}
	if err := post.UnmarshalJSON(ctx.PostBody()); err != nil {
		return pkgErrors.BadRequstError
	}

	post.Id = id
	post, err = del.service.UpdatePost(&post)
	if err != nil {
		return err
	}

	body, err := post.MarshalJSON()
	if err != nil {
		return pkgErrors.InternalServerError
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

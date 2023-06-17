package http

import (
	"encoding/json"
	"errors"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/t1d333/vk_edu_db_project/internal/forum"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/pkg/middleware"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type delivery struct {
	serv   forum.Service
	logger *zap.Logger
}

func RegisterHandlers(router *routing.Router, logger *zap.Logger, serv forum.Service) {
	del := delivery{serv, logger}
	router.Post("/api/forum/create", middleware.ErrorMiddlaware(del.Create))
	router.Get("/api/forum/<slug>/details", middleware.ErrorMiddlaware(del.GetForum))
	router.Get("/api/forum/<slug>/users", middleware.ErrorMiddlaware(del.GetUsers))
	router.Get("/api/forum/<slug>/threads", middleware.ErrorMiddlaware(del.GetThreads))
}

func (del *delivery) Create(ctx *routing.Context) error {
	forum := models.Forum{}

	if err := forum.UnmarshalJSON(ctx.PostBody()); err != nil {
		return pkgErrors.BadRequstError
	}

	forum, err := del.serv.Create(&forum)
	if err != nil {
		if !errors.Is(pkgErrors.ForumAlreadyExistsError, err) {
			return err
		}
		ctx.SetStatusCode(fasthttp.StatusConflict)
	} else {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	}

	body, err := forum.MarshalJSON()
	if err != nil {
		return pkgErrors.InternalServerError
	}

	ctx.SetBody(body)
	return nil
}

func (del *delivery) GetForum(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	forum, err := del.serv.GetForum(slug)
	if err != nil {
		return err
	}

	body, err := forum.MarshalJSON()
	if err != nil {
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

func (del *delivery) CreateThread(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	thread := models.Thread{}
	if err := thread.UnmarshalJSON(ctx.PostBody()); err != nil {
		return pkgErrors.BadRequstError
	}

	thread.Forum = slug
	thread, err := del.serv.CreateThread(&thread)

	if err != nil {
		switch {
		case errors.Is(pkgErrors.UserNotFoundError, err):
			return err
		case errors.Is(pkgErrors.ForumNotFound, err):
			return err
		case errors.Is(pkgErrors.ThreadAlreadyExistsError, err):
			ctx.SetStatusCode(fasthttp.StatusConflict)
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	}

	body, err := thread.MarshalJSON()
	if err != nil {
		del.logger.Error("", zap.Error(err))
		return pkgErrors.InternalServerError
	}

	ctx.SetBody(body)
	return nil
}

func (del *delivery) GetUsers(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 100
	}
	sinceTmp := ctx.QueryArgs().Peek("since")
	since := string(sinceTmp)

	desc := ctx.QueryArgs().GetBool("desc")
	users, err := del.serv.GetUsers(slug, limit, since, desc)
	if err != nil {
		return err
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	body, err := users.MarshalJSON()
	if err != nil {
		del.logger.Error("", zap.Error(err))
	}
	ctx.SetBody(body)
	return nil
}

func (del *delivery) GetThreads(ctx *routing.Context) error {
	slug := ctx.Param("slug")
	limit, _ := ctx.QueryArgs().GetUint("limit")
	sinceTmp := ctx.QueryArgs().Peek("since")
	since := string(sinceTmp)
	desc := ctx.QueryArgs().GetBool("desc")
	threads, err := del.serv.GetThreads(slug, limit, since, desc)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(threads)
	ctx.SetBody(body)
	return nil
}

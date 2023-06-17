package http

import (
	"errors"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/pkg/middleware"
	"github.com/t1d333/vk_edu_db_project/internal/thread"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type delivery struct {
	serv   thread.Service
	logger *zap.Logger
}

func RegisterHandlers(router *routing.Router, logger *zap.Logger, serv thread.Service) {
	del := delivery{serv, logger}

	router.Post("/api/forum/<slug>/create", middleware.ErrorMiddlaware(del.CreateThread))
	router.Post("/api/thread/<slug_or_id>/create", middleware.ErrorMiddlaware(del.CreatePost))
	router.Get("/api/thread/<slug_or_id>/details", middleware.ErrorMiddlaware(del.GetThread))
	router.Post("/api/thread/<slug_or_id>/details", middleware.ErrorMiddlaware(del.UpdateThread))
	router.Get("/api/thread/<slug_or_id>/posts", middleware.ErrorMiddlaware(del.GetPosts))
	router.Post("/api/thread/<slug_or_id>/vote", middleware.ErrorMiddlaware(del.AddVote))
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

func (del *delivery) CreatePost(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	var posts models.PostList
	if err := posts.UnmarshalJSON(ctx.PostBody()); err != nil {
		return pkgErrors.BadRequstError
	}

	posts, err := del.serv.CreatePosts(slugOrId, posts)
	if err != nil {
		return err
	}

	body, err := posts.MarshalJSON()
	if err != nil {
		return pkgErrors.InternalServerError
	}
	ctx.SetStatusCode(fasthttp.StatusCreated)
	ctx.SetBody(body)
	return nil
}

func (del *delivery) GetThread(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	thread, err := del.serv.GetThread(slugOrId)
	if err != nil {
		return err
	}
	body, err := thread.MarshalJSON()
	if err != nil {
		del.logger.Error("", zap.Error(err))
		return pkgErrors.InternalServerError
	}
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

func (del *delivery) UpdateThread(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	thread := &models.Thread{}
	if err := thread.UnmarshalJSON(ctx.PostBody()); err != nil {
		return pkgErrors.BadRequstError
	}

	updatedThread, err := del.serv.UpdateThread(slugOrId, thread)
	if err != nil {
		return err
	}

	body, err := updatedThread.MarshalJSON()
	if err != nil {
		del.logger.Error("", zap.Error(err))
		return pkgErrors.InternalServerError
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

func (del *delivery) GetPosts(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	limit, err := ctx.QueryArgs().GetUint("limit")
	if err != nil {
		limit = 100
	}

	sort := ctx.QueryArgs().Peek("sort")
	desc := ctx.QueryArgs().GetBool("desc")

	since, err := ctx.QueryArgs().GetUint("since")
	if err != nil {
		since = 0
	}

	posts, err := del.serv.GetPosts(slugOrId, limit, since, string(sort), desc)
	if err != nil {
		return err
	}

	body, err := posts.MarshalJSON()
	if err != nil {
		return pkgErrors.InternalServerError
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

func (del *delivery) AddVote(ctx *routing.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	vote := models.Vote{}
	if err := vote.UnmarshalJSON(ctx.PostBody()); err != nil {
		return pkgErrors.BadRequstError
	}

	thread, err := del.serv.AddVote(slugOrId, &vote)
	if err != nil {
		return err
	}

	body, err := thread.MarshalJSON()
	if err != nil {
		return pkgErrors.InternalServerError
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

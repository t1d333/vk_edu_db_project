package http

import (
	"errors"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/pkg/middleware"
	"github.com/t1d333/vk_edu_db_project/internal/user"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type delivery struct {
	service user.Service
	logger  *zap.Logger
}

func RegisterHandlers(router *routing.Router, logger *zap.Logger, serv user.Service) {
	del := delivery{serv, logger}
	router.Post("/api/user/<nickname>/create", middleware.ErrorMiddlaware(del.Create))
	router.Post("/api/user/<nickname>/profile", middleware.ErrorMiddlaware(del.Update))
	router.Get("/api/user/<nickname>/profile", middleware.ErrorMiddlaware(del.Get))
}

func (del *delivery) Get(ctx *routing.Context) error {
	nickname := ctx.Param("nickname")
	user, err := del.service.Get(nickname)
	if err != nil {
		return err
	}

	body, err := user.MarshalJSON()
	if err != nil {
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
	return nil
}

func (del *delivery) Create(ctx *routing.Context) error {
	body := ctx.PostBody()
	user := models.User{}
	user.Nickname = ctx.Param("nickname")

	if err := user.UnmarshalJSON(body); err != nil {
		del.logger.Error("", zap.Error(err))
		return pkgErrors.BadRequstError
	}

	users, err := del.service.Create(&user)

	if err != nil {
		if errors.Is(pkgErrors.UserAlreadyExistsError, err) {
			ctx.SetStatusCode(fasthttp.StatusConflict)
		} else {
			return err
		}
	} else {
		ctx.SetStatusCode(fasthttp.StatusCreated)
	}

	if err == nil {
		body, _ = users[0].MarshalJSON()
	} else {
		body, _ = users.MarshalJSON()
	}

	ctx.SetBody(body)
	return nil
}

func (del *delivery) Update(ctx *routing.Context) error {
	body := ctx.PostBody()
	user := models.User{}
	user.Nickname = ctx.Param("nickname")

	if err := user.UnmarshalJSON(body); err != nil {
		del.logger.Error("", zap.Error(err))
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return err
	}

	user, err := del.service.Update(&user)
	if err != nil {
		return err
	}

	ctx.SetStatusCode(fasthttp.StatusOK)

	body, err = user.MarshalJSON()

	if err != nil {
		return err
	}
	ctx.SetBody(body)

	return nil
}

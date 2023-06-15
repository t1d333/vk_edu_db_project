package service

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
	"github.com/t1d333/vk_edu_db_project/internal/user"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	rep    user.Repository
}

func NewService(logger *zap.Logger, rep user.Repository) user.Service {
	return &service{logger, rep}
}

func (serv *service) Create(user *models.User) (models.UserList, error) {
	return serv.rep.Create(user)
}

func (serv *service) Update(user *models.User) (models.User, error) {
	return serv.rep.Update(user)
}

func (serv *service) Get(nickname string) (models.User, error) {
	return serv.rep.Get(nickname)
}

package service

import (
	"github.com/t1d333/vk_edu_db_project/internal/forum"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	"go.uber.org/zap"
)

type service struct {
	rep    forum.Repository
	logger *zap.Logger
}

func NewService(logger *zap.Logger, rep forum.Repository) forum.Service {
	return &service{rep, logger}
}

func (serv *service) Create(forum *models.Forum) (models.Forum, error) {
	return serv.rep.Create(forum)
}

func (serv *service) GetForum(slug string) (models.Forum, error) {
	return serv.rep.GetForum(slug)
}

func (serv *service) GetUsers(slug string, limit int, since string, desc bool) (models.UserList, error) {
	return serv.rep.GetUsers(slug, limit, since, desc)
}

func (serv *service) GetThreads(slug string, limit int, since string, desc bool) ([]models.Thread, error) {
	return serv.rep.GetThreads(slug, limit, since, desc)
}

func (serv *service) CreateThread(thread *models.Thread) (models.Thread, error) {
	return serv.rep.CreateThread(thread)
}

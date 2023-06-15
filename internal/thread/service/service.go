package service

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
	"github.com/t1d333/vk_edu_db_project/internal/thread"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	rep    thread.Repository
}

func NewService(logger *zap.Logger, rep thread.Repository) thread.Service {
	return &service{logger, rep}
}

func (serv *service) CreateThread(thread *models.Thread) (models.Thread, error) {
	return serv.rep.CreateThread(thread)
}

func (serv *service) CreatePosts(slugOrId string, posts []models.Post) (models.PostList, error) {
	return serv.rep.CreatePosts(slugOrId, posts)
}

func (serv *service) GetThread(slugOrId string) (models.Thread, error) {
	return serv.rep.GetThread(slugOrId)
}

func (serv *service) UpdateThread(slugOrId string, thread *models.Thread) (models.Thread, error) {
	return serv.rep.UpdateThread(slugOrId, thread)
}

func (serv *service) GetPosts(slugOrId string) ([]models.Post, error) {
	return serv.rep.GetPosts(slugOrId)
}

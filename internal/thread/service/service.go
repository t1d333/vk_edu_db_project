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

func (serv *service) GetPosts(slugOrId string, limit, since int, sort string, desc bool) (models.PostList, error) {
	switch sort {
	case "tree":
		return serv.rep.GetPostsTree(slugOrId, limit, since, desc)
	case "parent_tree":
		return serv.rep.GetPostsParentTree(slugOrId, limit, since, desc)
	default:
		return serv.rep.GetPostsFlat(slugOrId, limit, since, desc)
	}
}

func (serv *service) AddVote(slugOrId string, vote *models.Vote) (models.Thread, error) {
	thread, err := serv.rep.GetThread(slugOrId)
	if err != nil {
		return thread, err
	}

	if _, err := serv.rep.GetVote(&thread, vote); err == nil {
		return serv.rep.UpdateVote(slugOrId, &thread, vote)
	} else {
		return serv.rep.AddVote(&thread, vote)
	}
}

package service

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
	"github.com/t1d333/vk_edu_db_project/internal/post"
	apimodels "github.com/t1d333/vk_edu_db_project/internal/post/api_models"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	rep    post.Repository
}

func NewService(logger *zap.Logger, rep post.Repository) post.Service {
	return &service{logger, rep}
}

func (serv *service) GetPost(id int, related []string) (apimodels.GetPostResponse, error) {
	post, err := serv.rep.GetPost(id)
	if err != nil {
		return apimodels.GetPostResponse{}, err
	}
	tmp := apimodels.GetPostResponse{
		Post:   &post,
		Author: nil,
		Thread: nil,
		Forum:  nil,
	}

	for _, tag := range related {
		switch tag {
		case "user":
			user, err := serv.rep.GetPostAuthor(tmp.Post)
			if err != nil {
				return tmp, err
			}
			tmp.Author = &user

		case "thread":
			thread, err := serv.rep.GetPostThread(tmp.Post)
			if err != nil {
				return tmp, err
			}
			tmp.Thread = &thread

		case "forum":
			forum, err := serv.rep.GetPostForum(tmp.Post)
			if err != nil {
				return tmp, err
			}
			tmp.Forum = &forum
		}
	}
	return tmp, nil
}

func (serv *service) UpdatePost(post *models.Post) (models.Post, error) {
	return serv.rep.UpdatePost(post)
}

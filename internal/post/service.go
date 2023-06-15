package post

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
	"github.com/t1d333/vk_edu_db_project/internal/post/api_models"
)

type Service interface {
	GetPost(id int, related []string) (apimodels.GetPostResponse, error)
	UpdatePost(post *models.Post) (models.Post, error)
}

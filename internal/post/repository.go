package post

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
)

type Repository interface {
	GetPost(id int) (models.Post, error)
	GetPostAuthor(post *models.Post) (models.User, error)
	GetPostThread(post *models.Post) (models.Thread, error)
	GetPostForum(post *models.Post) (models.Forum, error)
	UpdatePost(post *models.Post) (models.Post, error)
}

package thread

import "github.com/t1d333/vk_edu_db_project/internal/models"

type Service interface {
	CreateThread(thread *models.Thread) (models.Thread, error)
	CreatePosts(slugOrId string, posts []models.Post) (models.PostList, error)
	GetThread(slugOrId string) (models.Thread, error)
	UpdateThread(slugOrId string, thread *models.Thread) (models.Thread, error)
	GetPosts(slugOrId string) ([]models.Post, error)
}

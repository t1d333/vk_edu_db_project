package forum

import "github.com/t1d333/vk_edu_db_project/internal/models"

type Service interface {
	Create(forum *models.Forum) (models.Forum, error)
	GetUsers(slug string, limit int, since string, desc bool) (models.UserList, error)
	GetForum(slug string) (models.Forum, error)
	GetThreads(slug string, limit int, since string, desc bool) ([]models.Thread, error)
	CreateThread(thread *models.Thread) (models.Thread, error)
}

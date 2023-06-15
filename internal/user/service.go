package user

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
)

type Service interface {
	Get(nickname string) (models.User, error)
	Create(user *models.User) (models.UserList, error)
	Update(user *models.User) (models.User, error)
}

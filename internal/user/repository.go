package user

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
)

type Repository interface {
	Create(user *models.User) ([]models.User, error)
	Get(nickname string) (models.User, error)
	Update(user *models.User) (models.User, error)
}

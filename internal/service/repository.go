package service

import "github.com/t1d333/vk_edu_db_project/internal/models"

type Repository interface {
	GetStatus() (models.Status, error)
	Clear() error
}

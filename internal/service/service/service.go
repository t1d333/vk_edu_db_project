package service

import (
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgService "github.com/t1d333/vk_edu_db_project/internal/service"
	"go.uber.org/zap"
)

type service struct {
	logger *zap.Logger
	rep    pkgService.Repository
}

func NewService(logger *zap.Logger, rep pkgService.Repository) pkgService.Service {
	return &service{logger, rep}
}

func (serv *service) GetStatus() (models.Status, error) {
	return serv.rep.GetStatus()
}

func (serv *service) Clear() error {
	return serv.rep.Clear()
}

package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/service"
	"go.uber.org/zap"
)

type repository struct {
	logger *zap.Logger
	conn   *pgxpool.Pool
}

func NewRepository(logger *zap.Logger, conn *pgxpool.Pool) service.Repository {
	return &repository{logger, conn}
}

func (rep *repository) GetStatus() (models.Status, error) {
	tmp := models.Status{}
	row := rep.conn.QueryRow(context.Background(), getStatusCmd)

	if err := row.Scan(&tmp.User, &tmp.Forum, &tmp.Thread, &tmp.Post); err != nil {
		rep.logger.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.InternalDBError
	}
	return tmp, nil
}

func (rep *repository) Clear() error {
	_, err := rep.conn.Exec(context.Background(), clearDb)
	if err != nil {
		rep.logger.Error("DB error", zap.Error(err))
		return pkgErrors.InternalDBError
	}
	return nil
}

package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/user"
	"go.uber.org/zap"
)

type repository struct {
	logger *zap.Logger
	conn   *pgxpool.Pool
}

func NewRepository(logger *zap.Logger, conn *pgxpool.Pool) user.Repository {
	return &repository{logger, conn}
}

func (rep *repository) Create(user *models.User) ([]models.User, error) {
	row := rep.conn.QueryRow(context.Background(), createUserCmd, user.Nickname, user.Fullname, user.About, user.Email)
	if err := row.Scan(&user.Id, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			users := make([]models.User, 0)
			if pgErr.ConstraintName == "users_email_key" || pgErr.ConstraintName == "users_pkey" {
				if user, err := rep.Get(user.Nickname); err == nil {
					users = append(users, user)
				}

				if user, err := rep.GetByEmail(user.Email); err == nil {
					if len(users) > 0 {
						if user.Id != users[0].Id {
							users = append(users, user)
						}
					} else {
						users = append(users, user)
					}
				}
				return users, pkgErrors.UserAlreadyExistsError
			} else {
				rep.logger.Error("Query: "+createUserCmd, zap.Error(err))
				return []models.User{}, pkgErrors.InternalDBError
			}
		}
	}
	return []models.User{*user}, nil
}

func (rep *repository) Update(user *models.User) (models.User, error) {
	row := rep.conn.QueryRow(context.Background(), updateUserCmd, user.Nickname, user.Fullname, user.About, user.Email)

	if err := row.Scan(&user.Id, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(pgx.ErrNoRows, err) {
			return models.User{}, pkgErrors.UserNotFoundError
		}

		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "users_nickname_key" || pgErr.ConstraintName == "users_email_key" {
				user, _ := rep.Get(user.Nickname)
				return user, pkgErrors.UserAlreadyExistsError
			} else {
				rep.logger.Error("Internal DB error", zap.Error(err))
				return *user, pkgErrors.InternalDBError
			}
		}
	}

	return *user, nil
}

func (rep *repository) Get(nickname string) (models.User, error) {
	row := rep.conn.QueryRow(context.Background(), getUserCmd, nickname)
	user := models.User{}
	if err := row.Scan(&user.Id, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return user, pkgErrors.UserNotFoundError
		} else {
			rep.logger.Error("Internal DB error", zap.Error(err))
			return user, pkgErrors.InternalDBError
		}
	}

	return user, nil
}

func (rep *repository) GetByEmail(email string) (models.User, error) {
	row := rep.conn.QueryRow(context.Background(), getUserByEmailCmd, email)
	user := models.User{}
	if err := row.Scan(&user.Id, &user.Nickname, &user.Fullname, &user.About, &user.Email); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return user, pkgErrors.UserNotFoundError
		} else {
			rep.logger.Error("Internal DB error", zap.Error(err))
			return user, pkgErrors.InternalDBError
		}
	}

	return user, nil
}

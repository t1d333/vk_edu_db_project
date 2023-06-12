package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/t1d333/vk_edu_db_project/internal/forum"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"go.uber.org/zap"
)

type repository struct {
	logger *zap.Logger
	conn   *pgx.Conn
}

func NewRepository(logger *zap.Logger, conn *pgx.Conn) forum.Repository {
	return &repository{logger, conn}
}

func (rep *repository) Create(forum *models.Forum) (models.Forum, error) {
	userRow := rep.conn.QueryRow(context.Background(), getForumUserCmd, forum.User)

	if err := userRow.Scan(&forum.User); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return models.Forum{}, pkgErrors.UserNotFoundError
		}
		return models.Forum{}, pkgErrors.InternalDBError
	}

	row := rep.conn.QueryRow(context.Background(), createCmd, forum.Title, forum.User, forum.Slug)

	if err := row.Scan(&forum.Id, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "forums_slug_key" {
				tmp, _ := rep.GetForum(forum.Slug)
				return tmp, pkgErrors.ForumAlreadyExistsError
			}

			if pgErr.ConstraintName == "user_nickname" {
				return *forum, pkgErrors.UserNotFoundError
			}

			rep.logger.Error("Internal DB error", zap.Error(err))
			return *forum, pkgErrors.InternalDBError
		}
	}
	return *forum, nil
}

func (rep *repository) GetForum(slug string) (models.Forum, error) {
	row := rep.conn.QueryRow(context.Background(), getForumCmd, slug)
	forum := models.Forum{}
	if err := row.Scan(&forum.Id, &forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return forum, pkgErrors.ForumNotFound
		} else {
			rep.logger.Error("Internal DB error", zap.Error(err))
			return forum, pkgErrors.InternalDBError
		}
	}
	return forum, nil
}

func (rep *repository) CreateThread(thread *models.Thread) (models.Thread, error) {
	created, _ := time.Parse(time.RFC3339Nano, thread.Created)
	row := rep.conn.QueryRow(context.Background(), createThreadCmd, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, created)
	tmp := models.Thread{}
	var dt pgtype.Date
	err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &dt)
	tmp.Created = dt.Time.Format(time.RFC3339Nano)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.ConstraintName == "threads_forum_fkey" {
				return models.Thread{}, pkgErrors.ForumNotFound
			}

			if pgErr.ConstraintName == "threads_author_fkey" {
				return models.Thread{}, pkgErrors.UserNotFoundError
			}

			return tmp, pkgErrors.InternalDBError
		}
	}

	return tmp, nil
}

func (rep *repository) GetUsers(slug string) ([]models.User, error) {
	return []models.User{}, nil
}

func (rep *repository) GetThreads(slug string, limit int, since string, desc bool) ([]models.Thread, error) {
	getCmd := ""
	var rows pgx.Rows
	var err error

	if _, err := rep.GetForum(slug); err != nil {
		if errors.Is(pkgErrors.ForumNotFound, err) {
			return []models.Thread{}, pkgErrors.ForumNotFound
		} else {
			return []models.Thread{}, pkgErrors.InternalDBError
		}
	}

	if desc {
		if since == "" {
			getCmd = getThreadsDescCmd
		} else {
			getCmd = getThreadsDescWithFilterCmd
		}
	} else {
		if since == "" {
			getCmd = getThreadsAscCmd
		} else {
			getCmd = getThreadsAscWithFilterCmd
		}
	}

	if since == "" {
		rows, err = rep.conn.Query(context.Background(), getCmd, slug, limit)
	} else {
		rows, err = rep.conn.Query(context.Background(), getCmd, slug, limit, since)
	}

	if err != nil {
		rep.logger.Error("DB error", zap.Error(err))
	}

	threads := make([]models.Thread, 0)
	tmp := models.Thread{}
	var dt pgtype.Date

	for rows.Next() {

		if err := rows.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &dt); err != nil {
			rep.logger.Error("DB error", zap.Error(err))
			return threads, pkgErrors.InternalDBError
		}

		tmp.Created = dt.Time.Format(time.RFC3339Nano)
		threads = append(threads, tmp)
	}

	return threads, nil
}

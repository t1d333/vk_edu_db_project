package postgres

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/thread"
	"go.uber.org/zap"
)

type repository struct {
	logger *zap.Logger
	conn   *pgx.Conn
}

func NewRepository(logger *zap.Logger, conn *pgx.Conn) thread.Repository {
	return &repository{logger, conn}
}

func (rep *repository) CreateThread(thread *models.Thread) (models.Thread, error) {
	if thread.Slug != "" {
		tmp, err := rep.GetThread(thread.Slug)
		if err == nil {
			return tmp, pkgErrors.ThreadAlreadyExistsError
		}
	}

	created, _ := time.Parse(time.RFC3339Nano, thread.Created)
	row := rep.conn.QueryRow(context.Background(), createThreadCmd, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, created)
	tmp := models.Thread{}
	var dt pgtype.Date
	err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &dt)
	tmp.Created = dt.Time.Format(time.RFC3339Nano)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "threads_forum_fkey":
				return models.Thread{}, pkgErrors.ForumNotFound
			case "threads_author_fkey":
				return models.Thread{}, pkgErrors.UserNotFoundError
			default:
				return tmp, pkgErrors.InternalDBError
			}
		}
	}

	return tmp, nil
}

func (rep *repository) CreatePosts(slugOrId string, posts []models.Post) ([]models.Post, error) {
	result := make([]models.Post, 0)
	postTmp := models.Post{}
	created := time.Now().Format(time.RFC3339Nano)
	var dt pgtype.Date
	var row pgx.Row
	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return []models.Post{}, err
	}

	for _, post := range posts {
		row = rep.conn.QueryRow(context.Background(), createPostCmd, post.Parent, post.Author, post.Message, thread.Id, thread.Forum, created)
		if err := row.Scan(&postTmp.Id, &postTmp.Parent, &postTmp.Author, &postTmp.Message, &postTmp.IsEdited, &postTmp.Forum, &postTmp.Thread, &dt); err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				if pgErr.Message == "Invalid parent" {
					return []models.Post{}, pkgErrors.ParentPostNotFoundError
				}
				switch pgErr.ConstraintName {
				case "posts_forum_fkey":
					return []models.Post{}, pkgErrors.ForumNotFound
				case "posts_thread_fkey":
					return []models.Post{}, pkgErrors.ThreadNotFoundError
				case "thread_check":
					return []models.Post{}, pkgErrors.ThreadNotFoundError
				case "posts_author_fkey":
					return []models.Post{}, pkgErrors.UserNotFoundError
				}
			}

			rep.logger.Error("DB error", zap.Error(err))
			return []models.Post{}, pkgErrors.InternalDBError
		}
		post.Created = dt.Time.Format(time.RFC3339Nano)
		result = append(result, postTmp)
	}
	return result, nil
}

func (rep *repository) GetThread(slugOrId string) (models.Thread, error) {
	tmp := models.Thread{}
	var row pgx.Row

	if id, err := strconv.Atoi(slugOrId); err == nil {
		row = rep.conn.QueryRow(context.Background(), getThreadByIdCmd, id)
	} else {
		row = rep.conn.QueryRow(context.Background(), getThreadBySlugCmd, slugOrId)
	}

	var dt pgtype.Date
	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &dt); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ThreadNotFoundError
		}
		return tmp, pkgErrors.InternalDBError
	}
	tmp.Created = dt.Time.Format(time.RFC3339Nano)

	return tmp, nil
}

func (rep *repository) UpdateThread(slugOrId string, thread *models.Thread) (models.Thread, error) {
	tmp := models.Thread{}
	var row pgx.Row

	if id, err := strconv.Atoi(slugOrId); err == nil {
		row = rep.conn.QueryRow(context.Background(), updateThreadByIdCmd, id, thread.Message, thread.Title)
	} else {
		row = rep.conn.QueryRow(context.Background(), updateThreadBySlugCmd, slugOrId, thread.Message, thread.Title)
	}

	var dt pgtype.Date

	// RETURNING  id, title, author, forum, message, slug, votes, created;
	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Votes, &dt); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ThreadNotFoundError
		}
		rep.logger.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.InternalDBError
	}
	tmp.Created = dt.Time.Format(time.RFC3339Nano)
	return tmp, nil
}

func (rep *repository) GetPosts(slugOrId string) ([]models.Post, error) {
	return []models.Post{}, nil
}

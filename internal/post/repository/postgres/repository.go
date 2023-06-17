package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/post"
	"go.uber.org/zap"
)

type repository struct {
	logger *zap.Logger
	conn   *pgxpool.Pool
}

func NewRepository(logger *zap.Logger, conn *pgxpool.Pool) post.Repository {
	return &repository{logger, conn}
}

func (rep *repository) GetPost(id int) (models.Post, error) {
	tmp := models.Post{}
	row := rep.conn.QueryRow(context.Background(), getPostById, id)
	if err := row.Scan(&tmp.Id, &tmp.Parent, &tmp.Author, &tmp.Message, &tmp.IsEdited, &tmp.Forum, &tmp.Thread, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.PostNotFoundError
		}
		return tmp, pkgErrors.InternalDBError
	}

	return tmp, nil
}

func (rep *repository) GetPostAuthor(post *models.Post) (models.User, error) {
	tmp := models.User{}
	row := rep.conn.QueryRow(context.Background(), getPostAuthor, post.Author)
	if err := row.Scan(&tmp.Id, &tmp.Nickname, &tmp.Fullname, &tmp.About, &tmp.Email); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.UserNotFoundError
		}
		return tmp, pkgErrors.InternalDBError
	}
	return tmp, nil
}

func (rep *repository) GetPostForum(post *models.Post) (models.Forum, error) {
	tmp := models.Forum{}
	row := rep.conn.QueryRow(context.Background(), getPostForum, post.Forum)
	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.User, &tmp.Slug, &tmp.Posts, &tmp.Threads); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.UserNotFoundError
		}
		return tmp, pkgErrors.InternalDBError
	}
	return tmp, nil
}

func (rep *repository) GetPostThread(post *models.Post) (models.Thread, error) {
	tmp := models.Thread{}
	row := rep.conn.QueryRow(context.Background(), getPostThread, post.Thread)
	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Votes, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.UserNotFoundError
		}
		return tmp, pkgErrors.InternalDBError
	}
	return tmp, nil
}

func (rep *repository) UpdatePost(post *models.Post) (models.Post, error) {
	tmp := models.Post{}

	row := rep.conn.QueryRow(context.Background(), updatePost, post.Id, post.Message)
	if err := row.Scan(&tmp.Id, &tmp.Parent, &tmp.Author, &tmp.Message, &tmp.IsEdited, &tmp.Forum, &tmp.Thread, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.PostNotFoundError
		}
		rep.logger.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.InternalDBError
	}

	return tmp, nil
}

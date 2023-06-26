package postgres

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/t1d333/vk_edu_db_project/internal/models"
	pkgErrors "github.com/t1d333/vk_edu_db_project/internal/pkg/errors"
	"github.com/t1d333/vk_edu_db_project/internal/thread"
	"go.uber.org/zap"
)

type repository struct {
	logger *zap.Logger
	conn   *pgxpool.Pool
}

func NewRepository(logger *zap.Logger, conn *pgxpool.Pool) thread.Repository {
	return &repository{logger, conn}
}

func (rep *repository) CreateThread(thread *models.Thread) (models.Thread, error) {
	if thread.Slug != "" {
		tmp, err := rep.GetThread(thread.Slug)
		if err == nil {
			return tmp, pkgErrors.ThreadAlreadyExistsError
		}
	}

	row := rep.conn.QueryRow(context.Background(), createThreadCmd, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created)
	tmp := models.Thread{}
	err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Created)
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

	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return result, err
	}

	if len(posts) == 0 {
		return result, nil
	}

	postTmp := models.Post{}
	created := time.Unix(0, time.Now().UnixNano()/1e6*1e6)
	cmd := createPostBeginCmd
	args := make([]interface{}, 0, 6*len(posts))

	postTmp.Created = created
	for ind, post := range posts {
		tmpId := 0

		row := rep.conn.QueryRow(context.Background(), checkPostAuthor, post.Author)
		if err := row.Scan(&tmpId); err != nil {
			return result, pkgErrors.UserNotFoundError
		}

		if post.Parent != 0 {
			row = rep.conn.QueryRow(context.Background(), checkPostParent, post.Parent, thread.Id)
			if err := row.Scan(&tmpId); err != nil {
				return result, pkgErrors.ParentPostNotFoundError
			}
		}

		cmd += fmt.Sprintf(" ($%d, $%d, $%d, $%d, $%d, $%d)", 6*ind+1, 6*ind+2, 6*ind+3, 6*ind+4, 6*ind+5, 6*ind+6)
		args = append(args, post.Parent, post.Author, post.Message, thread.Id, thread.Forum, created)
		if ind != len(posts)-1 {
			cmd += ","
		}
	}
	cmd += " RETURNING id, parent, author, message, isEdited, forum, thread;"

	rows, err := rep.conn.Query(context.Background(), cmd, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			rep.logger.Error("TEST", zap.Error(err))
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

	for rows.Next() {
        if err := rows.Scan(&postTmp.Id, &postTmp.Parent, &postTmp.Author, &postTmp.Message, &postTmp.IsEdited, &postTmp.Forum, &postTmp.Thread); err != nil {
			return []models.Post{}, pkgErrors.InternalDBError
		}
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

	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Votes, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ThreadNotFoundError
		}
		return tmp, pkgErrors.InternalDBError
	}

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

	if err := row.Scan(&tmp.Id, &tmp.Title, &tmp.Author, &tmp.Forum, &tmp.Message, &tmp.Slug, &tmp.Votes, &tmp.Created); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.ThreadNotFoundError
		}
		rep.logger.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.InternalDBError
	}

	return tmp, nil
}

func (rep *repository) GetPostsFlat(slugOrId string, limit, since int, desc bool) (models.PostList, error) {
	var rows pgx.Rows
	var err error

	tmp := make([]models.Post, 0)
	post := models.Post{}
	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return []models.Post{}, err
	}

	if desc {
		if since != 0 {
			rows, err = rep.conn.Query(context.Background(), getPostsDescWithSinceCmd, thread.Id, since, limit)
			if err != nil {
				return tmp, pkgErrors.InternalDBError
			}
		} else {
			rows, err = rep.conn.Query(context.Background(), getPostsDescCmd, thread.Id, limit)
			if err != nil {
				return tmp, pkgErrors.InternalDBError
			}
		}
	} else {
		rows, err = rep.conn.Query(context.Background(), getPostsAscCmd, thread.Id, since, limit)
	}

	if err != nil {
		return tmp, pkgErrors.InternalDBError
	}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return tmp, pkgErrors.InternalDBError
		}
		tmp = append(tmp, post)
	}

	return tmp, nil
}

func (rep *repository) GetPostsTree(slugOrId string, limit, since int, desc bool) (models.PostList, error) {
	var rows pgx.Rows
	var err error

	tmp := make([]models.Post, 0)
	post := models.Post{}
	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return []models.Post{}, err
	}

	cmd := ""

	if desc {
		cmd = getPostsTreeDescCmd
	} else {
		cmd = getPostsTreeAscCmd
	}

	if since != 0 {
		switch cmd {
		case getPostsTreeAscCmd:
			cmd = getPostsTreeWithSinceAscCmd
		case getPostsTreeDescCmd:
			cmd = getPostsTreeWithSinceDescCmd
		}
	}

	if since != 0 {
		rows, err = rep.conn.Query(context.Background(), cmd, thread.Id, since, limit)
	} else {
		rows, err = rep.conn.Query(context.Background(), cmd, thread.Id, limit)
	}

	if err != nil {
		return tmp, pkgErrors.InternalDBError
	}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return tmp, pkgErrors.InternalDBError
		}
		tmp = append(tmp, post)
	}

	return tmp, nil
}

func (rep *repository) GetPostsParentTree(slugOrId string, limit, since int, desc bool) (models.PostList, error) {
	var rows pgx.Rows
	var err error

	tmp := make([]models.Post, 0)
	post := models.Post{}
	thread, err := rep.GetThread(slugOrId)
	if err != nil {
		return []models.Post{}, err
	}

	cmd := ""

	if desc {
		cmd = getPostsParentTreeDescCmd
	} else {
		cmd = getPostsParentTreeAscCmd
	}

	if since != 0 {
		switch cmd {
		case getPostsParentTreeAscCmd:
			cmd = getPostsParentTreeWithSinceAscCmd
		case getPostsParentTreeDescCmd:
			cmd = getPostsParentTreeWithSinceDescCmd
		}
	}

	if since != 0 {
		rows, err = rep.conn.Query(context.Background(), cmd, thread.Id, since, limit)
	} else {
		rows, err = rep.conn.Query(context.Background(), cmd, thread.Id, limit)
	}

	if err != nil {
		rep.logger.Error("DB error", zap.Error(err))
		return tmp, pkgErrors.InternalDBError
	}

	for rows.Next() {
		if err := rows.Scan(&post.Id, &post.Parent, &post.Author, &post.Message, &post.IsEdited, &post.Forum, &post.Thread, &post.Created); err != nil {
			return tmp, pkgErrors.InternalDBError
		}
		tmp = append(tmp, post)
	}

	return tmp, nil
}

func (rep *repository) GetVote(thread *models.Thread, vote *models.Vote) (models.Vote, error) {
	tmp := *vote
	row := rep.conn.QueryRow(context.Background(), getVoteCmd, vote.Nickname, thread.Id)
	if err := row.Scan(&tmp.Voice); err != nil {
		if errors.Is(pgx.ErrNoRows, err) {
			return tmp, pkgErrors.VoiceNotFoundError
		}
		return tmp, pkgErrors.InternalDBError
	}
	return tmp, nil
}

func (rep *repository) AddVote(thread *models.Thread, vote *models.Vote) (models.Thread, error) {
	row := rep.conn.QueryRow(context.Background(), addVoteCmd, vote.Nickname, vote.Voice, thread.Id)
	id := 0

	if err := row.Scan(&id); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.ConstraintName {
			case "votes_nickname_thread_key":
				return models.Thread{}, pkgErrors.VoiceArleadyExistsError
			case "votes_nickname_fkey":
				return models.Thread{}, pkgErrors.UserNotFoundError
			default:
				rep.logger.Error("DB error", zap.Error(err))
				return models.Thread{}, pkgErrors.InternalDBError
			}
		}
	}

	thread.Votes += vote.Voice
	return *thread, nil
}

func (rep *repository) UpdateVote(slugOrId string, thread *models.Thread, vote *models.Vote) (models.Thread, error) {
	row := rep.conn.QueryRow(context.Background(), updateVoteCmd, vote.Voice, vote.Nickname, thread.Id)
	id := 0
	if err := row.Scan(&id); err != nil {
		if !errors.Is(pgx.ErrNoRows, err) {
			rep.logger.Error("DB error", zap.Error(err))
			return models.Thread{}, pkgErrors.InternalDBError
		}
	}

	return rep.GetThread(slugOrId)
}

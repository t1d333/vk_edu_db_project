package thread

import "github.com/t1d333/vk_edu_db_project/internal/models"

type Repository interface {
	CreateThread(thread *models.Thread) (models.Thread, error)
	CreatePosts(slugOrId string, posts []models.Post) ([]models.Post, error)
	GetThread(slugOrId string) (models.Thread, error)
	UpdateThread(slugOrId string, thread *models.Thread) (models.Thread, error)
	GetPostsFlat(slugOrId string, limit, since int, desc bool) (models.PostList, error)
	GetPostsTree(slugOrId string, limit, since int, desc bool) (models.PostList, error)
	GetPostsParentTree(slugOrId string, limit, since int, desc bool) (models.PostList, error)
	AddVote(thread *models.Thread, vote *models.Vote) (models.Thread, error)
	GetVote(thread *models.Thread, vote *models.Vote) (models.Vote, error)
	UpdateVote(slugOrId string, thread *models.Thread, vote *models.Vote) (models.Thread, error)
}

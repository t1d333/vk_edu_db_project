package apimodels

import "github.com/t1d333/vk_edu_db_project/internal/models"

type GetPostParams struct {
	Author bool
	Forum  bool
	Thread bool
}

type GetPostResponse struct {
	Post   *models.Post   `json:"post"`
	Author *models.User   `json:"author,omitempty"`
	Forum  *models.Forum  `json:"forum,omitempty"`
	Thread *models.Thread `json:"thread,omitempty"`
}

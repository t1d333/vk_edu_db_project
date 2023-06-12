package models

type Forum struct {
	Id      int    `json:"-"`
	Title   string `json:"title"`
	Slug    string `json:"slug"`
	User    string `json:"user"`
	Threads int    `json:"threads"`
	Posts   int    `json:"posts"`
}

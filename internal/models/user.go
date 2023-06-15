package models

//easyjson:json
type UserList []User

//easyjson:json
type User struct {
	Id       int    `json:"-"`
	Nickname string `json:"nickname"`
	Fullname string `json:"fullname"`
	About    string `json:"about"`
	Email    string `json:"email"`
}

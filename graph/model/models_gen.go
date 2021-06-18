// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Email struct {
	Email string `json:"email"`
}

type NewAutoLogin struct {
	Token string `json:"token"`
}

type NewLogin struct {
	Username string `json:"username"`
	Pass     string `json:"pass"`
}

type Paging struct {
	Lastid int `json:"lastid"`
}

type Status struct {
	Status bool `json:"status"`
}

type Token struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

type User struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Avatarurl string `json:"avatarurl"`
	Email     string `json:"email"`
}

type Video struct {
	ID        string `json:"id"`
	Thumbnail string `json:"thumbnail"`
	Link      string `json:"link"`
	Title     string `json:"title"`
	User      *User  `json:"user"`
}

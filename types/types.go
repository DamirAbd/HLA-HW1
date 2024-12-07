package types

import (
	"time"
)

type User struct {
	ID         string    `json:"id" validate:"required"`
	FirstName  string    `json:"first_name" validate:"required"`
	SecondName string    `json:"second_name" validate:"required"`
	BirthDate  string    `json:"birthdate" validate:"required"`
	Biography  string    `json:"biography" validate:"required"`
	City       string    `json:"city" validate:"required"`
	CreatedAt  time.Time `json:"createdAt"`
	Password   string    `json:"-"`
}

type UserStore interface {
	GetUserByID(ID string) (*User, error)
	CreateUser(User) error
	GetUsersByName(FirstName string, SecondName string) ([]*UserForm, error)
	SetFriend(ID string, IDFriend string) error
	DeleteFriend(ID string, IDFriend string) error
	GetFriends(ID string) ([]*UserForm, error)
}

type RegisterUserPayload struct {
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	BirthDate  string `json:"birthdate" validate:"required"`
	Biography  string `json:"biography" validate:"required"`
	City       string `json:"city" validate:"required"`
	Password   string `json:"password" validate:"required,min=3,max=130"`
}

type LoginUserPayload struct {
	ID       string `json:"ID" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type UserForm struct {
	ID         string `json:"ID" validate:"required"`
	FirstName  string `json:"first_name" validate:"required"`
	SecondName string `json:"second_name" validate:"required"`
	BirthDate  string `json:"birthdate" validate:"required"`
	Biography  string `json:"biography" validate:"required"`
	City       string `json:"city" validate:"required"`
}

type Post struct {
	ID      string `json:"ID"`
	AutorId string `json:"authorID"`
	Post    string `json:"post"`
}

type PostStore interface {
	CreatePost(Post) error
	DeletePost(ID string) error
	UpdatePost(ID string, Post string) error
	GetPostByID(id string) (*Post, error)
	GetPostsByUsers(ids []string) ([]*Post, error)
}

type PostPayload struct {
	Post string `json:"text"`
}

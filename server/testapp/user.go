package main

import (
	"errors"
	"github.com/bobappleyard/anathema/server/a"
	"github.com/bobappleyard/anathema/server/hterror"
)

var ErrNoUser = hterror.WithStatusCode(404, errors.New("user not found"))

type User struct {
	ID int `json:"id"`
}

type userRepository interface {
	GetUser(id int) (*userEntity, error)
}

type userResource struct {
	a.Resource `path:"/user/{id}"`

	ID   int `path:"id"`
	Repo userRepository
}

func (userResource) Init(group a.Group) {
	group.GET(userResource.GetUser)
}

func (u userResource) GetUser() (User, error) {
	user, err := u.Repo.GetUser(u.ID)
	if err != nil {
		return User{}, err
	}
	if user == nil {
		return User{}, ErrNoUser
	}
	return User{ID: user.ID}, nil
}

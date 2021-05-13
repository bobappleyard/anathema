package main

import "github.com/bobappleyard/anathema/server/a"

type repository struct {
	a.Service

	users map[int] *userEntity
}

type userEntity struct {
	ID int
}

func (r *repository) Inject() {
	r.users = map[int]*userEntity{
		1: &userEntity{
			ID: 1,
		},
	}
}

func (r *repository) GetUser(id int) (*userEntity, error) {
	return r.users[id], nil
}




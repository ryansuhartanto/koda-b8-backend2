package repo

import (
	"errors"
	"slices"

	"github.com/ryansuhartanto/koda-b8-backend1/model"
)

type RepoUser struct {
	data []model.User
}

func NewRepoUser(data []model.User) *RepoUser {
	return &RepoUser{data}
}

func (r *RepoUser) Create(new model.User) error {
	if index := slices.IndexFunc(r.data, func(user model.User) bool {
		return user.Email == new.Email
	}); index >= 0 {
		return errors.New("Email already exists")
	}

	r.data = append(r.data, new)

	return nil
}

func (r *RepoUser) FindAll() ([]model.User, error) {
	return slices.Clone(r.data), nil
}

func (r *RepoUser) Auth(email string, password model.Password) (*model.User, error) {
	index := slices.IndexFunc(r.data, func(user model.User) bool {
		return user.Email == email
	})

	if index < 0 {
		return nil, errors.New("Email is not registered")
	}

	user := r.data[index]

	if user.Password != password {
		return nil, errors.New("Wrong password")
	}

	return &user, nil
}

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

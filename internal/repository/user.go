package repository

import (
	"slices"

	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
)

type UserRepository struct {
	data []model.User
}

func NewUserRepository(data []model.User) *UserRepository {
	return &UserRepository{data}
}

func (r *UserRepository) Add(new model.User) {
	r.data = append(r.data, new)
}

func (r *UserRepository) FindAll() []model.User {
	return slices.Clone(r.data)
}

func (r *UserRepository) Find(email string) *model.User {
	index := slices.IndexFunc(r.data, func(user model.User) bool {
		return user.Email == email
	})

	if index < 0 {
		return nil
	}

	user := r.data[index]
	return &user
}

package service

import (
	"errors"

	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{repository}
}

func (s *UserService) List() []model.User {
	return s.repository.FindAll()
}

func (s *UserService) Register(new model.User) error {
	if user := s.repository.Find(new.Email); user != nil {
		return errors.New("Email already exists")
	}

	s.repository.Add(new)
	return nil
}

func (r *UserService) Auth(email string, password model.Password) (*model.User, error) {
	user := r.repository.Find(email)

	if user == nil {
		return nil, errors.New("Email is not registered")
	}

	if user.Password != password {
		return nil, errors.New("Wrong password")
	}

	return user, nil
}

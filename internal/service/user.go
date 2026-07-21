package service

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository
	ctx        context.Context
}

func NewUserService(
	repository *repository.UserRepository,
	ctx context.Context,
) *UserService {
	return &UserService{
		repository,
		ctx,
	}
}

func (s *UserService) List() ([]model.User, error) {
	return s.repository.FindAll(s.ctx)
}

var ErrEmailConflict = errors.New("service: email already exists")

func (s *UserService) Register(new model.User) error {
	user, err := s.repository.Find(s.ctx, new.Auth.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}

	if user != nil {
		return ErrEmailConflict
	}

	err = s.repository.Add(s.ctx, new)
	if err != nil {
		return err
	}

	return nil
}

var ErrEmailUnregistered = errors.New("service: email is not registered")
var ErrPasswordIncorrect = errors.New("service: password is incorrect")

func (r *UserService) Login(email model.Email, password model.Password) (*model.User, error) {
	user, err := r.repository.Find(r.ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrEmailUnregistered
	}

	if user.Password != password {
		return nil, ErrPasswordIncorrect
	}

	return user, nil
}

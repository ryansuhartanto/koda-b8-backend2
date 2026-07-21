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

func (s *UserService) List() ([]model.UserIdentified, error) {
	return s.repository.FindAll(s.ctx)
}

var ErrEmailConflict = errors.New("service: email already exists")

func (s *UserService) Register(new model.User) (*model.UserIdentified, error) {
	user, err := s.repository.FindEmail(s.ctx, new.Auth.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if user != nil {
		return nil, ErrEmailConflict
	}

	user, err = s.repository.Add(s.ctx, new)
	if err != nil {
		return nil, err
	}

	return user, nil
}

var ErrEmailUnregistered = errors.New("service: email is not registered")
var ErrPasswordIncorrect = errors.New("service: password is incorrect")

func (r *UserService) Login(auth model.Auth) (*model.UserIdentified, error) {
	user, err := r.repository.FindEmail(r.ctx, auth.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrEmailUnregistered
	}

	if user.User.Password != auth.Password {
		return nil, ErrPasswordIncorrect
	}

	return user, nil
}

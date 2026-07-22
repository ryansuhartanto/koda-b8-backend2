package service

import (
	"context"
	"encoding/base64"
	"errors"

	"golang.org/x/crypto/bcrypt"

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
	user, err := s.repository.FindEmail(s.ctx, new.Credentials.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if user != nil {
		return nil, ErrEmailConflict
	}

	rawPassword, err := base64.StdEncoding.DecodeString(string(new.Password))
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(rawPassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	new.Password = model.Password(encryptedPassword)

	user, err = s.repository.Add(s.ctx, new)
	if err != nil {
		return nil, err
	}

	return user, nil
}

var ErrEmailUnregistered = errors.New("service: email is not registered")
var ErrPasswordIncorrect = errors.New("service: password is incorrect")

func (s *UserService) Login(cre model.Credentials) (*model.UserIdentified, error) {
	user, err := s.repository.FindEmail(s.ctx, cre.Email)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if user == nil {
		return nil, ErrEmailUnregistered
	}

	rawPassword, err := base64.StdEncoding.DecodeString(string(cre.Password))
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), rawPassword); err != nil {
		return nil, errors.Join(ErrPasswordIncorrect, err)
	}

	return user, nil
}

func (s *UserService) Edit(id model.Id, new model.User) (*model.UserIdentified, error) {
	user, err := s.repository.Update(s.ctx, id, new)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Delete(id model.Id) error {
	err := s.repository.Delete(s.ctx, id)
	if err != nil {
		return nil
	}

	return nil
}

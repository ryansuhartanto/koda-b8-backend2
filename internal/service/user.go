package service

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/repository"
)

type UserService struct {
	repository *repository.UserRepository

	jwtKey []byte
}

func NewUserService(
	repository *repository.UserRepository,
	jwtKey []byte,
) *UserService {
	return &UserService{
		repository,
		jwtKey,
	}
}

func (s *UserService) List(ctx context.Context) ([]model.UserIdentified, error) {
	return s.repository.FindAll(ctx)
}

var ErrEmailConflict = errors.New("service: email already exists")

func (s *UserService) Register(ctx context.Context, new model.User) (*model.AuthResult, error) {
	user, err := s.repository.FindEmail(ctx, new.Credentials.Email)
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

	user, err = s.repository.Add(ctx, new)
	if err != nil {
		return nil, err
	}

	result, err := model.NewAuthResult(user, s.jwtKey)
	if err != nil {
		return nil, err
	}

	return result, nil
}

var ErrEmailUnregistered = errors.New("service: email is not registered")
var ErrPasswordIncorrect = errors.New("service: password is incorrect")

func (s *UserService) Login(ctx context.Context, cre model.Credentials) (*model.AuthResult, error) {
	user, err := s.repository.FindEmail(ctx, cre.Email)
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

	result, err := model.NewAuthResult(user, s.jwtKey)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *UserService) Edit(ctx context.Context, id model.Id, new model.User) (*model.UserIdentified, error) {
	rawPassword, err := base64.StdEncoding.DecodeString(string(new.Password))
	if err != nil {
		return nil, err
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(rawPassword, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	new.Password = model.Password(encryptedPassword)

	user, err := s.repository.Update(ctx, id, new)
	if err != nil {
		return nil, err
	}

	return user, nil
}

var ErrImageUnsupported = errors.New("service: unsupported image format")

func (s *UserService) UpdatePicture(ctx context.Context, id model.Id, data []byte) error {
	if data == nil {
		if err := s.repository.UpdatePicture(ctx, id, nil); err != nil {
			return err
		}

		return nil
	}

	contentType := http.DetectContentType(data)
	if !strings.HasPrefix(contentType, "image/") {
		return ErrImageUnsupported
	}

	ext, err := mime.ExtensionsByType(contentType)
	if err != nil {
		return err
	}
	if len(ext) == 0 {
		return ErrImageUnsupported
	}

	file := filepath.Join("uploads", fmt.Sprintf("user-picture-%v%v", id.Id, ext[0]))
	if err := os.WriteFile(file, data, 0644); err != nil {
		return err
	}

	if err := s.repository.UpdatePicture(ctx, id, &file); err != nil {
		return err
	}

	return nil
}

func (s *UserService) Delete(ctx context.Context, id model.Id) error {
	err := s.repository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

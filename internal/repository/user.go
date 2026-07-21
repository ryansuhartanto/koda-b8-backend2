package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/db"
	"github.com/ryansuhartanto/koda-b8-backend1/internal/model"
)

type UserRepository struct {
	querier db.Querier
}

func NewUserRepository(querier db.Querier) *UserRepository {
	return &UserRepository{querier}
}

func (r *UserRepository) Add(ctx context.Context, new model.User) (*model.UserIdentified, error) {
	sql := "INSERT INTO users (name, email, password) VALUES (@name, @email, @password) RETURNING *"
	args := pgx.StructArgs(new)
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fn := pgx.RowToAddrOfStructByName[model.UserIdentified]
	user, err := pgx.CollectOneRow(rows, fn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]model.UserIdentified, error) {
	sql := "SELECT * FROM users"
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fn := pgx.RowToStructByName[model.UserIdentified]
	users, err := pgx.CollectRows(rows, fn)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) FindEmail(ctx context.Context, email model.Email) (*model.UserIdentified, error) {
	sql := "SELECT * FROM users WHERE email = @email"
	args := pgx.NamedArgs{
		"email": email,
	}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fn := pgx.RowToAddrOfStructByName[model.UserIdentified]
	user, err := pgx.CollectOneRow(rows, fn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Find(ctx context.Context, id model.Id) (*model.UserIdentified, error) {
	sql := "SELECT * FROM users WHERE id = @id"
	args := pgx.NamedArgs{
		"id": id,
	}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fn := pgx.RowToAddrOfStructByName[model.UserIdentified]
	user, err := pgx.CollectOneRow(rows, fn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

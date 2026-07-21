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

func (r *UserRepository) Add(ctx context.Context, new model.User) error {
	sql := "INSERT INTO users (name, email, password) VALUES (@name, @email, @password)"
	args := pgx.StructArgs(new)
	_, err := r.querier.Exec(ctx, sql, args)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindAll(ctx context.Context) ([]model.User, error) {
	sql := "SELECT * FROM users"
	rows, err := r.querier.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fn := pgx.RowToStructByName[model.User]
	users, err := pgx.CollectRows(rows, fn)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Find(ctx context.Context, email string) (*model.User, error) {
	sql := "SELECT * FROM users WHERE email = @email"
	args := pgx.NamedArgs{
		"email": email,
	}
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fn := pgx.RowToAddrOfStructByName[model.User]
	user, err := pgx.CollectOneRow(rows, fn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

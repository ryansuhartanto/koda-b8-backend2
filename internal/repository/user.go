package repository

import (
	"context"
	"fmt"

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
	args := StrictFlattenArgs(new)
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
	args := pgx.StrictNamedArgs{
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
	args := StrictFlattenArgs(id)
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

func (r *UserRepository) Update(ctx context.Context, id model.Id, new model.User) (*model.UserIdentified, error) {
	sql := "UPDATE users SET name = @name, email = @email, password = @password WHERE id = @id RETURNING *"
	args := StrictFlattenArgs(model.UserIdentified{Id: id, User: new})
	rows, err := r.querier.Query(ctx, sql, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	fmt.Println(sql, args)

	fn := pgx.RowToAddrOfStructByName[model.UserIdentified]
	user, err := pgx.CollectOneRow(rows, fn)
	if err != nil {
		return nil, err
	}

	return user, nil
}

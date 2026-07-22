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
	sql := `
		WITH new_user AS (
			INSERT INTO users (email, password) VALUES (@email, @password) RETURNING *
		), new_profile AS (
			INSERT INTO profiles (id, name, picture_url)
			VALUES ((SELECT id FROM new_user), @name, @picture_url)
			RETURNING *
		)
		SELECT * FROM new_user JOIN new_profile USING (id)
	`
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
	sql := `
		SELECT * FROM users JOIN profiles USING (id)
	`
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
	sql := `
		SELECT * FROM users JOIN profiles USING (id)
		WHERE users.email = @email
	`
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
	sql := `
		SELECT * FROM users JOIN profiles USING (id)
		WHERE users.id = @id
	`
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
	sql := `
		WITH updated_user AS (
			UPDATE users SET email = @email, password = @password WHERE id = @id RETURNING *
		), updated_profile AS (
			UPDATE profiles SET name = @name, picture_url = @picture_url WHERE id = @id RETURNING *
		)
		SELECT * FROM updated_user JOIN updated_profile USING (id)
	`
	args := StrictFlattenArgs(model.UserIdentified{Id: id, User: new})
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

func (r *UserRepository) Delete(ctx context.Context, id model.Id) error {
	sql := "DELETE FROM users WHERE id = @id"
	args := StrictFlattenArgs(id)
	_, err := r.querier.Exec(ctx, sql, args)
	if err != nil {
		return err
	}

	return nil
}

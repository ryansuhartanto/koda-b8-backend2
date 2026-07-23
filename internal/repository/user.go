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
		SELECT
			new_user.*,
			new_profile.name,
			new_profile.picture_url,
			new_profile.updated_at AS profile_updated_at
		FROM new_user JOIN new_profile USING (id)
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

func (r *UserRepository) FindAll(ctx context.Context, p model.Pagination) ([]model.UserIdentified, error) {
	sql := `
		SELECT
			users.*,
			profiles.name,
			profiles.picture_url,
			profiles.updated_at AS profile_updated_at
		FROM users JOIN profiles USING (id)
		ORDER BY users.id
		LIMIT @limit OFFSET @offset
	`
	args := StrictFlattenArgs(p)
	rows, err := r.querier.Query(ctx, sql, args)
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
		SELECT
			users.*,
			profiles.name,
			profiles.picture_url,
			profiles.updated_at AS profile_updated_at
		FROM users JOIN profiles USING (id)
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
		SELECT
			users.*,
			profiles.name,
			profiles.picture_url,
			profiles.updated_at AS profile_updated_at
		FROM users JOIN profiles USING (id)
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
			UPDATE users SET
				email = @email,
				password = @password
			WHERE id = @id
			RETURNING *
		), updated_profile AS (
			UPDATE profiles SET
				name = @name,
				picture_url = COALESCE(@picture_url,
				picture_url) WHERE id = @id
			RETURNING *
		)
		SELECT
			updated_user.*,
			updated_profile.name,
			updated_profile.picture_url,
			updated_profile.updated_at AS profile_updated_at
		FROM updated_user JOIN updated_profile USING (id)
	`
	args := StrictFlattenArgs(struct {
		model.Id
		model.User
	}{id, new})
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

func (r *UserRepository) UpdatePicture(ctx context.Context, id model.Id, url *string) error {
	sql := `UPDATE profiles SET picture_url = @picture_url WHERE id = @id`

	args := pgx.StrictNamedArgs{
		"id":          id.Id,
		"picture_url": url,
	}
	_, err := r.querier.Exec(ctx, sql, args)
	if err != nil {
		return err
	}

	return nil
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

package repository

import (
	"context"
	"errors"
	"fmt"

	"gitgub.com/Alksndr9/go-students-disciplines/internal/domain"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/infrastructure/db/pgxutils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrUserExists   = errors.New("user exists")
)

type Repo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *Repo {
	return &Repo{db: db}
}

func (r *Repo) SaveUser(ctx context.Context, user *domain.User) error {
	const op = "repository.psql.SaveUser"

	const query = `
			INSERT INTO users (
				username,
				password_hash,
				email,
				phone_number,
				role
			)
			VALUES (
				@username,
				@password_hash,
				@email,
				@phone_number,
				@role
			)`

	args := pgxutils.UserArgs(user)

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, ErrUserExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) GetUserByID(ctx context.Context, id uint64) (*domain.User, error) {
	const op = "repository.psql.GetUserByID"

	const query = `
			SELECT * FROM users
			WHERE id = $1`

	row, err := r.db.Query(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	res, err := pgx.CollectOneRow(
		row,
		pgx.RowToStructByName[domain.User],
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &res, nil
}

func (r *Repo) UpdateUser(ctx context.Context, id uint64, user *domain.User) error {
	const op = "repository.psql.UpdateUser"

	const query = `
			UPDATE users
			SET 
				username = @username,
				password_hash = @password_hash,
				email = @email,
				phone_number = @phone_number,
				role = @role
			WHERE id = @id`

	args := pgxutils.UserArgs(user)
	args["id"] = id

	_, err := r.db.Exec(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repo) DeleteUser(ctx context.Context, id uint64) error {
	const op = "repository.psql.DeleteUser"

	const query = `
			DELETE FROM users
			WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

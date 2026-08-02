package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/light-group/light-ex-backend/internal/modules/auth/domain"
	authrepo "github.com/light-group/light-ex-backend/internal/modules/auth/repository"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// DBTX is implemented by both *sql.DB and *sql.Tx.
type DBTX interface {
	ExecContext(
		ctx context.Context,
		query string,
		args ...any,
	) (sql.Result, error)

	QueryRowContext(
		ctx context.Context,
		query string,
		args ...any,
	) *sql.Row
}

type Repository struct {
	db DBTX
}

// Compile-time interface check.
var _ authrepo.Repository = (*Repository)(nil)

// New creates a PostgreSQL authentication repository.
func New(db DBTX) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(
	ctx context.Context,
	user *domain.User,
) error {

	const query = `
INSERT INTO users (
	id,
	email,
	password_hash,
	email_verified,
	status,
	created_at,
	updated_at
)
VALUES ($1,$2,$3,$4,$5,$6,$7)
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.EmailVerified,
		user.Status,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *Repository) Update(
	ctx context.Context,
	user *domain.User,
) error {

	const query = `
UPDATE users
SET
	email=$2,
	password_hash=$3,
	email_verified=$4,
	status=$5,
	updated_at=$6
WHERE id=$1
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.EmailVerified,
		user.Status,
		user.UpdatedAt,
	)

	return err
}

func (r *Repository) FindByID(
	ctx context.Context,
	id uuid.UUID,
) (*domain.User, error) {

	const query = `
SELECT
	id,
	email,
	password_hash,
	email_verified,
	status,
	created_at,
	updated_at
FROM users
WHERE id=$1
`

	return r.scanUser(
		r.db.QueryRowContext(ctx, query, id),
	)
}

func (r *Repository) FindByEmail(
	ctx context.Context,
	email string,
) (*domain.User, error) {

	const query = `
SELECT
	id,
	email,
	password_hash,
	email_verified,
	status,
	created_at,
	updated_at
FROM users
WHERE email=$1
`

	return r.scanUser(
		r.db.QueryRowContext(ctx, query, email),
	)
}

func (r *Repository) ExistsByEmail(
	ctx context.Context,
	email string,
) (bool, error) {

	const query = `
SELECT EXISTS(
	SELECT 1
	FROM users
	WHERE email=$1
)
`

	var exists bool

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(&exists)

	return exists, err
}

func (r *Repository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {

	const query = `
DELETE FROM users
WHERE id=$1
`

	_, err := r.db.ExecContext(
		ctx,
		query,
		id,
	)

	return err
}

func (r *Repository) scanUser(
	row *sql.Row,
) (*domain.User, error) {

	var user domain.User

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.EmailVerified,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}

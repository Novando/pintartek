package repository

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/domain/user/entity"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUser struct {
	ctx   context.Context
	query *pgx.Queries
	db    *pgxpool.Pool
}

func NewPostgresUserRepository(q *pgx.Queries, db *pgxpool.Pool) *PostgresUser {
	return &PostgresUser{
		query: q,
		db:    db,
	}
}

const createPostgresUser = `-- name: Create user :one
	INSERT INTO users (email, password, public_key, created_at, updated_at)
	VALUES ($1::varchar, $2::varchar, $3::text, NOW(), NOW())
	RETURNING id
`

func (r *PostgresUser) Create(arg CreateParam) (id pgtype.UUID, err error) {
	row := r.db.QueryRow(r.ctx, createPostgresUser, arg.Email, arg.Password, arg.PrivateKey)
	err = row.Scan(&id)
	return
}

const getPostgresUserByID = `-- name: Get user by the ID :one
	SELECT id, email, password, public_key, created_at, updated_at, deleted_at
	FROM users
	WHERE id = $1::uuid AND deleted_at IS NOT NULL
`

func (r *PostgresUser) GetByID(id pgtype.UUID) (data entity.User, err error) {
	row := r.db.QueryRow(r.ctx, getPostgresUserByID, id)
	err = row.Scan(
		&data.ID,
		&data.Email,
		&data.Password,
		&data.PublicKey,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)
	return
}

const updatePasswordPostgresUser = `-- name: Update password of a user :exec
	UPDATE users SET password = $1::varchar, updated_at = NOW() WHERE id = $2::uuid
`

func (r *PostgresUser) UpdatePassword(id pgtype.UUID, password string) error {
	_, err := r.db.Exec(r.ctx, updatePasswordPostgresUser, password, id)
	return err
}

const updatePrivateKeyPostgresUser = `-- name: Update private key of a user :exec
	UPDATE users SET public_key = $1::text, updated_at = NOW() WHERE id = $2::uuid
`

func (r *PostgresUser) UpdatePrivateKey(id pgtype.UUID, pk string) error {
	_, err := r.db.Exec(r.ctx, updatePrivateKeyPostgresUser, pk, id)
	return err
}

const deletePostgresUser = `-- name: Soft delete a user :exec
	UPDATE users SET deleted_at = NOW() WHERE id = $1::uuid
`

func (r *PostgresUser) Delete(id pgtype.UUID) error {
	_, err := r.db.Exec(r.ctx, deletePostgresUser, id)
	return err
}

const permanentDeletePostgresUser = `-- name: Permanent delete a user :exec
	DELETE FROM users WHERE id = $1::uuid
`

func (r *PostgresUser) PermanentDelete(id pgtype.UUID) error {
	_, err := r.db.Exec(r.ctx, permanentDeletePostgresUser, id)
	return err
}

package repository

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/domain/user/entity"
	"github.com/Novando/pintartek/pkg/common/consts"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUser struct {
	ctx   context.Context
	query *pgx.Queries
	db    *pgxpool.Pool
}

func NewPostgresUserRepository(
	c context.Context,
	q *pgx.Queries,
	db *pgxpool.Pool,
) *PostgresUser {
	return &PostgresUser{
		ctx:   c,
		query: q,
		db:    db,
	}
}

const createPostgresUser = `-- name: Create user :one
	INSERT INTO users (id, email, password, public_key, access_token, backup_token, created_at, updated_at)
	VALUES ($1::uuid, $2::varchar, $3::varchar, $4::varchar, $5::varchar, $6::varchar, NOW(), NOW())
	RETURNING id
`

func (r *PostgresUser) Create(arg CreateParam) (id pgtype.UUID, err error) {
	row := r.db.QueryRow(r.ctx, createPostgresUser,
		arg.ID,
		arg.Email,
		arg.Password,
		arg.PublicKey,
		arg.AccessToken,
		arg.BackupToken,
	)
	err = row.Scan(&id)
	return
}

const getPostgresUserByID = `-- name: Get user by the ID :one
	SELECT id, email, password, public_key, access_token, backup_token, created_at, updated_at, deleted_at
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
		&data.AccessToken,
		&data.BackupToken,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)
	return
}

const getPostgresUserByEmail = `-- name: Get user by an email :one
	SELECT id, email, password, public_key, access_token, backup_token, created_at, updated_at, deleted_at
	FROM users
	WHERE email = $1::varchar AND deleted_at IS NOT NULL
`

func (r *PostgresUser) GetByEmail(email string) (data entity.User, err error) {
	row := r.db.QueryRow(r.ctx, getPostgresUserByEmail, email)
	err = row.Scan(
		&data.ID,
		&data.Email,
		&data.Password,
		&data.PublicKey,
		&data.AccessToken,
		&data.BackupToken,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)
	if err != nil && err.Error() == pgx.ErrNoRows() {
		err = consts.ErrNoData
	}
	return
}

const updatePasswordPostgresUser = `-- name: Update password of a user :exec
	UPDATE users SET password = $1::varchar, updated_at = NOW() WHERE id = $2::uuid
`

func (r *PostgresUser) UpdatePassword(id pgtype.UUID, password string) error {
	_, err := r.db.Exec(r.ctx, updatePasswordPostgresUser, password, id)
	return err
}

const updatePublicKeyPostgresUser = `-- name: Update private key of a user :exec
	UPDATE users SET public_key = $1::text, updated_at = NOW() WHERE id = $2::uuid
`

func (r *PostgresUser) UpdatePublicKey(id pgtype.UUID, pk string) error {
	_, err := r.db.Exec(r.ctx, updatePublicKeyPostgresUser, pk, id)
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

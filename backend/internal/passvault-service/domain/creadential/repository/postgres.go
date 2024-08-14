package repository

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/domain/client/entity"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresCredential struct {
	ctx   context.Context
	query *pgx.Queries
	db    *pgxpool.Pool
}

func NewPostgresCredentialRepository(q *pgx.Queries, db *pgxpool.Pool) *PostgresCredential {
	return &PostgresCredential{
		query: q,
		db:    db,
	}
}

const createPostgresCredential = `-- name: Create credential :exec
	INSERT INTO credential (full_name, created_at, updated_at)
	VALUES ($1::varchar, NOW(), NOW())
	RETURNING id
`

func (r *PostgresCredential) Create(name string) (id pgtype.UUID, err error) {
	row := r.db.QueryRow(r.ctx, createPostgresCredential, name)
	err = row.Scan(&id)
	return
}

const getPostgresCredentialByID = `-- name: Get client by the ID :one
	SELECT id, full_name, created_at, updated_at, deleted_at
	FROM users
	WHERE id = $1::uuid AND deleted_at IS NOT NULL
`

func (r *PostgresCredential) GetByID(id pgtype.UUID) (data entity.Client, err error) {
	row := r.db.QueryRow(r.ctx, getPostgresCredentialByID, id)
	err = row.Scan(
		&data.ID,
		&data.FullName,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)
	return
}

const updatePostgresCredential = `-- name: Update client data :exec
	UPDATE clients SET full_name = $1::varchar, updated_at = NOW() WHERE id = $2::uuid
`

func (r *PostgresCredential) Update(id pgtype.UUID, name string) error {
	_, err := r.db.Exec(r.ctx, updatePostgresCredential, name, id)
	return err
}

const permanentDeletePostgresCredential = `-- name: Permanent delete a client :exec
	DELETE FROM clients WHERE id = $1::uuid
`

func (r *PostgresCredential) PermanentDelete(id pgtype.UUID) error {
	_, err := r.db.Exec(r.ctx, permanentDeletePostgresCredential, id)
	return err
}

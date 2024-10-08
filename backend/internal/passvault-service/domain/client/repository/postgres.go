package repository

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/domain/client/entity"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresClient struct {
	ctx   context.Context
	query *pgx.Queries
	db    *pgxpool.Pool
}

func NewPostgresClientRepository(c context.Context, q *pgx.Queries, db *pgxpool.Pool) *PostgresClient {
	return &PostgresClient{
		ctx:   c,
		query: q,
		db:    db,
	}
}

const createPostgresClient = `-- name: Create client :exec
	INSERT INTO clients (user_id, full_name, created_at, updated_at)
	VALUES ($1::uuid, $2::varchar, NOW(), NOW())
	RETURNING id
`

func (r *PostgresClient) Create(name string, userId pgtype.UUID) (id pgtype.UUID, err error) {
	row := r.db.QueryRow(r.ctx, createPostgresClient, userId, name)
	err = row.Scan(&id)
	return
}

const getByIDPostgresClient = `-- name: Get client by the ID :one
	SELECT id, user_id, full_name, created_at, updated_at, deleted_at
	FROM clients
	WHERE id = $1::uuid AND deleted_at IS NULL
`

func (r *PostgresClient) GetByID(id pgtype.UUID) (data entity.Client, err error) {
	row := r.db.QueryRow(r.ctx, getByIDPostgresClient, id)
	err = row.Scan(
		&data.ID,
		&data.UserID,
		&data.FullName,
		&data.CreatedAt,
		&data.UpdatedAt,
		&data.DeletedAt,
	)
	return
}

const updatePostgresClient = `-- name: Update client data :exec
	UPDATE clients SET full_name = $1::varchar, updated_at = NOW() WHERE id = $2::uuid
`

func (r *PostgresClient) Update(id pgtype.UUID, name string) error {
	_, err := r.db.Exec(r.ctx, updatePostgresClient, name, id)
	return err
}

const deletePostgresClient = `-- name: Soft delete a client :exec
	UPDATE clients SET deleted_at = NOW() WHERE id = $1::uuid
`

func (r *PostgresClient) Delete(id pgtype.UUID) error {
	_, err := r.db.Exec(r.ctx, deletePostgresClient, id)
	return err
}

const permanentDeletePostgresClient = `-- name: Permanent delete a client :exec
	DELETE FROM clients WHERE id = $1::uuid
`

func (r *PostgresClient) PermanentDelete(id pgtype.UUID) error {
	_, err := r.db.Exec(r.ctx, permanentDeletePostgresClient, id)
	return err
}

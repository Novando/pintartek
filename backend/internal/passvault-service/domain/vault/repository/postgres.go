package repository

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/domain/vault/entity"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresVault struct {
	ctx   context.Context
	query *pgx.Queries
	db    *pgxpool.Pool
}

func NewPostgresVaultRepository(
	c context.Context,
	q *pgx.Queries,
	db *pgxpool.Pool,
) *PostgresVault {
	return &PostgresVault{
		ctx:   c,
		query: q,
		db:    db,
	}
}

const createPostgresVault = `-- name: Create vault :one
	INSERT INTO vaults(name, credential, created_at, updated_at)
	VALUES ($1::uuid, $2::varchar, $3::varchar, NOW(), NOW())
	RETURNING id
`

func (r *PostgresVault) Create(arg CreateParam) (id pgtype.UUID, err error) {
	row := r.db.QueryRow(r.ctx, createPostgresVault,
		arg.PivotID,
		arg.Name,
		arg.Credential,
	)
	err = row.Scan(&id)
	return
}

const getByIDPostgresVault = `-- name: Get vault by the ID :one
	SELECT id, name, credential, created_at, updated_at
	FROM vaults
	WHERE id = $1::uuid AND deleted_at IS NULL
`

func (r *PostgresVault) GetByID(id pgtype.UUID) (data entity.Vault, err error) {
	row := r.db.QueryRow(r.ctx, getByIDPostgresVault, id)
	err = row.Scan(
		&data.ID,
		&data.Name,
		&data.Credential,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	return
}

const updatePostgresVault = `-- name: Update vault data :exec
	UPDATE vaults SET
		name = $1::varchar,
		credential = $2::varchar,
		updated_at = NOW()
	WHERE id = $3::uuid
`

func (r *PostgresVault) Update(id pgtype.UUID, arg UpdateParam) error {
	_, err := r.db.Exec(r.ctx, updatePostgresVault,
		arg.Name,
		arg.Credential,
		id,
	)
	return err
}

const permanentDeletePostgresVault = `-- name: Permanent delete a vault :exec
	DELETE FROM vaults WHERE id = $1::uuid
`

func (r *PostgresVault) PermanentDelete(id pgtype.UUID) error {
	_, err := r.db.Exec(r.ctx, permanentDeletePostgresVault, id)
	return err
}

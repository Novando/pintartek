package repository

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/domain/vault-group/aggregate"
	"github.com/Novando/pintartek/pkg/common/structs"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresVaultGroup struct {
	ctx   context.Context
	query *pgx.Queries
	db    *pgxpool.Pool
}

func NewPostgresVaultGroupRepository(
	c context.Context,
	q *pgx.Queries,
	db *pgxpool.Pool,
) *PostgresVaultGroup {
	return &PostgresVaultGroup{
		ctx:   c,
		query: q,
		db:    db,
	}
}

const createPostgresVaultGroup = `-- name: Create user-vault pivot relation :exec
	INSERT INTO user_vault_pivots(user_id, vault_id)
	VALUES ($1::uuid, $2::uuid)
`

func (r *PostgresVaultGroup) Create(arg CreateParam) error {
	_, err := r.db.Exec(r.ctx, createPostgresVaultGroup, arg.UserID, arg.VaultID)
	return err
}

const permanentDeletePostgresVaultGroup = `-- name: Permanent delete a user-vault pivot relation :exec
	DELETE FROM user_vault_pivots WHERE id = $1::int
`

func (r *PostgresVaultGroup) PermanentDelete(id uint64) error {
	_, err := r.db.Exec(r.ctx, permanentDeletePostgresVaultGroup, id)
	return err
}

const getAllVaultByUserIDPostgresVaultGroup = `-- name: Get all vault by user ID :many
	SELECT
		v.id AS id,
		u.id AS user_id,
		name, 
		credential,
		v.created_at AS created_at,
		v.updated_at AS updated_at
	FROM vaults v
	LEFT JOIN user_vault_pivots uvp ON v.id = uvp.vault_id
	LEFT JOIN users u ON uvp.user_id = u.id
	WHERE u.id = $1::uuid
	LIMIT $2::int OFFSET $3::int
`

func (r *PostgresVaultGroup) GetAllVaultByUserID(
	userID pgtype.UUID,
	arg structs.StdPagination,
) (data []aggregate.VaultList, err error) {
	rows, err := r.db.Query(r.ctx, getAllVaultByUserIDPostgresVaultGroup,
		userID,
		arg.Size,
		arg.Page,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i aggregate.VaultList
		if err = rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Name,
			&i.Credential,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		data = append(data, i)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return
}

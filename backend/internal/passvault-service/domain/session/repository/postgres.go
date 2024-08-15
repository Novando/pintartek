package repository

import (
	"context"
	"github.com/Novando/pintartek/internal/passvault-service/domain/session/entity"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type PostgresSession struct {
	ctx   context.Context
	query *pgx.Queries
	db    *pgxpool.Pool
}

func NewPostgresSessionRepository(
	c context.Context,
	q *pgx.Queries,
	db *pgxpool.Pool,
) *PostgresSession {
	return &PostgresSession{
		ctx:   c,
		query: q,
		db:    db,
	}
}

const createPostgresSession = `-- name: Create session :one
	INSERT INTO sessions (id, user_id, secret_key, expired_at)
	VALUES ($1::uuid, $2::uuid, $3::varchar, $4)
	RETURNING id
`

func (r *PostgresSession) Create(arg CreateParam) (id pgtype.UUID, err error) {
	expiry := time.Now().Add(time.Minute * 30)
	row := r.db.QueryRow(r.ctx, createPostgresSession,
		arg.ID,
		arg.UserID,
		arg.SecretKey,
		pgtype.Timestamptz{Time: expiry, Valid: true},
	)
	err = row.Scan(&id)
	return
}

const getPostgresSessionByID = `-- name: Get session by the ID :one
	SELECT user_id, secret_key
	FROM sessions
	WHERE id = $1::uuid AND expired_at >= NOW()
`

func (r *PostgresSession) GetByID(id pgtype.UUID) (data entity.Session, err error) {
	row := r.db.QueryRow(r.ctx, getPostgresSessionByID, id)
	err = row.Scan(
		&data.UserID,
		&data.SecretKey,
	)
	return
}

const permanentDeletePostgresSession = `-- name: Permanent delete a session :exec
	DELETE FROM sessions WHERE id = $1::uuid
`

func (r *PostgresSession) PermanentDelete(id pgtype.UUID) error {
	_, err := r.db.Exec(r.ctx, permanentDeletePostgresSession, id)
	return err
}

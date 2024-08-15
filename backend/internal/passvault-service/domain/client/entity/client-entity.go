package entity

import "github.com/jackc/pgx/v5/pgtype"

type Client struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	DeletedAt pgtype.Timestamptz
	FullName  string
}

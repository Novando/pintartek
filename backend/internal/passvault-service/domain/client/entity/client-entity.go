package entity

import "github.com/jackc/pgx/v5/pgtype"

type Client struct {
	ID        pgtype.UUID
	OwnerID   pgtype.UUID
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
	DeletedAt pgtype.Timestamptz
	FullName  string
}

package entity

import "github.com/jackc/pgx/v5/pgtype"

type Credential struct {
	ID         pgtype.UUID
	OwnerID    pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
	Credential string
	Password   string
	Url        string
	Note       string
}

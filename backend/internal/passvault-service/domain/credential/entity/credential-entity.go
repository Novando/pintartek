package entity

import "github.com/jackc/pgx/v5/pgtype"

type Credential struct {
	OwnerID    pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
	Name       string
	Credential string
	Password   string
	Url        string
	Note       string
}

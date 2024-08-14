package entity

import "github.com/jackc/pgx/v5/pgtype"

type Vault struct {
	ID         pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
	Credential string
	Name       string
}

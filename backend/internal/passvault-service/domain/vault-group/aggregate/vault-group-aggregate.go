package aggregate

import "github.com/jackc/pgx/v5/pgtype"

type VaultList struct {
	ID         pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
	Name       string
	Credential string
	UserID     uint64
}

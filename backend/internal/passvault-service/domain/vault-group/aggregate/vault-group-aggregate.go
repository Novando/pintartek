package aggregate

import "github.com/jackc/pgx/v5/pgtype"

type VaultList struct {
	ID         pgtype.UUID
	UserID     pgtype.UUID
	CreatedAt  pgtype.Timestamptz
	UpdatedAt  pgtype.Timestamptz
	Name       string
	Credential string
}

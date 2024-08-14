package entity

import "github.com/jackc/pgx/v5/pgtype"

type UserVault struct {
	UserID  pgtype.UUID
	VaultID pgtype.UUID
	ID      uint64
}

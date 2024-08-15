package entity

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	ID          pgtype.UUID
	CreatedAt   pgtype.Timestamptz
	UpdatedAt   pgtype.Timestamptz
	DeletedAt   pgtype.Timestamptz
	Email       string
	Password    string
	PublicKey   string
	AccessToken string
	BackupToken string
}

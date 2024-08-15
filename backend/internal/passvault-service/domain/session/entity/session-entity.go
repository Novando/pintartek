package entity

import "github.com/jackc/pgx/v5/pgtype"

type Session struct {
	UserID    pgtype.UUID `json:"userId"`
	SecretKey string      `json:"secretKey"`
}

package repository

import (
	"github.com/Novando/pintartek/internal/passvault-service/domain/session/entity"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateParam struct {
	ID        pgtype.UUID
	UserID    pgtype.UUID
	SecretKey string
}

type Session interface {
	Create(CreateParam) (pgtype.UUID, error)
	GetByID(pgtype.UUID) (entity.Session, error)
	PermanentDelete(pgtype.UUID) error
}

package repository

import (
	"github.com/Novando/pintartek/internal/passvault-service/domain/user/entity"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	CreateParam struct {
		Email      string
		Password   string
		PrivateKey string
	}
)

type User interface {
	Create(arg CreateParam) (id pgtype.UUID, err error)
	GetByID(id pgtype.UUID) (data entity.User, err error)
	UpdatePassword(id pgtype.UUID, password string) error
	UpdatePrivateKey(id pgtype.UUID, pub string) error
	Delete(id pgtype.UUID) error
	PermanentDelete(id pgtype.UUID) error
}

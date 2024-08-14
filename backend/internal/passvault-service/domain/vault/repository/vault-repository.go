package repository

import (
	"github.com/Novando/pintartek/internal/passvault-service/domain/vault/entity"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	CreateParam struct {
		PivotID pgtype.UUID
		UpdateParam
	}
	UpdateParam struct {
		Credential string
		Name       string
	}
)

type Vault interface {
	Create(arg CreateParam) (id pgtype.UUID, err error)
	GetByID(id pgtype.UUID) (data entity.Vault, err error)
	Update(id pgtype.UUID, arg UpdateParam) error
	PermanentDelete(id pgtype.UUID) error
}

package repository

import (
	"github.com/Novando/pintartek/internal/passvault-service/domain/vault/entity"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpsertParam struct {
	Credential string
	Name       string
}

type Vault interface {
	Create(arg UpsertParam) (id pgtype.UUID, err error)
	GetByID(id pgtype.UUID) (data entity.Vault, err error)
	UpdateName(id pgtype.UUID, name string) error
	UpdateCredential(id pgtype.UUID, credential string) error
	PermanentDelete(id pgtype.UUID) error
}

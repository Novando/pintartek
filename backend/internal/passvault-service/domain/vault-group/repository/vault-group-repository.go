package repository

import (
	"github.com/Novando/pintartek/internal/passvault-service/domain/vault-group/aggregate"
	"github.com/Novando/pintartek/pkg/common/structs"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateParam struct {
	UserID  pgtype.UUID
	VaultID pgtype.UUID
}

type Vault interface {
	Create(arg CreateParam) error
	PermanentDelete(id uint64) error
	GetAllVaultByUserID(userID pgtype.UUID, arg structs.StdPagination) ([]aggregate.VaultList, error)
}

package repository

import (
	"github.com/Novando/pintartek/internal/passvault-service/domain/creadential/entity"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateParam struct {
	Url        pgtype.Text
	Note       pgtype.Text
	Credential string
	Password   string
}

type Credential interface {
	Create(arg CreateParam) (id pgtype.UUID, err error)
	GetByID(id pgtype.UUID) (data entity.Credential, err error)
	Update(id pgtype.UUID, name string) error
	PermanentDelete(id pgtype.UUID) error
}

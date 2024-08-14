package repository

import (
	"github.com/Novando/pintartek/internal/passvault-service/domain/client/entity"
	"github.com/jackc/pgx/v5/pgtype"
)

type Client interface {
	Create(name string) (id pgtype.UUID, err error)
	GetByID(id pgtype.UUID) (data entity.Client, err error)
	Update(id pgtype.UUID, name string) error
	Delete(id pgtype.UUID) error
	PermanentDelete(id pgtype.UUID) error
}

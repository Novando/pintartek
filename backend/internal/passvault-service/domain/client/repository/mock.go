package repository

import (
	"github.com/stretchr/testify/mock"
)

type PostgresMock struct {
	Mock mock.Mock
}

//func (r *PostgresMock) Create(name string, userId pgtype.UUID) (id pgtype.UUID, err error) {
//	r.Mock.Called()
//}

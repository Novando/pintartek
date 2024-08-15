package service

import (
	"context"
	"fmt"
	"github.com/Novando/pintartek/internal/passvault-service/app/dto"
	clientRepo "github.com/Novando/pintartek/internal/passvault-service/domain/client/repository"
	userRepo "github.com/Novando/pintartek/internal/passvault-service/domain/user/repository"
	"github.com/Novando/pintartek/pkg/common/consts"
	"github.com/Novando/pintartek/pkg/common/structs"
	"github.com/Novando/pintartek/pkg/crypto"
	"github.com/Novando/pintartek/pkg/helper"
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/Novando/pintartek/pkg/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserConfig func(su *UserService)

type UserService struct {
	log        *logger.Logger
	userRepo   userRepo.User
	clientRepo clientRepo.Client
}

// NewUserService Initialize user service
func NewUserService(config UserConfig, cfgs ...UserConfig) *UserService {
	serv := &UserService{}
	cfgs = append(cfgs, config)
	for _, cfg := range cfgs {
		cfg(serv)
	}
	return serv
}

// WithPostgres Using Postgres to store data
func WithPostgres(l *logger.Logger, q *pgx.Queries, db *pgxpool.Pool) UserConfig {
	return func(su *UserService) {
		su.log = l
		su.userRepo = userRepo.NewPostgresUserRepository(context.Background(), q, db)
		su.clientRepo = clientRepo.NewPostgresClientRepository(context.Background(), q, db)
	}
}

// Register create a new user, which duplicate email is forbidden.
// Create an access token that will be used to decrypt vault
func (s *UserService) Register(params dto.RegisterRequest) (res structs.StdResponse, code int) {
	code = fiber.StatusOK
	_, err := s.userRepo.GetByEmail(params.Email)
	if err != nil && err.Error() != consts.ErrNoData.Error() {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	if err == nil {
		res = structs.StdResponse{Message: "DATA_EXISTS", Data: "Email already registered"}
		code = fiber.StatusBadRequest
		return
	}
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(params.Password), 10)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}

	// Generating all key for accessing vault in latter time
	pub, pvt, err := crypto.GenerateKeyPairEd25519()
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	pvtStr := fmt.Sprintf("%x", pvt)
	newUuid := uuid.GenerateUUID()
	cipher := helper.AbsoluteCharLen(params.Password+fmt.Sprintf("%x", newUuid.Bytes), 16)
	accessToken, err := crypto.EncryptAES(cipher, helper.AbsoluteCharLen(params.Password, 16))
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	backupToken, err := crypto.EncryptAES(cipher, helper.AbsoluteCharLen(pvtStr, 32))
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}

	userId, err := s.userRepo.Create(userRepo.CreateParam{
		Email:       params.Email,
		Password:    string(hashedPass),
		PublicKey:   fmt.Sprintf("%x", pub),
		AccessToken: accessToken,
		BackupToken: backupToken,
	})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.clientRepo.Create(params.FullName, userId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	res = structs.StdResponse{Message: "CREATED", Data: dto.RegisterResponse{
		PrivateKey: pvtStr,
	}}
	return
}

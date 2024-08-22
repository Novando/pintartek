package service

import (
	"context"
	"encoding/json"
	"fmt"
	dtoUser "github.com/Novando/pintartek/internal/passvault-service/app/dto/user"
	clientRepo "github.com/Novando/pintartek/internal/passvault-service/domain/client/repository"
	sessionEntity "github.com/Novando/pintartek/internal/passvault-service/domain/session/entity"
	sessionRepo "github.com/Novando/pintartek/internal/passvault-service/domain/session/repository"
	userRepo "github.com/Novando/pintartek/internal/passvault-service/domain/user/repository"
	"github.com/Novando/pintartek/pkg/common/consts"
	"github.com/Novando/pintartek/pkg/common/structs"
	"github.com/Novando/pintartek/pkg/crypto"
	"github.com/Novando/pintartek/pkg/helper"
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/Novando/pintartek/pkg/redis"
	"github.com/Novando/pintartek/pkg/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type UserConfig func(su *UserService)

type UserService struct {
	log         *logger.Logger
	userRepo    userRepo.User
	clientRepo  clientRepo.Client
	sessionRepo sessionRepo.Session
}

// NewUserService Initialize user service
func NewUserService(config UserConfig, cfgs ...UserConfig) *UserService {
	serv := &UserService{}
	cfgs = append([]UserConfig{config}, cfgs...)
	for _, cfg := range cfgs {
		cfg(serv)
	}
	return serv
}

// WithMock Using Postgres to store data
func WithMock(ur *userRepo.UserMock, sr *sessionRepo.SessionMock, cr *clientRepo.ClientMock) UserConfig {
	return func(su *UserService) {
		su.log = logger.InitZerolog()
		su.userRepo = ur
		su.sessionRepo = sr
		su.clientRepo = cr
	}
}

// WithUserPostgres Using Postgres to store data
func WithUserPostgres(c context.Context, q *pgx.Queries, db *pgxpool.Pool, l *logger.Logger) UserConfig {
	return func(su *UserService) {
		su.log = l
		su.userRepo = userRepo.NewPostgresUserRepository(c, q, db)
		su.clientRepo = clientRepo.NewPostgresClientRepository(c, q, db)
		su.sessionRepo = sessionRepo.NewPostgresSessionRepository(c, q, db)
	}
}

// WithUserRedis Using redis to store session data
func WithUserRedis(r *redis.Redis) UserConfig {
	return func(su *UserService) {
		su.sessionRepo = sessionRepo.NewRedisSessionRepository(r)
	}
}

// Register create a new user, which duplicate email is forbidden.
// Create an access token that will be used to decrypt vault
func (s *UserService) Register(params dtoUser.RegisterRequest) (res structs.StdResponse, code int) {
	_, err := s.userRepo.GetByEmail(params.Email)
	if err != nil && err.Error() != consts.ErrNoData.Error() {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
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
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	newUserUuid := uuid.GenerateUUID()

	// Generating all key for accessing vault in latter time
	pub, pvt, err := crypto.GenerateKeyPairEd25519()
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	pvtStr := fmt.Sprintf("%x", pvt)
	newUuid := uuid.GenerateUUID()
	cipher := helper.AbsoluteCharLen(params.Password+fmt.Sprintf("%x", newUuid.Bytes), 16)
	sessionData, err := json.Marshal(sessionEntity.Session{
		UserID:    newUserUuid,
		SecretKey: cipher,
	})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	accessToken, err := crypto.EncryptAES(string(sessionData), helper.AbsoluteCharLen(params.Password, 16))
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	backupToken, err := crypto.EncryptAES(string(sessionData), helper.AbsoluteCharLen(pvtStr, 32))
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}

	userId, err := s.userRepo.Create(userRepo.CreateParam{
		ID:          newUserUuid,
		Email:       params.Email,
		Password:    string(hashedPass),
		PublicKey:   fmt.Sprintf("%x", pub),
		AccessToken: accessToken,
		BackupToken: backupToken,
	})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.clientRepo.Create(params.FullName, userId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	res = structs.StdResponse{Message: "CREATED", Data: dtoUser.RegisterResponse{
		PrivateKey: pvtStr,
	}}
	code = fiber.StatusOK
	return
}

// Login create a new session, which allow user to access their respective vaults
func (s *UserService) Login(params dtoUser.LoginRequest) (res structs.StdResponse, code int) {
	userData, err := s.userRepo.GetByEmail(params.Email)
	if err != nil {
		msg := "CREDENTIAL_ERROR"
		code = fiber.StatusUnauthorized
		if err.Error() != consts.ErrNoData.Error() {
			s.log.Error(err.Error())
			msg = "REQUEST_ERROR"
			code = fiber.StatusBadRequest
		}
		res = structs.StdResponse{Message: msg, Data: err.Error()}
		return
	}
	if err = bcrypt.CompareHashAndPassword([]byte(userData.Password), []byte(params.Password)); err != nil {
		res = structs.StdResponse{Message: "CREDENTIAL_ERROR", Data: "invalid credential"}
		code = fiber.StatusUnauthorized
		return
	}
	tokenData, err := crypto.DecryptAES(userData.AccessToken, helper.AbsoluteCharLen(params.Password, 16))
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	var sessionData sessionEntity.Session
	if err = json.Unmarshal([]byte(tokenData), &sessionData); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	sessionId, err := s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        uuid.GenerateUUID(),
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	res = structs.StdResponse{
		Message: "SUCCESS",
		Data:    dtoUser.LoginResponse{AccessToken: fmt.Sprintf("%x", sessionId.Bytes)},
	}
	code = fiber.StatusOK
	return
}

// Logout delete an active session of current user
func (s *UserService) Logout(token string) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	err = s.sessionRepo.PermanentDelete(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	res = structs.StdResponse{Message: "SUCCESS", Data: "logged out"}
	code = fiber.StatusOK
	return
}

package service

import (
	"errors"
	"github.com/Novando/pintartek/internal/passvault-service/app/dto/user"
	clientRepo "github.com/Novando/pintartek/internal/passvault-service/domain/client/repository"
	sessionRepo "github.com/Novando/pintartek/internal/passvault-service/domain/session/repository"
	userEntity "github.com/Novando/pintartek/internal/passvault-service/domain/user/entity"
	userRepo "github.com/Novando/pintartek/internal/passvault-service/domain/user/repository"
	"github.com/Novando/pintartek/pkg/common/consts"
	"github.com/Novando/pintartek/pkg/helper"
	"github.com/Novando/pintartek/pkg/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type testUserService struct {
	serv        *UserService
	userMock    *userRepo.UserMock
	clientMock  *clientRepo.ClientMock
	sessionMock *sessionRepo.SessionMock
}

func initTestUserService(t *testing.T) testUserService {
	ms := sessionRepo.NewMockSessionRepository(t)
	mu := userRepo.NewMockUserRepository(t)
	mc := clientRepo.NewMockClientRepository(t)
	return testUserService{
		userMock:    mu,
		clientMock:  mc,
		sessionMock: ms,
		serv:        NewUserService(WithMock(mu, ms, mc)),
	}
}

func TestUserService_Register_Success(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.RegisterRequest{
		Email:           "test@test.com",
		FullName:        "Test User",
		Password:        "passwordpassword",
		ConfirmPassword: "passwordpassword",
	}
	userCreateParam := userRepo.CreateParam{
		ID:          pgtype.UUID{},
		Email:       registerUserParam.Email,
		Password:    "hashedPassword",
		PublicKey:   "publicKey",
		AccessToken: "accessToken",
		BackupToken: "backupToken",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, consts.ErrNoData)
	ts.userMock.Mock.On("Create", userCreateParam).Return(pgtype.UUID{}, nil)
	ts.clientMock.Mock.On("Create", registerUserParam.FullName, pgtype.UUID{}).Return(pgtype.UUID{}, nil)

	res, code := ts.serv.Register(registerUserParam)
	iface, err := helper.InterfaceToMapInterface(res.Data)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "CREATED", res.Message)
	assert.NotEmpty(t, iface["privateKey"])
	assert.Equal(t, http.StatusOK, code)
}

func TestUserService_Register_UserExists(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.RegisterRequest{
		Email:           "test@test.com",
		FullName:        "Test User",
		Password:        "passwordpassword",
		ConfirmPassword: "passwordpassword",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, nil)

	res, code := ts.serv.Register(registerUserParam)
	assert.Equal(t, "DATA_EXISTS", res.Message)
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestUserService_Register_FailCreateUser(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.RegisterRequest{
		Email:           "test@test.com",
		FullName:        "Test User",
		Password:        "passwordpassword",
		ConfirmPassword: "passwordpassword",
	}
	userCreateParam := userRepo.CreateParam{
		ID:          pgtype.UUID{},
		Email:       registerUserParam.Email,
		Password:    "hashedPassword",
		PublicKey:   "publicKey",
		AccessToken: "accessToken",
		BackupToken: "backupToken",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, consts.ErrNoData)
	ts.userMock.Mock.On("Create", userCreateParam).Return(pgtype.UUID{}, errors.New("err"))

	res, code := ts.serv.Register(registerUserParam)
	assert.Equal(t, "PROCESS_ERROR", res.Message)
	assert.Equal(t, http.StatusInternalServerError, code)
}

func TestUserService_Register_FailCreateSession(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.RegisterRequest{
		Email:           "test@test.com",
		FullName:        "Test User",
		Password:        "passwordpassword",
		ConfirmPassword: "passwordpassword",
	}
	userCreateParam := userRepo.CreateParam{
		ID:          pgtype.UUID{},
		Email:       registerUserParam.Email,
		Password:    "hashedPassword",
		PublicKey:   "publicKey",
		AccessToken: "accessToken",
		BackupToken: "backupToken",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, consts.ErrNoData)
	ts.userMock.Mock.On("Create", userCreateParam).Return(pgtype.UUID{}, nil)
	ts.clientMock.Mock.On("Create", registerUserParam.FullName, pgtype.UUID{}).Return(pgtype.UUID{}, errors.New("err"))

	res, code := ts.serv.Register(registerUserParam)
	assert.Equal(t, "PROCESS_ERROR", res.Message)
	assert.Equal(t, http.StatusInternalServerError, code)
}

func TestUserService_Register_FailGetUser(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.RegisterRequest{Email: "test@test.com"}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, errors.New("err"))

	res, code := ts.serv.Register(registerUserParam)
	assert.Equal(t, "REQUEST_ERROR", res.Message)
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestUserService_Login_NoUser(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.LoginRequest{
		Email:    "test@test.com",
		Password: "passwordpassword",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, consts.ErrNoData)

	res, code := ts.serv.Login(registerUserParam)
	assert.Equal(t, "CREDENTIAL_ERROR", res.Message)
	assert.Equal(t, http.StatusUnauthorized, code)
}

func TestUserService_Login_FailGetUser(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.LoginRequest{
		Email:    "test@test.com",
		Password: "passwordpassword",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, errors.New("err"))

	res, code := ts.serv.Login(registerUserParam)
	assert.Equal(t, "REQUEST_ERROR", res.Message)
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestUserService_Login_PasswordMismatch(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.LoginRequest{
		Email:    "test@test.com",
		Password: "password",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, nil)

	res, code := ts.serv.Login(registerUserParam)
	assert.Equal(t, "CREDENTIAL_ERROR", res.Message)
	assert.Equal(t, http.StatusUnauthorized, code)
}

func TestUserService_Login_FailCreateSession(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.LoginRequest{
		Email:    "test@test.com",
		Password: "passwordpassword",
	}
	createSessionParam := sessionRepo.CreateParam{
		ID:        pgtype.UUID{},
		UserID:    pgtype.UUID{},
		SecretKey: "secretKey",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, nil)
	ts.sessionMock.Mock.On("Create", createSessionParam).Return(userEntity.User{}, errors.New("err"))

	res, code := ts.serv.Login(registerUserParam)
	assert.Equal(t, "PROCESS_ERROR", res.Message)
	assert.Equal(t, http.StatusInternalServerError, code)
}

func TestUserService_Login_Success(t *testing.T) {
	ts := initTestUserService(t)
	registerUserParam := user.LoginRequest{
		Email:    "test@test.com",
		Password: "passwordpassword",
	}
	createSessionParam := sessionRepo.CreateParam{
		ID:        pgtype.UUID{},
		UserID:    pgtype.UUID{},
		SecretKey: "secretKey",
	}

	ts.userMock.Mock.On("GetByEmail", registerUserParam.Email).Return(userEntity.User{}, nil)
	ts.sessionMock.Mock.On("Create", createSessionParam).Return(userEntity.User{}, nil)

	res, code := ts.serv.Login(registerUserParam)
	assert.Equal(t, "SUCCESS", res.Message)
	assert.Equal(t, http.StatusOK, code)
}

func TestUserService_Logout_Success(t *testing.T) {
	ts := initTestUserService(t)
	token := "114886bb644e4ef09113952e2bb56b75"
	tokenBytes, _ := uuid.ParseUUID(token)

	ts.sessionMock.Mock.On("PermanentDelete", pgtype.UUID{Bytes: tokenBytes, Valid: true}).Return(userEntity.User{}, nil)

	res, code := ts.serv.Logout(token)
	assert.Equal(t, "SUCCESS", res.Message)
	assert.Equal(t, http.StatusOK, code)
}

func TestUserService_Logout_FailParseToken(t *testing.T) {
	ts := initTestUserService(t)
	token := "114886bb644e"
	res, code := ts.serv.Logout(token)
	assert.Equal(t, "REQUEST_ERROR", res.Message)
	assert.Equal(t, http.StatusBadRequest, code)
}

func TestUserService_Logout_FailDeleteSession(t *testing.T) {
	ts := initTestUserService(t)
	token := "114886bb644e4ef09113952e2bb56b75"
	tokenBytes, _ := uuid.ParseUUID(token)

	ts.sessionMock.Mock.On("PermanentDelete", pgtype.UUID{Bytes: tokenBytes, Valid: true}).Return(userEntity.User{}, errors.New("err"))

	res, code := ts.serv.Logout(token)
	assert.Equal(t, "PROCESS_ERROR", res.Message)
	assert.Equal(t, http.StatusInternalServerError, code)
}

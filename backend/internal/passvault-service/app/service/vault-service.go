package service

import (
	"context"
	"encoding/json"
	"fmt"
	vaultDto "github.com/Novando/pintartek/internal/passvault-service/app/dto/vault"
	sessionRepo "github.com/Novando/pintartek/internal/passvault-service/domain/session/repository"
	vaultGroupRepo "github.com/Novando/pintartek/internal/passvault-service/domain/vault-group/repository"
	vaultRepo "github.com/Novando/pintartek/internal/passvault-service/domain/vault/repository"
	"github.com/Novando/pintartek/pkg/common/consts"
	"github.com/Novando/pintartek/pkg/common/structs"
	"github.com/Novando/pintartek/pkg/crypto"
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/Novando/pintartek/pkg/postgresql/pgx"
	"github.com/Novando/pintartek/pkg/redis"
	"github.com/Novando/pintartek/pkg/uuid"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VaultConfig func(su *VaultService)

type VaultService struct {
	log            *logger.Logger
	vaultRepo      vaultRepo.Vault
	sessionRepo    sessionRepo.Session
	vaultGroupRepo vaultGroupRepo.VaultGroup
}

// NewVaultService Initialize user service
func NewVaultService(config VaultConfig, cfgs ...VaultConfig) *VaultService {
	serv := &VaultService{}
	cfgs = append([]VaultConfig{config}, cfgs...)
	for _, cfg := range cfgs {
		cfg(serv)
	}
	return serv
}

// WithVaultPostgres Using Postgres to store data
func WithVaultPostgres(c context.Context, q *pgx.Queries, db *pgxpool.Pool, l *logger.Logger) VaultConfig {
	return func(sv *VaultService) {
		sv.log = l
		sv.vaultRepo = vaultRepo.NewPostgresVaultRepository(c, q, db)
		sv.sessionRepo = sessionRepo.NewPostgresSessionRepository(c, q, db)
		sv.vaultGroupRepo = vaultGroupRepo.NewPostgresVaultGroupRepository(c, q, db)
	}
}

// WithVaultRedis Using redis to store session data
func WithVaultRedis(r *redis.Redis) VaultConfig {
	return func(sv *VaultService) {
		sv.sessionRepo = sessionRepo.NewRedisSessionRepository(r)
	}
}

// Create build a new vault that contain secret credentials
func (s *VaultService) Create(sessionToken string, param vaultDto.VaultRequest) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(sessionToken)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	sessionData, err := s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}

	mapRes, credential, err := s.processJson(param.Credential, sessionData.SecretKey, fmt.Sprintf("%x", uuid.GenerateUUID().Bytes))
	if err != nil {
		s.log.Error(err.Error())
		msg := "PROCESS_ERROR"
		code = fiber.StatusInternalServerError
		if err.Error() == consts.ErrCrypto.Error() {
			msg = "ACCESS_DENIED"
			code = fiber.StatusUnauthorized
		}
		res = structs.StdResponse{Message: msg, Data: err.Error()}
		return
	}

	vaultId, err := s.vaultRepo.Create(vaultRepo.UpsertParam{
		Name:       param.Name,
		Credential: credential,
	})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
	}
	if err = s.vaultGroupRepo.Create(vaultGroupRepo.CreateParam{VaultID: vaultId, UserID: sessionData.UserID}); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        pgtype.UUID{Bytes: uuid.GenerateUUID().Bytes, Valid: true},
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	res = structs.StdResponse{Message: "CREATED", Data: mapRes}
	code = fiber.StatusOK
	return
}

// GetAll return all vault owned by a user
func (s *VaultService) GetAll(token string) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	sessionData, err := s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}
	vaultData, err := s.vaultGroupRepo.GetAllVaultByUserID(sessionData.UserID, structs.StdPagination{Page: 0, Size: 1000})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	dto := []vaultDto.VaultResponse{}
	for _, item := range vaultData {
		dto = append(dto, vaultDto.VaultResponse{
			ID:        fmt.Sprintf("%x", item.ID.Bytes),
			Name:      item.Name,
			CreatedAt: item.CreatedAt.Time,
			UpdatedAt: item.UpdatedAt.Time,
		})
	}
	_, err = s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        pgtype.UUID{Bytes: uuid.GenerateUUID().Bytes, Valid: true},
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	res = structs.StdResponse{Message: "FETCHED", Data: dto, Count: int64(len(dto))}
	code = fiber.StatusOK
	return
}

// GetOne decrypt the credential of a vault
func (s *VaultService) GetOne(token, vaultId string) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	sessionData, err := s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}
	vaultBytes, err := uuid.ParseUUID(vaultId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	vaultData, err := s.vaultRepo.GetByID(pgtype.UUID{Bytes: vaultBytes, Valid: true})
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	credentials, err := crypto.DecryptAES(vaultData.Credential, sessionData.SecretKey)
	if err != nil {
		res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
		code = fiber.StatusUnauthorized
		return
	}
	var mapCredentials map[string]interface{}
	err = json.Unmarshal([]byte(credentials), &mapCredentials)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        pgtype.UUID{Bytes: uuid.GenerateUUID().Bytes, Valid: true},
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	res = structs.StdResponse{Message: "FETCHED", Data: mapCredentials}
	code = fiber.StatusOK
	return
}

// UpdateVaultName update the name of a vault
func (s *VaultService) UpdateVaultName(
	token,
	vaultId string,
	param vaultDto.VaultEditRequest,
) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	sessionData, err := s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}
	vaultBytes, err := uuid.ParseUUID(vaultId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	if err = s.vaultRepo.UpdateName(pgtype.UUID{Bytes: vaultBytes, Valid: true}, param.Name); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        pgtype.UUID{Bytes: uuid.GenerateUUID().Bytes, Valid: true},
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	res = structs.StdResponse{Message: "UPDATED"}
	code = fiber.StatusOK
	return
}

// UpdateCredential update the credential of a vault
func (s *VaultService) UpdateCredential(
	token string,
	vaultId string,
	credentialId string,
	param vaultDto.Credential,
) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	sessionData, err := s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}
	vaultBytes, err := uuid.ParseUUID(vaultId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	vaultUuid := pgtype.UUID{Bytes: vaultBytes, Valid: true}
	vaultData, err := s.vaultRepo.GetByID(vaultUuid)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	credentials, err := crypto.DecryptAES(vaultData.Credential, sessionData.SecretKey)
	if err != nil {
		res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
		code = fiber.StatusUnauthorized
		return
	}
	mapRes, credential, err := s.processJson(param, sessionData.SecretKey, credentialId, credentials)
	if err != nil {
		s.log.Error(err.Error())
		msg := "PROCESS_ERROR"
		code = fiber.StatusInternalServerError
		if err.Error() == consts.ErrCrypto.Error() {
			msg = "ACCESS_DENIED"
			code = fiber.StatusUnauthorized
		}
		res = structs.StdResponse{Message: msg, Data: err.Error()}
		return
	}
	if err = s.vaultRepo.UpdateCredential(vaultUuid, credential); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        pgtype.UUID{Bytes: uuid.GenerateUUID().Bytes, Valid: true},
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	res = structs.StdResponse{Message: "UPDATED", Data: mapRes}
	code = fiber.StatusOK
	return
}

// CreateCredential create the credential to a vault
func (s *VaultService) CreateCredential(
	token string,
	vaultId string,
	param vaultDto.Credential,
) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	sessionData, err := s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}
	vaultBytes, err := uuid.ParseUUID(vaultId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	vaultUuid := pgtype.UUID{Bytes: vaultBytes, Valid: true}
	vaultData, err := s.vaultRepo.GetByID(vaultUuid)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	credentials, err := crypto.DecryptAES(vaultData.Credential, sessionData.SecretKey)
	if err != nil {
		res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
		code = fiber.StatusUnauthorized
		return
	}
	mapRes, credential, err := s.processJson(
		param,
		sessionData.SecretKey,
		fmt.Sprintf("%x", uuid.GenerateUUID().Bytes),
		credentials,
	)
	if err != nil {
		s.log.Error(err.Error())
		msg := "PROCESS_ERROR"
		code = fiber.StatusInternalServerError
		if err.Error() == consts.ErrCrypto.Error() {
			msg = "ACCESS_DENIED"
			code = fiber.StatusUnauthorized
		}
		res = structs.StdResponse{Message: msg, Data: err.Error()}
		return
	}
	if err = s.vaultRepo.UpdateCredential(vaultUuid, credential); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        pgtype.UUID{Bytes: uuid.GenerateUUID().Bytes, Valid: true},
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	res = structs.StdResponse{Message: "CREATED", Data: mapRes}
	code = fiber.StatusOK
	return
}

// DeleteCredential delete a credential from a vault
func (s *VaultService) DeleteCredential(
	token string,
	vaultId string,
	credentialId string,
) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	sessionData, err := s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}
	vaultBytes, err := uuid.ParseUUID(vaultId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	vaultUuid := pgtype.UUID{Bytes: vaultBytes, Valid: true}
	vaultData, err := s.vaultRepo.GetByID(vaultUuid)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	credentials, err := crypto.DecryptAES(vaultData.Credential, sessionData.SecretKey)
	if err != nil {
		res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
		code = fiber.StatusUnauthorized
		return
	}
	mapRes, credential, err := s.processJson(nil, sessionData.SecretKey, credentialId, credentials)
	if err != nil {
		s.log.Error(err.Error())
		msg := "PROCESS_ERROR"
		code = fiber.StatusInternalServerError
		if err.Error() == consts.ErrCrypto.Error() {
			msg = "ACCESS_DENIED"
			code = fiber.StatusUnauthorized
		}
		res = structs.StdResponse{Message: msg, Data: err.Error()}
		return
	}
	if err = s.vaultRepo.UpdateCredential(vaultUuid, credential); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	_, err = s.sessionRepo.Create(sessionRepo.CreateParam{
		ID:        pgtype.UUID{Bytes: uuid.GenerateUUID().Bytes, Valid: true},
		UserID:    sessionData.UserID,
		SecretKey: sessionData.SecretKey,
	})
	res = structs.StdResponse{Message: "DELETED", Data: mapRes}
	code = fiber.StatusOK
	return
}

// Delete delete a vault permanently
func (s *VaultService) Delete(
	token string,
	vaultId string,
) (res structs.StdResponse, code int) {
	tokenBytes, err := uuid.ParseUUID(token)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	_, err = s.sessionRepo.GetByID(pgtype.UUID{Bytes: tokenBytes, Valid: true})
	if err != nil {
		if err.Error() == consts.ErrNoData.Error() {
			res = structs.StdResponse{Message: "ACCESS_DENIED", Data: err.Error()}
			code = fiber.StatusUnauthorized
		} else {
			s.log.Error(err.Error())
			res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
			code = fiber.StatusInternalServerError
		}
		return
	}
	vaultBytes, err := uuid.ParseUUID(vaultId)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "REQUEST_ERROR", Data: err.Error()}
		code = fiber.StatusBadRequest
		return
	}
	if err = s.vaultRepo.PermanentDelete(pgtype.UUID{Bytes: vaultBytes, Valid: true}); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
		return
	}
	res = structs.StdResponse{Message: "DELETED", Data: fmt.Sprintf("vaultId %v has been deleted", vaultId)}
	code = fiber.StatusOK
	return
}

// processJson restructure the JSON and append/update new credential value,
// and encrypt the credential. pass nil to `credential` to delete a field
func (s *VaultService) processJson(
	credential interface{},
	cipher,
	credentialId string,
	existingCredential ...string,
) (mapRes map[string]interface{}, res string, err error) {
	// restructure using named JSON
	if mapRes == nil {
		mapRes = make(map[string]interface{})
	}
	if len(existingCredential) > 0 {
		err = json.Unmarshal([]byte(existingCredential[0]), &mapRes)
		if err != nil {
			return
		}
	}
	paramJson, err := json.Marshal(credential)
	if err != nil {
		return
	}
	var mapJson map[string]interface{}
	if err = json.Unmarshal(paramJson, &mapJson); err != nil {
		return
	}
	if credential != nil {
		mapRes[credentialId] = mapJson
	} else {
		delete(mapRes, credentialId)
	}
	paramJson, err = json.Marshal(mapRes)
	if err != nil {
		return
	}

	// encrypt the credentials
	res, err = crypto.EncryptAES(string(paramJson), cipher)
	if err != nil {
		s.log.Error(err.Error())
		err = consts.ErrCrypto
	}
	return
}

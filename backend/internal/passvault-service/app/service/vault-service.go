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
	code = 200
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

	// restructure using named JSON
	paramJson, err := json.Marshal(param.Credential)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
	}
	var mapJson, mapRes map[string]interface{}
	if err = json.Unmarshal(paramJson, &mapJson); err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
	}
	mapRes[fmt.Sprintf("%x", uuid.GenerateUUID().Bytes)] = mapJson
	paramJson, err = json.Marshal(mapRes)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
	}

	// encrypt the credentials
	credential, err := crypto.EncryptAES(string(paramJson), sessionData.SecretKey)
	if err != nil {
		s.log.Error(err.Error())
		res = structs.StdResponse{Message: "PROCESS_ERROR", Data: err.Error()}
		code = fiber.StatusInternalServerError
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
	}
	res = structs.StdResponse{Message: "CREATED", Data: mapRes}
	return
}

func (s *VaultService) GetAll(token string) (res structs.StdResponse, code int) {
	code = 200
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
	res = structs.StdResponse{Message: "FETCHED", Data: dto, Count: int64(len(dto))}
	return
}

//func (s *VaultService) mergeJson(json1, json2 string) (res string, err error) {
//	var map1, map2, mapRes map[string]interface{}
//
//	// Unmarshal the first JSON into map1.
//	err := json.Unmarshal([]byte(json1), &map1)
//	if err != nil {
//		fmt.Println("Error unmarshalling json1:", err)
//		return
//	}
//
//	// Unmarshal the second JSON into map2.
//	err = json.Unmarshal([]byte(json2), &map2)
//	if err != nil {
//		fmt.Println("Error unmarshalling json2:", err)
//		return
//	}
//
//	// Extract the "id" key from map2 to create the nested structure.
//	id := map2["id"].(string)
//	delete(map2, "id")
//
//	// Merge map2 into map1 under the key of the extracted "id".
//	map1[id] = map2
//
//	// Marshal the resulting map back to a JSON string.
//	mergedJSON, err := json.Marshal(map1)
//	if err != nil {
//		fmt.Println("Error marshalling result:", err)
//		return
//	}
//
//	// Print the final merged JSON string.
//	fmt.Println("Merged JSON:", string(mergedJSON))
//}

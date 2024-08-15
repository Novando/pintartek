package repository

import (
	"encoding/json"
	"fmt"
	"github.com/Novando/pintartek/internal/passvault-service/domain/session/entity"
	"github.com/Novando/pintartek/pkg/redis"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type RedisSession struct {
	rds *redis.Redis
}

func NewRedisSessionRepository(r *redis.Redis) *RedisSession {
	return &RedisSession{rds: r}
}

func (r *RedisSession) Create(arg CreateParam) (id pgtype.UUID, err error) {
	sessionData := entity.Session{UserID: arg.UserID, SecretKey: arg.SecretKey}
	val, err := json.Marshal(sessionData)
	if err != nil {
		return
	}
	if err = r.rds.Set(fmt.Sprintf("%x", arg.ID.Bytes), string(val), time.Minute*30); err != nil {
		return
	}
	id = arg.ID
	return
}

func (r *RedisSession) GetByID(id pgtype.UUID) (session entity.Session, err error) {
	val, err := r.rds.Get(fmt.Sprintf("%x", id.Bytes))
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(val), &session)
	return
}

func (r *RedisSession) PermanentDelete(id pgtype.UUID) error {
	return r.rds.Delete(fmt.Sprintf("%x", id.Bytes))
}

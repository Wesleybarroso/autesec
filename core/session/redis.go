package session

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

func NewRedisStore(addr, password string, db int, ttl time.Duration) *RedisStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &RedisStore{
		client: rdb,
		ttl:    ttl,
	}
}

func (r *RedisStore) Save(ctx context.Context, s *Session) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	return r.client.Set(ctx, s.SessionID, data, r.ttl).Err()
}

func (r *RedisStore) Get(ctx context.Context, sessionID string) (*Session, error) {
	val, err := r.client.Get(ctx, sessionID).Result()
	if err != nil {
		return nil, err
	}

	var s Session
	if err := json.Unmarshal([]byte(val), &s); err != nil {
		return nil, err
	}

	return &s, nil
}

func (r *RedisStore) Delete(ctx context.Context, sessionID string) error {
	return r.client.Del(ctx, sessionID).Err()
}

func (r *RedisStore) UpdateLastSeen ( ctx context.Context, sessionID string) error {
	
	//1. Buscar sessão atual
	s, err := r.Get(ctx, sessionID)
	if err != nil {
		return err 
	}
	//2. Atualiza o campo LastSeen
	s.LastSeen = time.Now()

	//3. Salva de volta no Redis 
	return r.Save(ctx, s)
}

package store

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/18211167516/gosession"
	"github.com/go-redis/redis/v8"
)

var Ctx context.Context

type RedisStore struct {
	reidsClient *redis.Client
}

func NewRedis(addr, password string, db, PoolSize int) *RedisStore {

	ClientRedis := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
		PoolSize: PoolSize, // 连接池大小
		//MinIdleConns: 5,
	})

	Ctx = ClientRedis.Context()

	err := ClientRedis.Ping(Ctx).Err()
	if err != nil {
		log.Panic("Cache store redis：", err)
	}

	return &RedisStore{
		reidsClient: ClientRedis,
	}
}

func RegisterRedis(addr, password string, db, PoolSize int) {
	gosession.Register("redis", NewRedis(addr, password, db, PoolSize))
}

func (s *RedisStore) FmtKey(Sid, key string) string {
	return fmt.Sprintf("%s_%s", Sid, key)
}

func (s *RedisStore) Get(Sid, key string) (interface{}, error) {
	return s.reidsClient.Get(Ctx, s.FmtKey(Sid, key)).Result()
}

func (s *RedisStore) Set(Sid, key string, value interface{}, d int) error {
	return s.reidsClient.Set(Ctx, s.FmtKey(Sid, key), value, time.Duration(d)*time.Second).Err()
}

func (s *RedisStore) Remove(Sid, key string) error {
	return s.reidsClient.Del(Ctx, s.FmtKey(Sid, key)).Err()
}

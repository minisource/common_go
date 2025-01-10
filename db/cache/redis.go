package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client

type RedisConfig struct {
	Host               string
	Port               string
	Password           string
	Db                 string
	DialTimeout        time.Duration
	ReadTimeout        time.Duration
	WriteTimeout       time.Duration
	IdleCheckFrequency time.Duration
	PoolSize           int
	PoolTimeout        time.Duration
}

func InitRedis(cfg *RedisConfig) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:               fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password:           cfg.Password,
		DB:                 0,
		DialTimeout:        cfg.DialTimeout * time.Second,
		ReadTimeout:        cfg.ReadTimeout * time.Second,
		WriteTimeout:       cfg.WriteTimeout * time.Second,
		PoolSize:           cfg.PoolSize,
		PoolTimeout:        cfg.PoolTimeout,
		IdleTimeout:        500 * time.Millisecond,
		IdleCheckFrequency: cfg.IdleCheckFrequency * time.Millisecond,
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() {
	redisClient.Close()
}

func Set[T any](c *redis.Client, key string, value T, duration time.Duration) error {
	v, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return c.Set(key, v, duration).Err()
}

func Get[T any](c *redis.Client, key string) (T, error) {
	var dest T = *new(T)
	v, err := c.Get(key).Result()
	if err != nil {
		return dest, err
	}
	err = json.Unmarshal([]byte(v), &dest)
	if err != nil {
		return dest, err
	}
	return dest, nil
}

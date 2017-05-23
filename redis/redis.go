package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/mitchellh/mapstructure"
)

type RedisConfig struct {
	Network        string
	Address        string
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	Pool           redis.Pool
}

func NewRedis(cfg *RedisConfig) (redis.Pool, error) {

	cfg.Pool.Dial = func() (redis.Conn, error) {
		return redis.DialTimeout(cfg.Network, cfg.Address, cfg.ConnectTimeout,
			cfg.ReadTimeout, cfg.WriteTimeout)
	}
	return cfg.Pool, nil
}

func Creator(cfg interface{}) (interface{}, error) {
	var redisConfig RedisConfig
	err := mapstructure.WeakDecode(cfg, &redisConfig)
	if err != nil {
		return nil, err
	}
	redisConfig.ConnectTimeout = redisConfig.ConnectTimeout * time.Millisecond
	redisConfig.ReadTimeout = redisConfig.ReadTimeout * time.Millisecond
	redisConfig.WriteTimeout = redisConfig.WriteTimeout * time.Millisecond
	return NewRedis(&redisConfig)
}

package go_redis

import (
	"github.com/go-redis/redis"
)

type Redis struct {
	clients map[string]*redis.Client
}

func (this *Redis) Get(name string) *redis.Client {
	return this.clients[name]
}

func New(cfg *Config) (*Redis, error) {
	reply := &Redis{
		clients: make(map[string]*redis.Client),
	}
	for name, config := range *cfg {
		c := redis.NewClient(config)
		_, err := c.Ping().Result()
		if err != nil {
			return nil, err
		}
		reply.clients[name] = c
	}
	return reply, nil
}

type Config map[string]*redis.Options

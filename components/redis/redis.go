package redis

import (
	"errors"
	"fmt"
	"log"
	"time"

	redigo "github.com/garyburd/redigo/redis"
)

var (
	redisPool = map[string]*redigo.Pool{}

	errNotFoundRedis = errors.New("Redis Instance Not Found")

	dialFunc = func(network, address string, dialOptions []redigo.DialOption) func() (redigo.Conn, error) {
		return func() (redigo.Conn, error) {
			conn, err := redigo.Dial(network, address, dialOptions...)
			if err != nil {
				log.Panic(err)
			}
			return conn, err
		}
	}

	testOnBorrowFunc = func(c redigo.Conn, t time.Time) error {
		_, err := c.Do("ping")
		if err != nil {
			return err
		}
		return nil
	}

	testFunc = func(p *redigo.Pool) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = errors.New(fmt.Sprintf("%+v", r))
			}

		}()
		c := p.Get()
		defer c.Close()
		return c.Err()
	}
)

type Redis struct {
	redisPool map[string]*redigo.Pool
}

func (this *Redis) Get(name string) redigo.Conn {
	p, ok := this.redisPool[name]
	if !ok {
		panic(errNotFoundRedis)
	}
	return p.Get()
}

func New(cfg *Config) (*Redis, error) {
	reply := &Redis{
		redisPool: make(map[string]*redigo.Pool),
	}
	for name, config := range *cfg {

		r := &redigo.Pool{}

		dialOptions := []redigo.DialOption{}

		if config.DialConnectionTimeout > 0 {
			dialOptions = append(dialOptions, redigo.DialConnectTimeout(time.Second*time.Duration(config.DialConnectionTimeout)))
		}

		if config.DialReadTimeout > 0 {
			dialOptions = append(dialOptions, redigo.DialReadTimeout(time.Second*time.Duration(config.DialReadTimeout)))
		}

		if config.DialWriteTimeout > 0 {
			dialOptions = append(dialOptions, redigo.DialWriteTimeout(time.Second*time.Duration(config.DialWriteTimeout)))
		}

		if config.DialPassword != "" {
			dialOptions = append(dialOptions, redigo.DialPassword(config.DialPassword))
		}

		dialOptions = append(dialOptions, redigo.DialDatabase(config.DB))

		r.Dial = dialFunc(config.Network, config.Address, dialOptions)

		if config.MaxIdle > 0 {
			r.MaxIdle = config.MaxIdle
		}
		if config.MaxActive > 0 {
			r.MaxActive = config.MaxActive
		}
		if config.TestOnBorrow {
			r.TestOnBorrow = testOnBorrowFunc
		}

		if config.IdleTimeout > 0 {
			r.IdleTimeout = time.Second * time.Duration(config.IdleTimeout)
		}
		r.Wait = config.Wait

		if err := testFunc(r); err != nil {
			return nil, err
		}
		reply.redisPool[name] = r

	}
	return reply, nil
}

type Config map[string]struct {
	Network               string
	Address               string
	DialConnectionTimeout int
	DialReadTimeout       int
	DialWriteTimeout      int
	DialPassword          string
	DB                    int
	MaxIdle               int
	MaxActive             int
	TestOnBorrow          bool
	IdleTimeout           int
	Wait                  bool
}

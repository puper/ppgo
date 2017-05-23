package app

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	"github.com/juju/errors"
	"github.com/puper/p2pwatch/container"
	"github.com/puper/p2pwatch/dbman"
)

var (
	defaultContainer = container.NewContainer()
)

func Register(name string, creator container.Creator) {
	defaultContainer.Register(name, creator)
}

func ConfigureAll(cfg map[string]interface{}) error {
	return defaultContainer.ConfigureAll(cfg)
}

func Get(name string) (interface{}, error) {
	return defaultContainer.Get(name)
}

func MustGet(name string) interface{} {
	instance, _ := defaultContainer.Get(name)
	return instance
}

func MustGetRedis(name string) redis.Conn {
	pool, _ := defaultContainer.Get(name)
	pool2, _ := pool.(redis.Pool)
	return pool2.Get()
}

func GetDBMan() (*dbman.DBMan, error) {
	instance, err := Get("dbman")
	if err != nil {
		return nil, err
	}
	if dm, ok := instance.(*dbman.DBMan); ok {
		return dm, nil
	}
	return nil, errors.NotValidf("can not trans interface to dbman")
}

type Model interface {
	ConnName() string
}

func ChooseDB(model Model, write bool) *gorm.DB {
	dbm, err := GetDBMan()
	if err != nil {
		return nil
	}
	db, _ := dbm.Get(model.ConnName(), write)
	return db
}

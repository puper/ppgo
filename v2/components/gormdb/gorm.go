package gormdb

import (
	"time"

	"github.com/puper/ppgo/helpers"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DB struct {
		wrappers map[string]*Wrapper
	}
	Wrapper struct {
		master *gorm.DB
		slave  []*gorm.DB
	}
	Config map[string]struct {
		Master          string
		Slave           []string
		ConnMaxLifeTime int // in second
		MaxIdleConns    int
		MaxOpenConns    int
		Callback        func(*gorm.DB) *gorm.DB
	}
	Model interface {
		ConnectionName() string
	}
)

func New(cfg *Config) (*DB, error) {
	man := &DB{
		wrappers: make(map[string]*Wrapper),
	}
	var err error
	for name, config := range *cfg {
		w := new(Wrapper)
		w.master, err = gorm.Open(mysql.Open(config.Master), &gorm.Config{})
		if err != nil {
			return nil, err
		}
		rawDb, err := w.master.DB()
		if err != nil {
			return nil, err
		}
		rawDb.SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
		rawDb.SetMaxIdleConns(config.MaxIdleConns)
		rawDb.SetMaxOpenConns(config.MaxOpenConns)
		if config.Callback != nil {
			w.master = config.Callback(w.master)
		}
		for _, s := range config.Slave {
			slave, err := gorm.Open(mysql.Open(s), &gorm.Config{})
			if err != nil {
				return nil, err
			}
			rawDb, err := slave.DB()
			if err != nil {
				return nil, err
			}
			rawDb.SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
			rawDb.SetMaxIdleConns(config.MaxIdleConns)
			rawDb.SetMaxOpenConns(config.MaxOpenConns)
			if config.Callback != nil {
				slave = config.Callback(slave)
			}
			w.slave = append(w.slave, slave)
		}
		man.wrappers[name] = w
	}
	return man, nil
}

func (this *Wrapper) Write() *gorm.DB {
	return this.master
}

func (this *Wrapper) Read() *gorm.DB {
	if len(this.slave) == 0 {
		return this.master
	}
	return this.slave[helpers.GlobalRand().Intn(len(this.slave))]
}

func (this *DB) Write(name string) *gorm.DB {
	return this.wrappers[name].Write()
}

func (this *DB) Read(name string) *gorm.DB {
	return this.wrappers[name].Read()
}

func (this *DB) WriteModel(m Model) *gorm.DB {
	return this.Write(m.ConnectionName()).Model(m)
}

func (this *DB) ReadModel(m Model) *gorm.DB {
	return this.Read(m.ConnectionName()).Model(m)
}

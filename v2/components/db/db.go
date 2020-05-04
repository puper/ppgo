package db

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/puper/ppgo/helpers"
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
		Driver          string
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
		w.master, err = gorm.Open(config.Driver, config.Master)
		w.master.DB().SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
		w.master.DB().SetMaxIdleConns(config.MaxIdleConns)
		w.master.DB().SetMaxOpenConns(config.MaxOpenConns)
		if config.Callback != nil {
			w.master = config.Callback(w.master)
		}
		if err != nil {
			return nil, err
		}
		for _, s := range config.Slave {
			slave, err := gorm.Open(config.Driver, s)
			if err != nil {
				return nil, err
			}
			slave.DB().SetConnMaxLifetime(time.Duration(config.ConnMaxLifeTime) * time.Second)
			slave.DB().SetMaxIdleConns(config.MaxIdleConns)
			slave.DB().SetMaxOpenConns(config.MaxOpenConns)
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

package dbman

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/mitchellh/mapstructure"
	"github.com/puper/errors"
)

type DBMan struct {
	Config *DBManConfig
	RWDBs  map[string]*RWDB
}

type DBManConfig struct {
	RWDBs map[string]*RWDBConfig
}

func NewDBMan(cfg *DBManConfig) (*DBMan, error) {
	dbman := &DBMan{
		Config: cfg,
		RWDBs:  make(map[string]*RWDB),
	}
	for k, v := range cfg.RWDBs {
		rwdb, err := NewRWDB(v)
		if err != nil {
			return nil, err
		}
		dbman.RWDBs[k] = rwdb
	}
	return dbman, nil

}

func (this *DBMan) Get(name string, write bool) (*gorm.DB, error) {
	if rwdb, ok := this.RWDBs[name]; ok {
		return rwdb.Get(write), nil
	}
	return nil, errors.NotFoundf("rwdb %s not found", name)
}

func Creator(cfg interface{}) (interface{}, error) {
	var c DBManConfig
	err := mapstructure.WeakDecode(cfg, &c)
	if err != nil {
		return nil, errors.Annotatef(err, "decode db config error: %v", cfg)
	}
	return NewDBMan(&c)
}

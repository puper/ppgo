package dbman

import (
	"math/rand"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/puper/errors"
)

type RWDB struct {
	Config *RWDBConfig
	Write  *gorm.DB
	Reads  []*gorm.DB
}

type RWDBConfig struct {
	Driver string
	Dsn    string
	Reads  []string
}

func NewRWDB(cfg *RWDBConfig) (*RWDB, error) {
	var err error
	rwdb := &RWDB{
		Config: cfg,
	}
	rwdb.Write, err = gorm.Open(cfg.Driver, cfg.Dsn)
	if err != nil {
		return nil, errors.Annotatef(err, "open write db %v,%v failed", cfg.Driver, cfg.Dsn)
	}
	for dsn := range cfg.Reads {
		rdb, err := gorm.Open(cfg.Driver, dsn)
		if err != nil {
			return nil, errors.Annotatef(err, "open read db %v,%v failed", cfg.Driver, dsn)
		}
		rwdb.Reads = append(rwdb.Reads, rdb)
	}
	return rwdb, nil
}

func (this *RWDB) Get(write bool) *gorm.DB {
	if write {
		return this.Write
	}
	l := len(this.Reads)
	if l == 0 {
		return this.Write
	}
	rand.Seed(time.Now().UnixNano())
	return this.Reads[rand.Intn(l)]

}

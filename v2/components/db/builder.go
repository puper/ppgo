package db

import (
	"github.com/jinzhu/gorm"
	"github.com/puper/ppgo/helpers"
	"github.com/puper/ppgo/v2/engine"
)

func Builder(configKey string, callback func(*gorm.DB) *gorm.DB) engine.Builder {
	return func(e *engine.Engine) (interface{}, error) {
		cfg := e.GetConfig().Get(configKey)
		c := &Config{}
		if err := helpers.StructDecode(cfg, c, "json"); err != nil {
			return nil, err
		}
		for k, v := range *c {
			v.Callback = callback
			(*c)[k] = v
		}
		return New(c)
	}
}

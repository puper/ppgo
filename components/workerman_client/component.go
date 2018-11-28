package workerman_client

import (
	"github.com/puper/ppgo/engine"
	"github.com/puper/ppgo/helpers"
)

type Component struct {
	engine.BaseComponent
}

func (this *Component) Create(cfg interface{}) (interface{}, error) {
	c := &Config{}
	if err := helpers.StructDecode(cfg, c, "json"); err != nil {
		return nil, err
	}
	return New(c)
}

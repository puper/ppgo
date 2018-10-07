package gin

import (
	"context"
	"time"

	"github.com/facebookgo/grace/gracehttp"
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

func (this *Component) Init(_, intstance interface{}) error {
	return nil
}

func (this *Component) Start(_, instance interface{}) error {
	gin := instance.(*Gin)
	return gracehttp.Serve(gin.svr)
}

func (this *Component) Stop(_, instance interface{}) error {
	gin := instance.(*Gin)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return gin.svr.Shutdown(ctx)
}

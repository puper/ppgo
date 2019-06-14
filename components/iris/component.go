package iris

import (
	"context"

	"github.com/kataras/iris"

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
	server := instance.(*Iris)
	return server.application.Run(
		iris.Addr(server.config.Addr),
		iris.WithoutServerError(
			iris.ErrServerClosed,
		),
		iris.WithoutPathCorrection,
	)
}

func (this *Component) Stop(_, instance interface{}) error {
	server := instance.(*Iris)
	return server.application.Shutdown(context.Background())
}

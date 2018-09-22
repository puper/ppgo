package grpc_server

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

func (this *Component) Start(_, instance interface{}) error {
	gs := instance.(*GRPCServer)
	return gs.svr.Serve(gs.lis)
}

func (this *Component) Stop(_, instance interface{}) error {
	gs := instance.(*GRPCServer)
	gs.svr.GracefulStop()
	return nil
}

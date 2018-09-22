package grpc_server

import (
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	lis net.Listener
	svr *grpc.Server
}

func (this *GRPCServer) GetEngine() *grpc.Server {
	return this.svr
}

func New(cfg *Config) (*GRPCServer, error) {
	var (
		err error
	)
	gs := new(GRPCServer)
	gs.lis, err = net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	gs.svr = grpc.NewServer()
	return gs, nil
}

type Config struct {
	Addr string
}

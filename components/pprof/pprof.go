package pprof

import (
	"net"
	"net/http"
	_ "net/http/pprof"
)

type Config struct {
	Addr string
}

type PProf struct {
	config *Config
	lis    net.Listener
}

func New(cfg *Config) (*PProf, error) {
	var (
		err error
	)
	lis, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	go http.Serve(lis, nil)
	return &PProf{
		config: cfg,
		lis:    lis,
	}, nil
}

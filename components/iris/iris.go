package iris

import (
	"github.com/kataras/iris"
)

type Config struct {
	Addr string
}

type Iris struct {
	application *iris.Application
	config      *Config
}

func (this *Iris) GetApplication() *iris.Application {
	return this.application
}

func New(cfg *Config) (*Iris, error) {
	app := iris.New()
	return &Iris{
		application: app,
		config:      cfg,
	}, nil
}

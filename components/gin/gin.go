package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Addr string
}

type Gin struct {
	engine *gin.Engine
	config *Config
	svr    *http.Server
}

func (this *Gin) GetEngine() *gin.Engine {
	return this.engine
}

func New(cfg *Config) (*Gin, error) {
	e := gin.Default()
	return &Gin{
		engine: e,
		config: cfg,
		svr: &http.Server{
			Addr:    cfg.Addr,
			Handler: e,
		},
	}, nil
}

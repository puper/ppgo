package gin

import (
	"net/http"
	"time"

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
	e := gin.New()
	return &Gin{
		engine: e,
		config: cfg,
		svr: &http.Server{
			Addr:    cfg.Addr,
			Handler: e,
		},
	}, nil
}

type AccessRecorder func(time.Time, time.Time, time.Duration, int, string, string, string)

func AccessRecordMiddleware(out AccessRecorder) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		c.Next()
		end := time.Now()
		latency := end.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		if raw != "" {
			path = path + "?" + raw
		}
		out(start, end, latency, statusCode, clientIP, method, path)
	}
}

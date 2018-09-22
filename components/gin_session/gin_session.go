package gin_session

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Name         string
	CookieSecret string
}

func New(cfg *Config) (gin.HandlerFunc, error) {
	store := cookie.NewStore([]byte(cfg.CookieSecret))
	return sessions.Sessions(cfg.Name, store), nil
}

package app

import (
	"code.int.thoseyears.com/golang/ppgo/components/dbman"
	pgin "code.int.thoseyears.com/golang/ppgo/components/gin"
	"code.int.thoseyears.com/golang/ppgo/components/log"
	"code.int.thoseyears.com/golang/ppgo/components/redis"
	"code.int.thoseyears.com/golang/ppgo/engine"
	"github.com/Sirupsen/logrus"
	redigo "github.com/garyburd/redigo/redis"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var (
	app *engine.Engine
)

func Create(cfg *engine.Config) *engine.Engine {
	app = engine.New(cfg)
	return app
}

func Get() *engine.Engine {
	return app
}

func GetDB() *dbman.DBMan {
	return app.GetInstance("db").(*dbman.DBMan)
}

func GetLog(name string) *logrus.Logger {
	return app.GetInstance("log").(*log.Log).Get(name)
}

func GetServer() *gin.Engine {
	return app.GetInstance("server").(*pgin.Gin).GetEngine()
}

func GetSessionMiddleware() gin.HandlerFunc {
	return app.GetInstance("session").(gin.HandlerFunc)
}

func GetSession(c *gin.Context) sessions.Session {
	return sessions.Default(c)
}

func GetRedis(name string) redigo.Conn {
	return app.GetInstance("redis").(*redis.Redis).Get(name)
}

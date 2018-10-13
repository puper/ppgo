package components

import (
	"github.com/puper/ppgo/components/dbman"
	"github.com/puper/ppgo/components/gin"
	"github.com/puper/ppgo/components/gin_session"
	"github.com/puper/ppgo/components/go_redis"
	"github.com/puper/ppgo/components/grpc_server"
	"github.com/puper/ppgo/components/log"
	"github.com/puper/ppgo/components/pprof"
	"github.com/puper/ppgo/components/redis"
	"github.com/puper/ppgo/engine"
)

func Components() map[string]engine.Component {
	return map[string]engine.Component{
		"db":          &dbman.Component{},
		"redis":       &redis.Component{},
		"go_redis":    &go_redis.Component{},
		"log":         &log.Component{},
		"gin":         &gin.Component{},
		"gin_session": &gin_session.Component{},
		"grpc_server": &grpc_server.Component{},
		"pprof":       &pprof.Component{},
	}
}

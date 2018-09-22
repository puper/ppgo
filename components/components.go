package components

import (
	"code.int.thoseyears.com/golang/ppgo/components/dbman"
	"code.int.thoseyears.com/golang/ppgo/components/gin"
	"code.int.thoseyears.com/golang/ppgo/components/gin_session"
	"code.int.thoseyears.com/golang/ppgo/components/go_redis"
	"code.int.thoseyears.com/golang/ppgo/components/grpc_server"
	"code.int.thoseyears.com/golang/ppgo/components/log"
	"code.int.thoseyears.com/golang/ppgo/components/redis"
	"code.int.thoseyears.com/golang/ppgo/engine"
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
	}
}

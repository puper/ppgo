package bootstrap

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/puper/p2pwatch/app"
	"github.com/puper/p2pwatch/handlers"

	"github.com/go-ozzo/ozzo-routing"
	"github.com/puper/p2pwatch/dbman"
	"github.com/puper/p2pwatch/endless"
	"github.com/puper/p2pwatch/listener"
	"github.com/puper/p2pwatch/redis"
	"github.com/spf13/viper"
)

func Init() error {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.Level(viper.GetInt("log.Config.Level")))
	out, err := os.OpenFile(viper.GetString("log.Config.FileName"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0660)
	if err != nil {
		log.Panicf("log file err %v", err)
	}
	log.SetOutput(out)
	app.Register("dbman", dbman.Creator)
	app.Register("redis", redis.Creator)
	return app.ConfigureAll(viper.GetStringMap("components"))
}

func Run() error {
	listener.SetConfig(viper.GetString("server.Addr"))
	r := routing.New()
	r.Any("/puper", handlers.Index)
	http.Handle("/", r)
	endless.ListenAndServe(viper.GetString("server.Addr"), nil)
	return nil
}

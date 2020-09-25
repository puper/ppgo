package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris/v12"

	"github.com/puper/ppgo/v2/components/db"
	"github.com/puper/ppgo/v2/components/irisapp"
	"github.com/puper/ppgo/v2/components/log"

	"github.com/puper/ppgo/v2/engine"
	"github.com/spf13/viper"
)

func main() {
	config := viper.New()
	config.SetConfigFile("/Users/puper/go/src/github.com/puper/ppgo/v2/examples/config.toml")
	config.ReadInConfig()
	e := engine.New(config)
	e.Register("log", log.Builder("log"))
	e.Register("db", db.Builder("db", nil), "log")
	e.Register("web", func(e *engine.Engine) (interface{}, error) {
		l, err := net.Listen("tcp4", e.GetConfig().GetString("web.addr"))
		if err != nil {
			return nil, err
		}
		app := &irisapp.Application{
			Application: iris.New(),
		}
		go app.Run(iris.Listener(l))
		return app, nil
	}, "log", "db")
	e.Build()
	defer e.Close()
	stop := make(chan struct{})
	go func() {
		sChan := make(chan os.Signal)
		for {
			signal.Notify(sChan, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
			sig := <-sChan
			switch sig {
			case os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				stop <- struct{}{}
			}

		}
	}()
	<-stop
}

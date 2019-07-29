package irisapp

import (
	"context"

	"github.com/kataras/iris"
)

type Application struct {
	*iris.Application
}

func (this *Application) Close() error {
	return this.Shutdown(context.Background())
}

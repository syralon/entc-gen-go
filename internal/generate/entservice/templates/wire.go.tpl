//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"{{.module}}/internal/conf"
	"{{.module}}/internal/data"
	"{{.module}}/internal/server"
	"{{.module}}/internal/controller"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		controller.ProviderSet,
		newApp,
	))
}

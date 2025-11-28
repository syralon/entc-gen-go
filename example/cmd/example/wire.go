//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/google/wire"

	"github.com/syralon/entc-gen-go/example/internal/conf"
	"github.com/syralon/entc-gen-go/example/internal/data"
	"github.com/syralon/entc-gen-go/example/internal/server"
	"github.com/syralon/entc-gen-go/example/internal/service"
)

// initApp init kratos application.
func initApp(*conf.Server, *conf.Data) (*kratos.App, func(), error) {
	panic(wire.Build(
		server.ProviderSet,
		data.ProviderSet,
		service.ProviderSet,
		newApp,
	))
}

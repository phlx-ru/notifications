//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"context"

	"notifications/internal/pkg/metrics"
	"notifications/internal/senders"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"notifications/internal/biz"
	"notifications/internal/conf"
	"notifications/internal/data"
	"notifications/internal/server"
	"notifications/internal/service"
)

// wireData init database
func wireData(*conf.Data, log.Logger) (data.Database, func(), error) {
	panic(wire.Build(data.ProviderDataSet))
}

// wireApp init kratos application.
func wireApp(context.Context, data.Database, *conf.Server, *conf.Auth, *senders.Senders, metrics.Metrics, log.Logger) (
	*kratos.App,
	error,
) {
	panic(wire.Build(server.ProviderSet, data.ProviderRepoSet, biz.ProviderSet, service.ProviderSet, newApp))
}

//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"notifications/internal/conf"
	"notifications/internal/senders"
	"notifications/internal/worker"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"notifications/internal/biz"
	"notifications/internal/data"
)

func wireWorker(*conf.Data, *senders.Senders, log.Logger) (*worker.Worker, func(), error) {
	panic(wire.Build(data.ProviderDataSet, data.ProviderRepoSet, biz.ProviderSet, newWorker))
}

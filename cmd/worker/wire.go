//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"notifications/internal/conf"
	"notifications/internal/pkg/metrics"
	"notifications/internal/senders"
	"notifications/internal/worker"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	"notifications/internal/biz"
	"notifications/internal/data"
)

// wireData init database
func wireData(*conf.Data, log.Logger) (data.Database, func(), error) {
	panic(wire.Build(data.ProviderDataSet))
}

func wireWorker(data.Database, *senders.Senders, metrics.Metrics, log.Logger) (*worker.Worker, error) {
	panic(wire.Build(data.ProviderRepoSet, biz.ProviderSet, newWorker))
}

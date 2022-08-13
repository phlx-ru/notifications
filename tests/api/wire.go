//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package api

import (
	"notifications/internal/pkg/metrics"
	"notifications/internal/senders"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
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

func wireNotificationRepo(data.Database, log.Logger, metrics.Metrics) biz.NotificationRepo {
	panic(wire.Build(data.ProviderRepoSet))
}

func wireHTTPServer(data.Database, *conf.Server, *senders.Senders, metrics.Metrics, log.Logger) *http.Server {
	panic(wire.Build(server.ProviderSet, data.ProviderRepoSet, biz.ProviderSet, service.ProviderSet))
}

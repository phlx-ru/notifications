package server

import (
	"time"

	v1notification "notifications/api/notification/v1"
	"notifications/internal/conf"
	"notifications/internal/middlewares"
	"notifications/internal/pkg/metrics"
	"notifications/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	c *conf.Server,
	notifier *service.NotificationService,
	metric metrics.Metrics,
	logger log.Logger,
) *http.Server {
	var opts = []http.ServerOption{
		http.Timeout(5 * time.Second), // TODO Check timeout
		http.Middleware(
			middlewares.Duration(metric, logger),
			tracing.Server(),
			recovery.Recovery(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	v1notification.RegisterNotificationHTTPServer(srv, notifier)
	registerPProfHandlers(srv)
	return srv
}

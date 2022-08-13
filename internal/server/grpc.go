package server

import (
	v1notification "notifications/api/notification/v1"
	"notifications/internal/conf"
	"notifications/internal/middlewares"
	"notifications/internal/pkg/metrics"
	"notifications/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(
	c *conf.Server,
	notifier *service.NotificationService,
	metric metrics.Metrics,
	logger log.Logger,
) *grpc.Server {
	var opts = []grpc.ServerOption{
		grpc.Timeout(c.Grpc.Timeout.AsDuration()),
		grpc.Middleware(
			middlewares.Duration(metric, logger),
			tracing.Server(),
			recovery.Recovery(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)
	v1notification.RegisterNotificationServer(srv, notifier)
	return srv
}

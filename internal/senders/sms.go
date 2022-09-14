package senders

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"notifications/internal/clients/smsaero"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"
)

const (
	metricSMSAeroSendTimings = `senders.sms-aero.send.timings`
	metricSMSAeroSendSuccess = `senders.sms-aero.send.success`
	metricSMSAeroSendFailure = `senders.sms-aero.send.failure`
)

type SMSAeroSender interface {
	Send(ctx context.Context, phone, text string) error
}

type SMSAero struct {
	client smsaero.Client
	metric metrics.Metrics
	logs   logger.Logger
}

func NewSMSAero(client smsaero.Client, metric metrics.Metrics, logs log.Logger) *SMSAero {
	return &SMSAero{
		client: client,
		metric: metric,
		logs:   logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "senders-sms-aero"),
	}
}

func (a *SMSAero) Send(ctx context.Context, phone, text string) error {
	defer a.metric.NewTiming().Send(metricSMSAeroSendTimings)

	err := a.client.Send(ctx, phone, text)
	if err != nil {
		a.metric.Increment(metricSMSAeroSendFailure)
		a.logs.WithContext(ctx).Errorf("failed sms-aero notification: %v", err)
	} else {
		a.metric.Increment(metricSMSAeroSendSuccess)
		a.logs.WithContext(ctx).Info("success sms-aero notification")
	}
	return err
}

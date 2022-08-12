package senders

import (
	"context"

	"notifications/internal/clients/telegram"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricTelegramSendSuccess = `senders.telegram.send.success`
	metricTelegramSendFailure = `senders.telegram.send.failure`
	metricTelegramSendTimings = `senders.telegram.send.timings`
)

type TelegramSenderOption func(request *telegram.SendMessageRequest)

// WithParseMode get `markdown` or `html`
func WithParseMode(parseMode string) TelegramSenderOption {
	return func(request *telegram.SendMessageRequest) {
		request.ParseMode = &parseMode
	}
}

func WithDisableWebPagePreview(disableWebPagePreview bool) TelegramSenderOption {
	return func(request *telegram.SendMessageRequest) {
		request.DisableWebPagePreview = &disableWebPagePreview
	}
}

func WithDisableNotification(disableNotification bool) TelegramSenderOption {
	return func(request *telegram.SendMessageRequest) {
		request.DisableNotification = &disableNotification
	}
}

func WithProtectContent(protectContent bool) TelegramSenderOption {
	return func(request *telegram.SendMessageRequest) {
		request.ProtectContent = &protectContent
	}
}

type TelegramSender interface {
	Send(ctx context.Context, chatID, text string, options ...TelegramSenderOption) error
}

type Telegram struct {
	client telegram.Client
	metric metrics.Metrics
	logs   logger.Logger
}

func NewTelegram(client telegram.Client, metric metrics.Metrics, logs log.Logger) *Telegram {
	return &Telegram{
		client: client,
		metric: metric,
		logs:   logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "senders-telegram"),
	}
}

func (t *Telegram) Send(ctx context.Context, chatID, text string, options ...TelegramSenderOption) error {
	defer t.metric.NewTiming().Send(metricTelegramSendTimings)

	request := &telegram.SendMessageRequest{
		ChatID: chatID,
		Text:   text,
	}

	for _, option := range options {
		option(request)
	}

	_, err := t.client.SendMessage(*request)
	if err != nil {
		t.metric.Increment(metricTelegramSendFailure)
		t.logs.WithContext(ctx).Errorf("failed telegram notification: %v", err)
	} else {
		t.metric.Increment(metricTelegramSendSuccess)
		t.logs.WithContext(ctx).Info("success telegram notification")
	}
	return err
}

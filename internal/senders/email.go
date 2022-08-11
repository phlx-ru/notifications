package senders

import (
	"context"
	"fmt"
	"net/smtp"
	"net/url"
	"strings"

	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jordan-wright/email"
)

const (
	metricEmailSendSuccess = `senders.email.send.success`
	metricEmailSendFailure = `senders.email.send.failure`
	metricEmailSendTimings = `senders.email.send.timings`
)

type EmailSender interface {
	SendText(ctx context.Context, to []string, subject, body string) error
	SendHTML(ctx context.Context, to []string, subject, body string) error
}

type Email struct {
	From              string
	RelayAddress      string
	RelayAuthUsername string
	RelayAuthPassword string
	RelayAuthHost     string
	metric            metrics.Metrics
	logs              logger.Logger
}

func (e *Email) SendText(ctx context.Context, to []string, subject, body string) error {
	return e.send(
		ctx,
		&email.Email{
			To:      to,
			From:    e.From,
			Subject: subject,
			Text:    []byte(body),
		},
	)
}

func (e *Email) SendHTML(ctx context.Context, to []string, subject, body string) error {
	return e.send(
		ctx,
		&email.Email{
			To:      to,
			From:    e.From,
			Subject: subject,
			HTML:    []byte(body),
		},
	)
}

func (e *Email) send(ctx context.Context, mail *email.Email) error {
	defer e.metric.NewTiming().Send(metricEmailSendTimings)
	auth := smtp.PlainAuth("", e.RelayAuthUsername, e.RelayAuthPassword, e.RelayAuthHost)
	err := mail.Send(e.RelayAddress, auth)
	if err != nil {
		e.logs.WithContext(ctx).Errorf("failed to send email to %s", strings.Join(mail.To, ", "))
		e.metric.Increment(metricEmailSendFailure)
	} else {
		e.logs.WithContext(ctx).Infof("email sent to %s", strings.Join(mail.To, ", "))
		e.metric.Increment(metricEmailSendSuccess)
	}
	return err
}

func NewEmail(from, addr, username, password string, metric metrics.Metrics, logs log.Logger) (*Email, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse smtp relay address '%s': %v", addr, err)
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("smtp relay address '%s' is incorrect", addr)
	}
	return &Email{
		From:              from,
		RelayAddress:      addr,
		RelayAuthUsername: username,
		RelayAuthPassword: password,
		RelayAuthHost:     u.Scheme,
		metric:            metric,
		logs:              logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "senders-email"),
	}, nil
}

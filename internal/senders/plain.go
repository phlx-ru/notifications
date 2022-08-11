package senders

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	emulateErrors  = false // TODO Move to request
	errorsPercent  = 50
	emulateLatency = false // TODO Move to request
	latencyLimit   = 1000 * time.Millisecond

	metricPlainSendSuccess = `senders.plain.send.success`
	metricPlainSendFailure = `senders.plain.send.failure`
	metricPlainSendTimings = `senders.plain.send.timings`
)

type PlainSender interface {
	Send(ctx context.Context, message string) error
}

type Plain struct {
	metric metrics.Metrics
	writer logger.Logger
	logs   logger.Logger
}

func FromPath(path string) (*os.File, error) {
	//nolint:gosec // G304: Potential file inclusion via variable (gosec)
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
}

func NewPlain(file *os.File, metric metrics.Metrics, logs log.Logger) *Plain {
	writer := log.NewHelper(
		log.With(
			log.NewStdLogger(file),
			"ts", log.DefaultTimestamp,
			"sender", "plain",
		),
	)
	return &Plain{
		writer: writer,
		metric: metric,
		logs:   logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "senders-plain"),
	}
}

func (p *Plain) Send(ctx context.Context, message string) error {
	defer p.metric.NewTiming().Send(metricPlainSendTimings)
	var err error
	if emulateLatency {
		time.Sleep(randomLatency(latencyLimit.Milliseconds()))
	}
	if emulateErrors && randomBool(errorsPercent) {
		err = fmt.Errorf(`plain message failed cause of errors percent is %d`, errorsPercent)
	}
	if err != nil {
		p.metric.Increment(metricPlainSendFailure)
		p.logs.WithContext(ctx).Errorf("failed plain notification '%s': %v", message, err)
	} else {
		p.writer.WithContext(ctx).Infof(`received message [%s]`, message)
		p.metric.Increment(metricPlainSendSuccess)
		p.logs.WithContext(ctx).Infof("success plain notification '%s'", message)
	}
	return err
}

func randomLatency(limitMilliSeconds int64) time.Duration {
	rand.Seed(time.Now().UnixNano())
	//nolint:gosec // G404: Use of weak random number generator (math/rand instead of crypto/rand)
	ms := (rand.Int63() + limitMilliSeconds) % limitMilliSeconds
	return time.Duration(ms) * time.Millisecond
}

func randomBool(successPercent int) bool {
	rand.Seed(time.Now().UnixNano())
	//nolint:gosec // G404: Use of weak random number generator (math/rand instead of crypto/rand)
	return (rand.Uint32() % 100) >= uint32(successPercent)
}

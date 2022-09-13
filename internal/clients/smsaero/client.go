package smsaero

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-kratos/kratos/v2/log"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"
	"notifications/internal/pkg/template"
	"notifications/internal/pkg/transport"
)

const (
	host     = `gate.smsaero.ru`
	sendPath = `v2/sms/send`
	sign     = `SMS Aero`

	// URLTemplate for GET method send SMS via https://smsaero.ru/api/v1/#send-sms
	URLTemplate = `https://{{ .email }}:{{ .apiKey }}@{{ .host }}/{{ .sendPath }}` +
		`?number={{ .number }}&text={{ urlquery .text  }}&sign={{ urlquery .sign }}`

	metricSendSuccess = `clients.sms-aero.send.success`
	metricSendFailure = `clients.sms-aero.send.failure`
	metricSendTimings = `clients.sms-aero.send.timings`
)

type Client interface {
	Send(ctx context.Context, phone, text string) error
}

type SMSAero struct {
	AuthEmail  string
	AuthAPIKey string
	client     transport.HTTPClient
	metric     metrics.Metrics
	logs       logger.Logger
}

func New(authEmail, authAPIKey string, client transport.HTTPClient, metric metrics.Metrics, logs log.Logger) *SMSAero {
	return &SMSAero{
		AuthEmail:  authEmail,
		AuthAPIKey: authAPIKey,
		client:     client,
		metric:     metric,
		logs:       logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "clients-sms-aero"),
	}
}

func (c *SMSAero) Send(ctx context.Context, phone, text string) error {
	defer c.metric.NewTiming().Send(metricSendTimings)
	var err error
	defer func() {
		if err != nil {
			c.metric.Increment(metricSendFailure)
			c.logs.Errorf(`failed to send: %v`, err)
		} else {
			c.metric.Increment(metricSendSuccess)
		}
	}()

	url := template.MustInterpolate(
		URLTemplate, map[string]any{
			"email":    c.AuthEmail,
			"apiKey":   c.AuthAPIKey,
			"host":     host,
			"sendPath": sendPath,
			"number":   phone,
			"text":     text,
			"sign":     sign,
		},
	)

	method := `GET`
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response map[string]any
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return err
	}

	successRaw, ok := response["success"]
	if !ok {
		err = fmt.Errorf("response has not success attribute: %v", string(responseBody))
		return err
	}
	success, ok := successRaw.(bool)
	if !ok {
		err = fmt.Errorf("response success if not bool: %v", string(responseBody))
	}
	if !success {
		err = fmt.Errorf("response is not success: %v", string(responseBody))
	}
	return err
}

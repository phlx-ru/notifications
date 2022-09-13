package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"
	"notifications/internal/pkg/transport"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	baseURLPattern = `https://api.telegram.org/bot%s/%s`

	metricSendMessageSuccess = `clients.telegram.sendMessage.success`
	metricSendMessageFailure = `clients.telegram.sendMessage.failure`
	metricSendMessageTimings = `clients.telegram.sendMessage.timings`
)

type Client interface {
	SendMessage(ctx context.Context, request SendMessageRequest) (*SendMessageResponse, error)
}

type Telegram struct {
	botToken string
	client   transport.HTTPClient
	metric   metrics.Metrics
	logs     logger.Logger
}

// SendMessageRequest based on https://core.telegram.org/bots/api#sendmessage
type SendMessageRequest struct {
	ChatID                string  `json:"chat_id"`                            // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	Text                  string  `json:"text"`                               // Text of the message to be sent, 1-4096 characters after entities parsing
	ParseMode             *string `json:"parse_mode,omitempty"`               // Mode for parsing entities in the message text. See formatting options (https://core.telegram.org/bots/api#formatting-options) for more details.
	DisableWebPagePreview *bool   `json:"disable_web_page_preview,omitempty"` // Disables link previews for links in this message
	DisableNotification   *bool   `json:"disable_notification,omitempty"`     // Sends the message silently (https://telegram.org/blog/channels-2-0#silent-messages). Users will receive a notification with no sound.
	ProtectContent        *bool   `json:"protect_content,omitempty"`          // Protects the contents of the sent message from forwarding and saving
}

// SendMessageResponse based on https://core.telegram.org/bots/api#message
type SendMessageResponse struct {
	MessageID int `json:"message_id"` // Unique message identifier inside this chat
	Date      int `json:"date"`       // Date the message was sent in Unix time
}

func New(botToken string, client transport.HTTPClient, metric metrics.Metrics, logs log.Logger) *Telegram {
	return &Telegram{
		botToken: botToken,
		client:   client,
		metric:   metric,
		logs:     logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "clients-telegram"),
	}
}

func (t *Telegram) SendMessage(ctx context.Context, request SendMessageRequest) (*SendMessageResponse, error) {
	defer t.metric.NewTiming().Send(metricSendMessageTimings)
	var err error
	defer func() {
		if err != nil {
			t.metric.Increment(metricSendMessageFailure)
			t.logs.Errorf(`failed to sendMessage: %v`, err)
		} else {
			t.metric.Increment(metricSendMessageSuccess)
		}
	}()

	url := fmt.Sprintf(baseURLPattern, t.botToken, `sendMessage`)
	method := `POST`

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response *SendMessageResponse
	err = json.Unmarshal(responseBody, &response)
	return response, err
}

package schema

import (
	"errors"
	"fmt"

	"notifications/internal/pkg/slices"
	"notifications/internal/pkg/strings"
)

var (
	TelegramParseModes = []string{`markdown`, `html`}
)

type PayloadTelegram struct {
	PayloadTyped `json:"-"`

	// Attributes based on https://core.telegram.org/bots/api#sendmessage
	ChatID                string `json:"chat_id"`                            // Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	Text                  string `json:"text"`                               // Text of the message to be sent, 1-4096 characters after entities parsing
	ParseMode             string `json:"parse_mode,omitempty"`               // Mode for parsing entities in the message text. See formatting options (https://core.telegram.org/bots/api#formatting-options) for more details.
	DisableWebPagePreview string `json:"disable_web_page_preview,omitempty"` // Disables link previews for links in this message
	DisableNotification   string `json:"disable_notification,omitempty"`     // Sends the message silently (https://telegram.org/blog/channels-2-0#silent-messages). Users will receive a notification with no sound.
	ProtectContent        string `json:"protect_content,omitempty"`          // Protects the contents of the sent message from forwarding and saving
}

func (p Payload) ToPayloadTelegram() (*PayloadTelegram, error) {
	return toPayloadTyped[PayloadTelegram](p)
}

func (pt PayloadTelegram) MustToPayload() Payload {
	return mustToPayloadCommon(pt)
}

func (pt PayloadTelegram) Validate() error {
	if pt.ChatID == "" {
		return errors.New(`payload telegram has empty field 'chat_id'`)
	}
	if pt.Text == "" {
		return errors.New(`payload telegram has empty field 'text'`)
	}
	if pt.ParseMode != "" && !slices.Includes(pt.ParseMode, TelegramParseModes) {
		return fmt.Errorf("payload telegram has unknown value of 'parse_mode': %s", pt.ParseMode)
	}
	if pt.DisableWebPagePreview != "" && !strings.IsBool(pt.DisableWebPagePreview) {
		return fmt.Errorf(
			"payload telegram has incorrect boolean value of 'disable_web_page_preview': %s",
			pt.DisableWebPagePreview,
		)
	}
	if pt.DisableNotification != "" && !strings.IsBool(pt.DisableNotification) {
		return fmt.Errorf(
			"payload telegram has incorrect boolean value of 'disable_notification': %s",
			pt.DisableNotification,
		)
	}
	if pt.ProtectContent != "" && !strings.IsBool(pt.ProtectContent) {
		return fmt.Errorf(
			"payload telegram has incorrect boolean value of 'protect_content': %s",
			pt.ProtectContent,
		)
	}
	return nil
}

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

	ChatID                string `json:"chat_id"`
	Text                  string `json:"text"`
	ParseMode             string `json:"parse_mode,omitempty"`
	DisableWebPagePreview string `json:"disable_web_page_preview,omitempty"`
	DisableNotification   string `json:"disable_notification,omitempty"`
	ProtectContent        string `json:"protect_content,omitempty"`
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

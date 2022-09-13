package schema

import (
	"errors"
	"fmt"
	"strings"

	"notifications/internal/pkg/sms"
	pkgStrings "notifications/internal/pkg/strings"
)

const (
	messageLimit = 2000
)

type PayloadSMS struct {
	PayloadTyped `json:"-"`

	Phone string `json:"phone"`           // Phone number in format 79009009090
	Text  string `json:"text"`            // Text of SMS message, limit of 160 symbols for latin and 70 symbols for cyrillic https://www.twilio.com/docs/glossary/what-sms-character-limit
	Split string `json:"split,omitempty"` // Split to few messages if text length exceeds limit
}

func (p Payload) ToPayloadSMS() (*PayloadSMS, error) {
	return toPayloadTyped[PayloadSMS](p)
}

func (ps PayloadSMS) MustToPayload() Payload {
	return mustToPayloadCommon(ps)
}

func (ps PayloadSMS) Validate() error {
	if ps.Text == "" {
		return errors.New(`payload sms has empty field text`)
	}
	if len(ps.Text) > messageLimit {
		return fmt.Errorf(`message exceeds symbols limit of %d`, messageLimit)
	}
	if ps.Phone == "" {
		return errors.New(`payload sms has empty field text`)
	}
	phoneExample := `79009009090`
	if len(ps.Phone) != len(phoneExample) {
		return fmt.Errorf(`phone must be in format %s`, phoneExample)
	}
	if !strings.HasPrefix(ps.Phone, `79`) {
		return fmt.Errorf(`phone must be start with 79`)
	}
	if ps.Split != "" && !pkgStrings.IsBool(ps.Split) {
		return fmt.Errorf("split is not contain bool value: %s", ps.Split)
	}
	if !pkgStrings.IsTrue(ps.Split) && sms.IsExceedsLimit(ps.Text) {
		return errors.New(`message exceeds symbols limit, add {"split":"true"} for split to few messages`)
	}
	return nil
}

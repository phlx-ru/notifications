package schema

import (
	"errors"
	"fmt"
	"net/mail"

	"notifications/internal/pkg/strings"
)

type PayloadEmail struct {
	PayloadTyped `json:"-"`

	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	IsHTML  string `json:"is_html"`
}

func (p Payload) ToPayloadEmail() (*PayloadEmail, error) {
	return toPayloadTyped[PayloadEmail](p)
}

func (pe PayloadEmail) MustToPayload() Payload {
	return mustToPayloadCommon(pe)
}

func (pe PayloadEmail) Validate() error {
	if pe.To == "" {
		return errors.New(`payload email has empty field 'to'`)
	}
	_, err := mail.ParseAddress(pe.To)
	if err != nil {
		return fmt.Errorf(`email '%s' is invalid: %w`, pe.To, err)
	}
	if pe.Subject == "" {
		return errors.New(`payload email has empty field 'subject'`)
	}
	if pe.Body == "" {
		return errors.New(`payload email has empty field 'body'`)
	}
	if pe.IsHTML != "" && !strings.IsBool(pe.IsHTML) {
		return fmt.Errorf(`payload email has incorrect boolean value for field 'is_html': %s`, pe.IsHTML)
	}
	return nil
}

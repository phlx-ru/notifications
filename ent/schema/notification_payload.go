package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
)

type Payload []byte

func (p Payload) String() string {
	return string(p)
}

func PayloadFromProto(proto map[string]string) (*Payload, error) {
	bytes, err := json.Marshal(proto)
	if err != nil {
		return nil, err
	}
	payload := Payload(bytes)
	return &payload, nil
}

type PayloadEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	IsHTML  bool   `json:"is_html"`
}

func (p Payload) ToPayloadEmail() (*PayloadEmail, error) {
	var pe PayloadEmail
	err := json.Unmarshal(p, &pe)
	if err != nil {
		return nil, err
	}
	return &pe, nil
}

func (pe *PayloadEmail) MustToPayload() Payload {
	bytes, err := json.Marshal(pe)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (pe *PayloadEmail) Validate() error {
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
	return nil
}

type PayloadPlain struct {
	Message string `json:"message"`
}

func (p Payload) ToPayloadPlain() (*PayloadPlain, error) {
	var pp PayloadPlain
	err := json.Unmarshal(p, &pp)
	if err != nil {
		return nil, err
	}
	return &pp, nil
}

func (pp *PayloadPlain) MustToPayload() Payload {
	bytes, err := json.Marshal(pp)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (pp *PayloadPlain) Validate() error {
	if pp.Message == "" {
		return errors.New(`payload plain has empty field 'message'`)
	}
	return nil
}

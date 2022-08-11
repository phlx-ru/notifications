package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
)

const (
	marshalErrString = "<failed to marshal payload>"
)

type Payload map[string]string

func (p Payload) String() string {
	s, err := json.Marshal(p)
	if err != nil {
		return marshalErrString
	}
	return string(s)
}

func PayloadFromProto(proto map[string]string) (*Payload, error) {
	payload := Payload(proto)
	return &payload, nil
}

type PayloadEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	IsHTML  string `json:"is_html"`
}

func toPayloadTyped[T PayloadPlain | PayloadEmail](source Payload) (*T, error) {
	var res T
	marshaled := source.String()
	if marshaled == marshalErrString {
		return nil, errors.New(marshalErrString)
	}
	err := json.Unmarshal([]byte(marshaled), &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func mustToPayloadCommon(source any) Payload {
	bytes, err := json.Marshal(source)
	if err != nil {
		panic(err)
	}
	var res Payload
	if err := json.Unmarshal(bytes, &res); err != nil {
		panic(err)
	}
	return res
}

func (p Payload) ToPayloadEmail() (*PayloadEmail, error) {
	return toPayloadTyped[PayloadEmail](p)
}

func (pe *PayloadEmail) MustToPayload() Payload {
	return mustToPayloadCommon(pe)
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
	return toPayloadTyped[PayloadPlain](p)
}

func (pp *PayloadPlain) MustToPayload() Payload {
	return mustToPayloadCommon(pp)
}

func (pp *PayloadPlain) Validate() error {
	if pp.Message == "" {
		return errors.New(`payload plain has empty field 'message'`)
	}
	return nil
}

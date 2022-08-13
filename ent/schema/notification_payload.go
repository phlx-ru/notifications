package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	marshalErrStringPattern = "<failed to marshal payload:%s>"
)

type Payload map[string]string

func TypeToValidatorMap(p *Payload) map[NotificationType]PayloadTypedValidator {
	return map[NotificationType]PayloadTypedValidator{
		TypePlain:    ToPayloadTypedValidator(p.ToPayloadEmail()),
		TypeEmail:    ToPayloadTypedValidator(p.ToPayloadEmail()),
		TypeTelegram: ToPayloadTypedValidator(p.ToPayloadTelegram()),
	}
}

type PayloadTypedValidator func() error

func ToPayloadTypedValidator(typed PayloadTyped, err error) PayloadTypedValidator {
	return func() error {
		if err != nil {
			return err
		}
		return typed.Validate()
	}
}

func (p Payload) String() string {
	s, err := json.Marshal(p)
	if err != nil {
		return fmt.Sprintf(marshalErrStringPattern, err.Error())
	}
	return string(s)
}

func (p Payload) Validate(as NotificationType) error {
	validators := TypeToValidatorMap(&p)

	validate, ok := validators[as]
	if !ok {
		return fmt.Errorf(`unknown or unimplemented type of notification '%s'`, as)
	}

	return validate()
}

func PayloadFromProto(proto map[string]string) (*Payload, error) {
	payload := Payload(proto)
	return &payload, nil
}

func toPayloadTyped[T PayloadTyped](source Payload) (*T, error) {
	var res T
	marshaled := source.String()
	marshalErrPrefix := strings.Split(marshalErrStringPattern, "%s")[0]
	if strings.HasPrefix(marshaled, marshalErrPrefix) {
		return nil, errors.New(marshaled)
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

type PayloadTyped interface {
	MustToPayload() Payload
	Validate() error
}

package schema

import "errors"

type PayloadPlain struct {
	PayloadTyped `json:"-"`

	Message string `json:"message"`
}

func (p Payload) ToPayloadPlain() (*PayloadPlain, error) {
	return toPayloadTyped[PayloadPlain](p)
}

func (pp PayloadPlain) MustToPayload() Payload {
	return mustToPayloadCommon(pp)
}

func (pp PayloadPlain) Validate() error {
	if pp.Message == "" {
		return errors.New(`payload plain has empty field 'message'`)
	}
	return nil
}

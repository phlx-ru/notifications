package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	kratosErrors "github.com/go-kratos/kratos/v2/errors"
)

type NotificationType string

type NotificationStatus string

const (
	TypeEmail    NotificationType = `email`
	TypeSMS      NotificationType = `sms`
	TypePush     NotificationType = `push`
	TypeWhatsApp NotificationType = `whatsapp`
	TypeTelegram NotificationType = `telegram`

	StatusDraft   NotificationStatus = `draft`
	StatusPending NotificationStatus = `pending`
	StatusSent    NotificationStatus = `sent`
	StatusRetry   NotificationStatus = `retry`
	StatusFail    NotificationStatus = `fail`
)

var (
	Types = []NotificationType{
		TypeEmail,
		TypeSMS,
		TypePush,
		TypeWhatsApp,
		TypeTelegram,
	}

	Statuses = []NotificationStatus{
		StatusDraft,
		StatusPending,
		StatusSent,
		StatusRetry,
		StatusFail,
	}
)

// Notification holds the schema definition for the Notification entity.
type Notification struct {
	ent.Schema
}

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

func (p Payload) ToPayloadEmail() (*PayloadEmail, error) {
	var pe PayloadEmail
	err := json.Unmarshal(p, &pe)
	if err != nil {
		return nil, err
	}
	return &pe, nil
}

type PayloadEmail struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	IsHTML  bool   `json:"is_html"`
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

// Fields of the Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.Int("sender_id"),

		field.String("type").
			Default(TypeEmail.String()).
			Validate(ValidateType).
			GoType(NotificationType(``)).
			Comment("types in (sms|email|whatsapp|push)"),

		field.JSON("payload", Payload{}).
			Comment("message payload as map<string, any> variable by type"),

		field.Int("ttl").
			Comment("time to live in seconds"),

		field.String("status").
			Default(StatusDraft.String()).
			Validate(ValidateStatus).
			GoType(NotificationStatus(``)).
			Comment("statuses in (draft|pending|sent|retry|fail)"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("creation time of notification"),

		field.Time("updated_at").
			Default(time.Now).
			Comment("last update time of notification"),

		field.Time("planned_at").
			Default(time.Now).
			Comment("time for start sending this notification"),

		field.Int("retries").
			Default(0).
			Comment("count of retries to send notification"),

		field.Time("sent_at").
			Optional().
			Nillable().
			Comment("time of notification was sent"),
	}
}

// Indexes of the schema.
func (Notification) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("type"),
		index.Fields("status"),
		index.Fields("planned_at"),
		index.Fields("sent_at"),
	}
}

// Edges of the Notification.
func (Notification) Edges() []ent.Edge {
	return nil
}

func (t NotificationType) String() string {
	return string(t)
}

func ValidateType(notificationType string) error {
	for _, validType := range Types {
		if NotificationType(notificationType) == validType {
			return nil
		}
	}

	return kratosErrors.Newf(
		http.StatusBadRequest, `VALIDATION_ERROR`,
		`invalid notification type: %s`, notificationType,
	)
}

func (s NotificationStatus) String() string {
	return string(s)
}

func ValidateStatus(notificationStatus string) error {
	for _, validStatus := range Statuses {
		if NotificationStatus(notificationStatus) == validStatus {
			return nil
		}
	}
	return kratosErrors.Newf(
		http.StatusBadRequest, `VALIDATION_ERROR`,
		`invalid notification status: %s`, notificationStatus,
	)
}

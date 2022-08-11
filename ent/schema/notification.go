package schema

import (
	"net/http"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	kratosErrors "github.com/go-kratos/kratos/v2/errors"
)

type NotificationType string

type NotificationStatus string

const (
	TypePlain    NotificationType = `plain`
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
		TypePlain,
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

// Fields of the Notification.
func (Notification) Fields() []ent.Field {
	return []ent.Field{
		field.Int("sender_id"),

		field.String("type").
			Default(TypeEmail.String()).
			Validate(ValidateType).
			GoType(NotificationType(``)).
			Comment("types in (plain|sms|email|whatsapp|push)"),

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
			UpdateDefault(time.Now).
			Annotations(
				&entsql.Annotation{
					Default: "CURRENT_TIMESTAMP",
				},
			).
			Comment("last update time of notification"),

		field.Time("planned_at").
			Default(time.Now).
			Comment("time for start sending this notification"),

		field.Time("retry_at").
			Optional().
			Nillable().
			Comment("time for new try to send this notification"),

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

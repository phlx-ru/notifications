// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"notifications/ent/notification"
	"notifications/ent/schema"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
)

// Notification is the model entity for the Notification schema.
type Notification struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// SenderID holds the value of the "sender_id" field.
	SenderID int `json:"sender_id,omitempty"`
	// types in (plain|sms|email|whatsapp|push)
	Type schema.NotificationType `json:"type,omitempty"`
	// message payload as map<string, any> variable by type
	Payload schema.Payload `json:"payload,omitempty"`
	// time to live in seconds
	TTL int `json:"ttl,omitempty"`
	// statuses in (draft|pending|sent|retry|fail)
	Status schema.NotificationStatus `json:"status,omitempty"`
	// creation time of notification
	CreatedAt time.Time `json:"created_at,omitempty"`
	// last update time of notification
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// time for start sending this notification
	PlannedAt time.Time `json:"planned_at,omitempty"`
	// time for new try to send this notification
	RetryAt *time.Time `json:"retry_at,omitempty"`
	// count of retries to send notification
	Retries int `json:"retries,omitempty"`
	// time of notification was sent
	SentAt *time.Time `json:"sent_at,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Notification) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case notification.FieldPayload:
			values[i] = new([]byte)
		case notification.FieldID, notification.FieldSenderID, notification.FieldTTL, notification.FieldRetries:
			values[i] = new(sql.NullInt64)
		case notification.FieldType, notification.FieldStatus:
			values[i] = new(sql.NullString)
		case notification.FieldCreatedAt, notification.FieldUpdatedAt, notification.FieldPlannedAt, notification.FieldRetryAt, notification.FieldSentAt:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type Notification", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Notification fields.
func (n *Notification) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case notification.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			n.ID = int(value.Int64)
		case notification.FieldSenderID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field sender_id", values[i])
			} else if value.Valid {
				n.SenderID = int(value.Int64)
			}
		case notification.FieldType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[i])
			} else if value.Valid {
				n.Type = schema.NotificationType(value.String)
			}
		case notification.FieldPayload:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field payload", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &n.Payload); err != nil {
					return fmt.Errorf("unmarshal field payload: %w", err)
				}
			}
		case notification.FieldTTL:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field ttl", values[i])
			} else if value.Valid {
				n.TTL = int(value.Int64)
			}
		case notification.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				n.Status = schema.NotificationStatus(value.String)
			}
		case notification.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				n.CreatedAt = value.Time
			}
		case notification.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				n.UpdatedAt = value.Time
			}
		case notification.FieldPlannedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field planned_at", values[i])
			} else if value.Valid {
				n.PlannedAt = value.Time
			}
		case notification.FieldRetryAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field retry_at", values[i])
			} else if value.Valid {
				n.RetryAt = new(time.Time)
				*n.RetryAt = value.Time
			}
		case notification.FieldRetries:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field retries", values[i])
			} else if value.Valid {
				n.Retries = int(value.Int64)
			}
		case notification.FieldSentAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field sent_at", values[i])
			} else if value.Valid {
				n.SentAt = new(time.Time)
				*n.SentAt = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this Notification.
// Note that you need to call Notification.Unwrap() before calling this method if this Notification
// was returned from a transaction, and the transaction was committed or rolled back.
func (n *Notification) Update() *NotificationUpdateOne {
	return (&NotificationClient{config: n.config}).UpdateOne(n)
}

// Unwrap unwraps the Notification entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (n *Notification) Unwrap() *Notification {
	_tx, ok := n.config.driver.(*txDriver)
	if !ok {
		panic("ent: Notification is not a transactional entity")
	}
	n.config.driver = _tx.drv
	return n
}

// String implements the fmt.Stringer.
func (n *Notification) String() string {
	var builder strings.Builder
	builder.WriteString("Notification(")
	builder.WriteString(fmt.Sprintf("id=%v, ", n.ID))
	builder.WriteString("sender_id=")
	builder.WriteString(fmt.Sprintf("%v", n.SenderID))
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", n.Type))
	builder.WriteString(", ")
	builder.WriteString("payload=")
	builder.WriteString(fmt.Sprintf("%v", n.Payload))
	builder.WriteString(", ")
	builder.WriteString("ttl=")
	builder.WriteString(fmt.Sprintf("%v", n.TTL))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", n.Status))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(n.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(n.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("planned_at=")
	builder.WriteString(n.PlannedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	if v := n.RetryAt; v != nil {
		builder.WriteString("retry_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteString(", ")
	builder.WriteString("retries=")
	builder.WriteString(fmt.Sprintf("%v", n.Retries))
	builder.WriteString(", ")
	if v := n.SentAt; v != nil {
		builder.WriteString("sent_at=")
		builder.WriteString(v.Format(time.ANSIC))
	}
	builder.WriteByte(')')
	return builder.String()
}

// Notifications is a parsable slice of Notification.
type Notifications []*Notification

func (n Notifications) config(cfg config) {
	for _i := range n {
		n[_i].config = cfg
	}
}

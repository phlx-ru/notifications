// Code generated by ent, DO NOT EDIT.

package notification

import (
	"notifications/ent/predicate"
	"notifications/ent/schema"
	"time"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// SenderID applies equality check predicate on the "sender_id" field. It's identical to SenderIDEQ.
func SenderID(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSenderID), v))
	})
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), vc))
	})
}

// TTL applies equality check predicate on the "ttl" field. It's identical to TTLEQ.
func TTL(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTTL), v))
	})
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), vc))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// PlannedAt applies equality check predicate on the "planned_at" field. It's identical to PlannedAtEQ.
func PlannedAt(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPlannedAt), v))
	})
}

// RetryAt applies equality check predicate on the "retry_at" field. It's identical to RetryAtEQ.
func RetryAt(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRetryAt), v))
	})
}

// Retries applies equality check predicate on the "retries" field. It's identical to RetriesEQ.
func Retries(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRetries), v))
	})
}

// SentAt applies equality check predicate on the "sent_at" field. It's identical to SentAtEQ.
func SentAt(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSentAt), v))
	})
}

// SenderIDEQ applies the EQ predicate on the "sender_id" field.
func SenderIDEQ(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSenderID), v))
	})
}

// SenderIDNEQ applies the NEQ predicate on the "sender_id" field.
func SenderIDNEQ(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSenderID), v))
	})
}

// SenderIDIn applies the In predicate on the "sender_id" field.
func SenderIDIn(vs ...int) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldSenderID), v...))
	})
}

// SenderIDNotIn applies the NotIn predicate on the "sender_id" field.
func SenderIDNotIn(vs ...int) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldSenderID), v...))
	})
}

// SenderIDGT applies the GT predicate on the "sender_id" field.
func SenderIDGT(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSenderID), v))
	})
}

// SenderIDGTE applies the GTE predicate on the "sender_id" field.
func SenderIDGTE(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSenderID), v))
	})
}

// SenderIDLT applies the LT predicate on the "sender_id" field.
func SenderIDLT(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSenderID), v))
	})
}

// SenderIDLTE applies the LTE predicate on the "sender_id" field.
func SenderIDLTE(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSenderID), v))
	})
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldType), vc))
	})
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldType), vc))
	})
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...schema.NotificationType) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldType), v...))
	})
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...schema.NotificationType) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldType), v...))
	})
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldType), vc))
	})
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldType), vc))
	})
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldType), vc))
	})
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldType), vc))
	})
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldType), vc))
	})
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldType), vc))
	})
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldType), vc))
	})
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldType), vc))
	})
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v schema.NotificationType) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldType), vc))
	})
}

// TTLEQ applies the EQ predicate on the "ttl" field.
func TTLEQ(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTTL), v))
	})
}

// TTLNEQ applies the NEQ predicate on the "ttl" field.
func TTLNEQ(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTTL), v))
	})
}

// TTLIn applies the In predicate on the "ttl" field.
func TTLIn(vs ...int) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldTTL), v...))
	})
}

// TTLNotIn applies the NotIn predicate on the "ttl" field.
func TTLNotIn(vs ...int) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldTTL), v...))
	})
}

// TTLGT applies the GT predicate on the "ttl" field.
func TTLGT(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTTL), v))
	})
}

// TTLGTE applies the GTE predicate on the "ttl" field.
func TTLGTE(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTTL), v))
	})
}

// TTLLT applies the LT predicate on the "ttl" field.
func TTLLT(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTTL), v))
	})
}

// TTLLTE applies the LTE predicate on the "ttl" field.
func TTLLTE(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTTL), v))
	})
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), vc))
	})
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStatus), vc))
	})
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...schema.NotificationStatus) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldStatus), v...))
	})
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...schema.NotificationStatus) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = string(vs[i])
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldStatus), v...))
	})
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStatus), vc))
	})
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStatus), vc))
	})
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStatus), vc))
	})
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStatus), vc))
	})
}

// StatusContains applies the Contains predicate on the "status" field.
func StatusContains(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldStatus), vc))
	})
}

// StatusHasPrefix applies the HasPrefix predicate on the "status" field.
func StatusHasPrefix(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldStatus), vc))
	})
}

// StatusHasSuffix applies the HasSuffix predicate on the "status" field.
func StatusHasSuffix(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldStatus), vc))
	})
}

// StatusEqualFold applies the EqualFold predicate on the "status" field.
func StatusEqualFold(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldStatus), vc))
	})
}

// StatusContainsFold applies the ContainsFold predicate on the "status" field.
func StatusContainsFold(v schema.NotificationStatus) predicate.Notification {
	vc := string(v)
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldStatus), vc))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// PlannedAtEQ applies the EQ predicate on the "planned_at" field.
func PlannedAtEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldPlannedAt), v))
	})
}

// PlannedAtNEQ applies the NEQ predicate on the "planned_at" field.
func PlannedAtNEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldPlannedAt), v))
	})
}

// PlannedAtIn applies the In predicate on the "planned_at" field.
func PlannedAtIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldPlannedAt), v...))
	})
}

// PlannedAtNotIn applies the NotIn predicate on the "planned_at" field.
func PlannedAtNotIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldPlannedAt), v...))
	})
}

// PlannedAtGT applies the GT predicate on the "planned_at" field.
func PlannedAtGT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldPlannedAt), v))
	})
}

// PlannedAtGTE applies the GTE predicate on the "planned_at" field.
func PlannedAtGTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldPlannedAt), v))
	})
}

// PlannedAtLT applies the LT predicate on the "planned_at" field.
func PlannedAtLT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldPlannedAt), v))
	})
}

// PlannedAtLTE applies the LTE predicate on the "planned_at" field.
func PlannedAtLTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldPlannedAt), v))
	})
}

// RetryAtEQ applies the EQ predicate on the "retry_at" field.
func RetryAtEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRetryAt), v))
	})
}

// RetryAtNEQ applies the NEQ predicate on the "retry_at" field.
func RetryAtNEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRetryAt), v))
	})
}

// RetryAtIn applies the In predicate on the "retry_at" field.
func RetryAtIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldRetryAt), v...))
	})
}

// RetryAtNotIn applies the NotIn predicate on the "retry_at" field.
func RetryAtNotIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldRetryAt), v...))
	})
}

// RetryAtGT applies the GT predicate on the "retry_at" field.
func RetryAtGT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRetryAt), v))
	})
}

// RetryAtGTE applies the GTE predicate on the "retry_at" field.
func RetryAtGTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRetryAt), v))
	})
}

// RetryAtLT applies the LT predicate on the "retry_at" field.
func RetryAtLT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRetryAt), v))
	})
}

// RetryAtLTE applies the LTE predicate on the "retry_at" field.
func RetryAtLTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRetryAt), v))
	})
}

// RetryAtIsNil applies the IsNil predicate on the "retry_at" field.
func RetryAtIsNil() predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldRetryAt)))
	})
}

// RetryAtNotNil applies the NotNil predicate on the "retry_at" field.
func RetryAtNotNil() predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldRetryAt)))
	})
}

// RetriesEQ applies the EQ predicate on the "retries" field.
func RetriesEQ(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldRetries), v))
	})
}

// RetriesNEQ applies the NEQ predicate on the "retries" field.
func RetriesNEQ(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldRetries), v))
	})
}

// RetriesIn applies the In predicate on the "retries" field.
func RetriesIn(vs ...int) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldRetries), v...))
	})
}

// RetriesNotIn applies the NotIn predicate on the "retries" field.
func RetriesNotIn(vs ...int) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldRetries), v...))
	})
}

// RetriesGT applies the GT predicate on the "retries" field.
func RetriesGT(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldRetries), v))
	})
}

// RetriesGTE applies the GTE predicate on the "retries" field.
func RetriesGTE(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldRetries), v))
	})
}

// RetriesLT applies the LT predicate on the "retries" field.
func RetriesLT(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldRetries), v))
	})
}

// RetriesLTE applies the LTE predicate on the "retries" field.
func RetriesLTE(v int) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldRetries), v))
	})
}

// SentAtEQ applies the EQ predicate on the "sent_at" field.
func SentAtEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSentAt), v))
	})
}

// SentAtNEQ applies the NEQ predicate on the "sent_at" field.
func SentAtNEQ(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSentAt), v))
	})
}

// SentAtIn applies the In predicate on the "sent_at" field.
func SentAtIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.In(s.C(FieldSentAt), v...))
	})
}

// SentAtNotIn applies the NotIn predicate on the "sent_at" field.
func SentAtNotIn(vs ...time.Time) predicate.Notification {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotIn(s.C(FieldSentAt), v...))
	})
}

// SentAtGT applies the GT predicate on the "sent_at" field.
func SentAtGT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSentAt), v))
	})
}

// SentAtGTE applies the GTE predicate on the "sent_at" field.
func SentAtGTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSentAt), v))
	})
}

// SentAtLT applies the LT predicate on the "sent_at" field.
func SentAtLT(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSentAt), v))
	})
}

// SentAtLTE applies the LTE predicate on the "sent_at" field.
func SentAtLTE(v time.Time) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSentAt), v))
	})
}

// SentAtIsNil applies the IsNil predicate on the "sent_at" field.
func SentAtIsNil() predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldSentAt)))
	})
}

// SentAtNotNil applies the NotNil predicate on the "sent_at" field.
func SentAtNotNil() predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldSentAt)))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Notification) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Notification) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Notification) predicate.Notification {
	return predicate.Notification(func(s *sql.Selector) {
		p(s.Not())
	})
}

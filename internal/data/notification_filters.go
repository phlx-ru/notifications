package data

import (
	"time"

	"notifications/ent/predicate"
	"notifications/ent/schema"

	entSql "entgo.io/ent/dialect/sql"
)

func FilterByType(types ...schema.NotificationType) predicate.Notification {
	return func(selector *entSql.Selector) {
		var args []any
		for _, typ := range types {
			args = append(args, typ)
		}
		selector.Where(entSql.P().In(`type`, args...))
	}
}

func FilterByStatus(statuses ...schema.NotificationStatus) predicate.Notification {
	return func(selector *entSql.Selector) {
		var args []any
		for _, status := range statuses {
			args = append(args, status)
		}
		selector.Where(entSql.P().In(`status`, args...))
	}
}

func FilterByPlannedAt(plannedAt time.Time) predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.Where(entSql.GTE(`planned_at`, plannedAt))
	}
}

func FilterForUpdateWithSkipLocked() predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.ForUpdate(entSql.WithLockAction(entSql.SkipLocked))
	}
}

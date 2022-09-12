package data

import (
	"time"

	"notifications/ent/predicate"
	"notifications/ent/schema"

	entSql "entgo.io/ent/dialect/sql"
)

func FilterByID(id int) predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.Where(entSql.P().EQ(`id`, id))
	}
}

func FilterByType(types ...schema.NotificationType) predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.Where(entSql.P().In(`type`, itemsToAny(types)...))
	}
}

func FilterByStatus(statuses ...schema.NotificationStatus) predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.Where(entSql.P().In(`status`, itemsToAny(statuses)...))
	}
}

// FilterByPlannedAt deprecated
func FilterByPlannedAt(plannedAt time.Time) predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.Where(entSql.LTE(`planned_at`, plannedAt))
	}
}

func FilterByPlannedAtOrRetryAt(now time.Time) predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.
			Where(
				entSql.Or(
					entSql.And(
						entSql.IsNull(`retry_at`),
						entSql.LTE(`planned_at`, now),
					),
					entSql.LTE(`retry_at`, now),
				),
			)
	}
}

func FilterForUpdateWithSkipLocked() predicate.Notification {
	return func(selector *entSql.Selector) {
		selector.ForUpdate(entSql.WithLockAction(entSql.SkipLocked))
	}
}

func itemsToAny[T comparable](items []T) []any {
	res := []any{}
	for _, item := range items {
		res = append(res, item)
	}
	return res
}

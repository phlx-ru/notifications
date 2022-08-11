package data

import (
	"context"
	databaseSql "database/sql"
	"errors"
	"runtime/debug"
	"time"

	"notifications/ent"
	"notifications/ent/schema"
	"notifications/internal/biz"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	metricSaveTimings                             = `data.notification.save.timings`
	metricUpdateTimings                           = `data.notification.update.timings`
	metricFindByIDTimings                         = `data.notification.findById.timings`
	metricCountWaitingNotificationsTimings        = `data.notification.countWaitingNotifications.timings`
	metricListWaitingNotificationsWithLockTimings = `data.notification.listWaitingNotificationsWithLock.timings`
	metricTransactionTimings                      = `data.notification.transaction.timings`
)

type notificationRepo struct {
	data   *Data
	metric metrics.Metrics
	logs   *log.Helper
}

// NewNotificationRepo .
func NewNotificationRepo(data *Data, logs log.Logger, metric metrics.Metrics) biz.NotificationRepo {
	return &notificationRepo{
		data:   data,
		metric: metric,
		logs:   logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "data-notification"),
	}
}

func (r *notificationRepo) Save(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
	defer r.metric.NewTiming().Send(metricSaveTimings)
	if n == nil {
		return nil, errors.New("notification is empty")
	}

	creating := r.client(ctx).Notification.Create().
		SetSenderID(n.SenderID).
		SetType(n.Type).
		SetPayload(n.Payload).
		SetTTL(n.TTL).
		SetStatus(n.Status).
		SetPlannedAt(n.PlannedAt).
		SetRetries(n.Retries)

	if n.SentAt != nil {
		creating = creating.SetSentAt(*n.SentAt)
	}

	if n.RetryAt != nil {
		creating = creating.SetRetryAt(*n.RetryAt)
	}

	return creating.Save(ctx)
}

func (r *notificationRepo) Update(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
	defer r.metric.NewTiming().Send(metricUpdateTimings)
	if n == nil {
		return nil, errors.New("notification is empty")
	}
	updating := r.client(ctx).Notification.
		UpdateOne(n).
		SetSenderID(n.SenderID).
		SetType(n.Type).
		SetPayload(n.Payload).
		SetTTL(n.TTL).
		SetStatus(n.Status).
		SetPlannedAt(n.PlannedAt).
		SetRetries(n.Retries)

	if n.SentAt != nil {
		updating = updating.SetSentAt(*n.SentAt)
	}

	if n.RetryAt != nil {
		updating = updating.SetRetryAt(*n.RetryAt)
	}

	return updating.Save(ctx)
}

func (r *notificationRepo) FindByID(ctx context.Context, id int64) (*ent.Notification, error) {
	defer r.metric.NewTiming().Send(metricFindByIDTimings)
	return r.client(ctx).Notification.Get(ctx, int(id))
}

func (r *notificationRepo) CountWaitingNotifications(ctx context.Context) (int, error) {
	defer r.metric.NewTiming().Send(metricCountWaitingNotificationsTimings)
	return r.client(ctx).Notification.Query().
		Where(
			FilterByStatus(schema.StatusPending, schema.StatusRetry),
			FilterByType(schema.Types...),
			FilterByPlannedAtOrRetryAt(time.Now()),
		).
		Count(ctx)
}

func (r *notificationRepo) ListWaitingNotificationsWithLock(ctx context.Context, limit int) (
	[]*ent.Notification,
	error,
) {
	defer r.metric.NewTiming().Send(metricListWaitingNotificationsWithLockTimings)
	return r.client(ctx).Notification.Query().
		Where(
			FilterByStatus(schema.StatusPending, schema.StatusRetry),
			FilterByType(schema.Types...),
			FilterByPlannedAtOrRetryAt(time.Now()),
			FilterForUpdateWithSkipLocked(),
		).
		Order(OrderByCreatedAt()).
		Limit(limit).
		Unique(false). // Cause: FOR UPDATE is not allowed with DISTINCT clause
		All(ctx)
}

func (r *notificationRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	defer r.metric.NewTiming().Send(metricTransactionTimings)
	tx, err := r.data.ent.BeginTx(ctx, txOptions)
	if err != nil {
		r.logs.Errorf(`failed to start tx: %v`, err)
		return err
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			r.logs.Errorf(`tx panic: recovered = %v; stack = %v`, recovered, string(debug.Stack()))
			if tx != nil {
				if err := tx.Rollback(); err != nil {
					r.logs.Errorf(`tx panic rollback error: %v`, err)
				}
			}
		}
	}()
	repoCtx := ent.NewContext(ctx, tx.Client())
	for _, process := range processes {
		if err := process(repoCtx); err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				r.logs.Errorf(`failed to rollback tx caused of err '%s' because of: %v`, err.Error(), rollbackErr)
				return rollbackErr
			}
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		r.logs.Errorf(`failed to commit tx: %v`, err)
		return err
	}
	return nil
}

// client return client by tx in context if it exists or default ent client
func (r *notificationRepo) client(ctx context.Context) *ent.Client {
	if client := ent.FromContext(ctx); client != nil {
		return client
	}
	return r.data.ent
}

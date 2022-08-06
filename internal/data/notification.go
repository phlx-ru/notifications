package data

import (
	"context"
	databaseSql "database/sql"
	"errors"
	"runtime/debug"
	"time"

	"notifications/ent"
	"notifications/ent/predicate"
	"notifications/ent/schema"
	"notifications/internal/biz"

	entSql "entgo.io/ent/dialect/sql"
	"github.com/go-kratos/kratos/v2/log"
)

type notificationRepo struct {
	data *Data
	log  *log.Helper
}

// NewNotificationRepo .
func NewNotificationRepo(data *Data, logger log.Logger) biz.NotificationRepo {
	return &notificationRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *notificationRepo) Save(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
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

	return creating.Save(ctx)
}

func (r *notificationRepo) Update(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
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

	updated, err := updating.Save(ctx)
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (r *notificationRepo) FindByID(ctx context.Context, id int64) (*ent.Notification, error) {
	return r.client(ctx).Notification.Get(ctx, int(id))
}

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

func OrderByCreatedAt() ent.OrderFunc {
	return ent.Asc(`created_at`)
}

func (r *notificationRepo) ListWaitingNotificationsWithLock(ctx context.Context, limit int) (
	[]*ent.Notification,
	error,
) {
	return r.ListWithOrderAndFilters(
		ctx,
		limit,
		OrderByCreatedAt(),
		FilterByType(schema.Types...),
		FilterByStatus(schema.StatusPending, schema.StatusRetry),
		FilterByPlannedAt(time.Now()),
		FilterForUpdateWithSkipLocked(),
	)
}

func (r *notificationRepo) ListWithOrderAndFilters(
	ctx context.Context,
	limit int,
	order ent.OrderFunc,
	filters ...predicate.Notification,
) ([]*ent.Notification, error) {
	return r.client(ctx).Notification.Query().
		Where(filters...).
		Order(order).
		Limit(limit).
		All(ctx)
}

func (r *notificationRepo) Transaction(
	ctx context.Context,
	txOptions *databaseSql.TxOptions,
	processes ...func(repoCtx context.Context) error,
) error {
	tx, err := r.data.ent.BeginTx(ctx, txOptions)
	if err != nil {
		r.log.Errorf(`failed to start tx: %v`, err)
		return err
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			r.log.Errorf(`tx panic: recovered = %v; stack = %v`, recovered, string(debug.Stack()))
			if tx != nil {
				if err := tx.Rollback(); err != nil {
					r.log.Errorf(`tx panic rollback error: %v`, err)
				}
			}
		}
	}()
	repoCtx := ent.NewContext(ctx, tx.Client())
	for _, process := range processes {
		if err := process(repoCtx); err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				r.log.Errorf(`failed to rollback tx caused of err '%s' because of: %v`, err.Error(), rollbackErr)
				return rollbackErr
			}
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		r.log.Errorf(`failed to commit tx: %v`, err)
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

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

func (r *notificationRepo) CountWaitingNotifications(ctx context.Context) (int, error) {
	return r.client(ctx).Notification.Query().
		Where(
			FilterByStatus(schema.StatusPending, schema.StatusRetry),
			FilterByType(schema.Types...),
			FilterByPlannedAt(time.Now()),
		).
		Count(ctx)
}

func (r *notificationRepo) ListWaitingNotificationsWithLock(ctx context.Context, limit int) (
	[]*ent.Notification,
	error,
) {
	return r.client(ctx).Notification.Query().
		Where(
			FilterByStatus(schema.StatusPending, schema.StatusRetry),
			FilterByType(schema.Types...),
			FilterByPlannedAt(time.Now()),
			FilterForUpdateWithSkipLocked(),
		).
		Order(OrderByCreatedAt()).
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

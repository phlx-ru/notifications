package data

import (
	"context"
	"errors"

	"notifications/ent"
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

	creating := r.data.ent.Notification.Create().
		SetSenderID(n.SenderID).
		SetType(n.Type).
		SetPayload(n.Payload).
		SetTTL(n.TTL).
		SetStatus(n.Status).
		SetPlannedAt(n.PlannedAt)

	if n.RetryAt != nil {
		creating = creating.SetRetryAt(*n.RetryAt)
	}

	if n.SentAt != nil {
		creating = creating.SetSentAt(*n.SentAt)
	}

	saved, err := creating.Save(ctx)
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (r *notificationRepo) Update(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
	if n == nil {
		return nil, errors.New("notification is empty")
	}
	updating := r.data.ent.Notification.
		UpdateOne(n).
		SetSenderID(n.SenderID).
		SetType(n.Type).
		SetPayload(n.Payload).
		SetTTL(n.TTL).
		SetStatus(n.Status).
		SetPlannedAt(n.PlannedAt)

	if n.RetryAt != nil {
		updating = updating.SetRetryAt(*n.RetryAt)
	}

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
	n, err := r.data.ent.Notification.Get(ctx, int(id))
	if err != nil {
		return nil, err
	}
	return n, nil
}

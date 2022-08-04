package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"

	"notifications/ent"
)

// NotificationRepo is a Notifications repo.
type NotificationRepo interface {
	Save(context.Context, *ent.Notification) (*ent.Notification, error)
	Update(context.Context, *ent.Notification) (*ent.Notification, error)
	FindByID(context.Context, int64) (*ent.Notification, error)
}

// NotificationUsecase is a Greeter usecase.
type NotificationUsecase struct {
	repo NotificationRepo
	log  *log.Helper
}

// NewNotificationUsecase new a Notification usecase.
func NewNotificationUsecase(repo NotificationRepo, logger log.Logger) *NotificationUsecase {
	return &NotificationUsecase{repo: repo, log: log.NewHelper(logger)}
}

// CreateNotification creates a Notification, and returns saved Notification.
func (uc *NotificationUsecase) CreateNotification(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
	uc.log.WithContext(ctx).Infof("get notification for creating with payload: %v", n.Payload.String())

	return uc.repo.Save(ctx, n)
}

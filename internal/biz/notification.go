package biz

import (
	"context"
	databaseSql "database/sql"
	"net/http"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"

	v1 "notifications/api/notification/v1"
	"notifications/ent"
	"notifications/ent/schema"
	"notifications/internal/senders"
)

const (
	RetryInterval = 5 * time.Second
)

// NotificationRepo is a Notifications repo.
type NotificationRepo interface {
	Save(context.Context, *ent.Notification) (*ent.Notification, error)

	Update(context.Context, *ent.Notification) (*ent.Notification, error)

	FindByID(context.Context, int64) (*ent.Notification, error)

	ListWaitingNotificationsWithLock(ctx context.Context, limit int) (
		[]*ent.Notification,
		error,
	)

	Transaction(
		ctx context.Context,
		txOptions *databaseSql.TxOptions,
		actions ...func(repoCtx context.Context) error,
	) error

	CountWaitingNotifications(ctx context.Context) (int, error)
}

// NotificationUsecase is a Greeter usecase.
type NotificationUsecase struct {
	repo    NotificationRepo
	log     *log.Helper
	senders *senders.Senders
}

type NotificationInDTO struct { // TODO REVIEW FOR DEPENDENCIES
	SendType v1.Type
	SenderID int64
	Payload  *schema.Payload
	TTL      int
}

type NotificationOutDTO struct {
	ID   int64
	Sent bool
}

// NewNotificationUsecase new a Notification usecase.
func NewNotificationUsecase(repo NotificationRepo, senders *senders.Senders, logger log.Logger) *NotificationUsecase {
	return &NotificationUsecase{
		repo:    repo,
		senders: senders,
		log:     log.NewHelper(logger),
	}
}

// CreateNotification creates a Notification, and returns saved Notification.
func (uc *NotificationUsecase) CreateNotification(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
	uc.log.WithContext(ctx).Infof("get notification for creating with Payload: %v", n.Payload.String())

	return uc.repo.Save(ctx, n)
}

func (uc *NotificationUsecase) CountOfWaitingForProcessNotifications(ctx context.Context) (int, error) {
	return uc.repo.CountWaitingNotifications(ctx)
}

// ProcessNotifications concurrency-safe notification processing
func (uc *NotificationUsecase) ProcessNotifications(ctx context.Context, limit int) (int64, int64, error) {
	found := int64(0)
	processed := int64(0)
	transactionOptions := &databaseSql.TxOptions{
		Isolation: databaseSql.LevelReadCommitted,
		ReadOnly:  false,
	}

	transaction := func(repoCtx context.Context) error {
		list, err := uc.repo.ListWaitingNotificationsWithLock(repoCtx, limit)
		if err != nil {
			return err
		}
		found = int64(len(list))

		for _, notification := range list {
			dto := transformNotificationModelToInDTO(notification)
			err = uc.SendNotification(ctx, dto)
			if err == nil {
				notification.Status = schema.StatusSent
				notification.SentAt = pointer.ToTime(time.Now())
				processed++
			} else {
				uc.log.Warnf(`unsuccessful attempt to send notification with id %d: %v`, notification.ID, err)
				live := notification.PlannedAt.Sub(notification.CreatedAt)
				timeToLive := time.Duration(notification.TTL) * time.Second
				if live < timeToLive {
					notification.Status = schema.StatusRetry
					notification.Retries++
					notification.PlannedAt = time.Now().Add(RetryInterval)
				} else {
					uc.log.Errorf(`failed to send notification with id %d: %v`, notification.ID, err)
					notification.Status = schema.StatusFail
				}
			}
			if _, err = uc.repo.Update(repoCtx, notification); err != nil {
				processed = 0
				return err
			}
		}

		// TODO Metrics

		return nil
	}

	return found, processed, uc.repo.Transaction(ctx, transactionOptions, transaction)
}

func (uc *NotificationUsecase) SendNotification(ctx context.Context, dto *NotificationInDTO) error {
	processors := map[v1.Type]NotificationProcessor{
		v1.Type_plain: uc.ProcessPlainNotification,
		v1.Type_email: uc.ProcessEmailNotification,
	}

	processor, ok := processors[dto.SendType]
	if !ok {
		return errors.Newf(
			http.StatusInternalServerError,
			`UNKNOWN_NOTIFICATION_TYPE`,
			`failed to send notification: unknown type '%s'`,
			dto.SendType.String(),
		)
	}

	return processor(ctx, dto.Payload)
}

func (uc *NotificationUsecase) SendNotificationAndSaveToRepo(ctx context.Context, dto *NotificationInDTO) (
	*NotificationOutDTO,
	error,
) {
	result := &NotificationOutDTO{
		ID:   0,
		Sent: false,
	}
	plannedAt := time.Now()
	err := uc.SendNotification(ctx, dto)
	if err != nil {
		return result, err
	}
	result.Sent = true
	sentAt := time.Now()
	// TODO metrics of process time

	model := transformNotificationInDTOToModel(
		dto, func(notification *ent.Notification) {
			notification.Status = schema.StatusSent
			notification.PlannedAt = plannedAt
			notification.SentAt = &sentAt
		},
	)

	notification, err := uc.repo.Save(ctx, model)
	if notification != nil && notification.ID != 0 {
		result.ID = int64(notification.ID)
	}
	return result, err
}

func (uc *NotificationUsecase) EnqueueNotification(ctx context.Context, dto *NotificationInDTO) (
	*NotificationOutDTO,
	error,
) {
	result := &NotificationOutDTO{
		ID:   0,
		Sent: false,
	}
	model := transformNotificationInDTOToModel(
		dto, func(notification *ent.Notification) {
			notification.Status = schema.StatusSent
			notification.PlannedAt = time.Now()
		},
	)
	notification, err := uc.repo.Save(ctx, model)
	if notification != nil && notification.ID != 0 {
		result.ID = int64(notification.ID)
	}
	return result, err
}

type NotificationProcessor func(context.Context, *schema.Payload) error

func (uc *NotificationUsecase) ProcessEmailNotification(_ context.Context, payload *schema.Payload) error {
	payloadEmail, err := payload.ToPayloadEmail()
	if err != nil {
		return err
	}
	if err := payloadEmail.Validate(); err != nil {
		return err
	}
	send := uc.senders.EmailSender.SendText
	if payloadEmail.IsHTML {
		send = uc.senders.EmailSender.SendHTML
	}
	return send([]string{payloadEmail.To}, payloadEmail.Subject, payloadEmail.Body)
}

func (uc *NotificationUsecase) ProcessPlainNotification(_ context.Context, payload *schema.Payload) error {
	payloadPlain, err := payload.ToPayloadPlain()
	if err != nil {
		return err
	}
	if err := payloadPlain.Validate(); err != nil {
		return err
	}
	return uc.senders.PlainSender.Send(payloadPlain.Message)
}

package biz

import (
	"context"
	databaseSql "database/sql"
	"net/http"
	"strings"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"

	v1 "notifications/api/notification/v1"
	"notifications/ent"
	"notifications/ent/schema"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"
	"notifications/internal/senders"
)

const (
	RetryInterval = 5 * time.Second

	metricCreateNotificationSuccess = `biz.notification.createNotification.success`
	metricCreateNotificationFailure = `biz.notification.createNotification.failure`
	metricCreateNotificationTimings = `biz.notification.createNotification.timings`

	metricCountOfPendingNotificationsSuccess = `biz.notification.countOfPendingNotifications.success`
	metricCountOfPendingNotificationsFailure = `biz.notification.countOfPendingNotifications.failure`
	metricCountOfPendingNotificationsTimings = `biz.notification.countOfPendingNotifications.timings`

	metricProcessNotificationsSuccess = `biz.notification.processNotifications.success`
	metricProcessNotificationsFailure = `biz.notification.processNotifications.failure`
	metricProcessNotificationsTimings = `biz.notification.processNotifications.timings`

	metricSendNotificationSuccess = `biz.notification.sendNotification.success`
	metricSendNotificationFailure = `biz.notification.sendNotification.failure`
	metricSendNotificationTimings = `biz.notification.sendNotification.timings`

	metricSendNotificationAndSaveToRepoSuccess = `biz.notification.sendNotificationAndSaveToRepo.success`
	metricSendNotificationAndSaveToRepoFailure = `biz.notification.sendNotificationAndSaveToRepo.failure`
	metricSendNotificationAndSaveToRepoTimings = `biz.notification.sendNotificationAndSaveToRepo.timings`

	metricEnqueueNotificationSuccess = `biz.notification.enqueueNotification.success`
	metricEnqueueNotificationFailure = `biz.notification.enqueueNotification.failure`
	metricEnqueueNotificationTimings = `biz.notification.enqueueNotification.timings`

	metricProcessEmailNotificationSuccess = `biz.notification.processEmailNotification.success`
	metricProcessEmailNotificationFailure = `biz.notification.processEmailNotification.failure`
	metricProcessEmailNotificationTimings = `biz.notification.processEmailNotification.timings`

	metricProcessPlainNotificationSuccess = `biz.notification.processPlainNotification.success`
	metricProcessPlainNotificationFailure = `biz.notification.processPlainNotification.failure`
	metricProcessPlainNotificationTimings = `biz.notification.processPlainNotification.timings`
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
	senders *senders.Senders
	metric  metrics.Metrics
	logs    logger.Logger
}

type NotificationInDTO struct { // TODO REVIEW FOR DEPENDENCIES
	SendType  v1.Type
	SenderID  int64
	Payload   *schema.Payload
	TTL       int
	PlannedAt *time.Time
}

type NotificationOutDTO struct {
	ID   int64
	Sent bool
}

// NewNotificationUsecase new a Notification usecase.
func NewNotificationUsecase(
	repo NotificationRepo,
	senders *senders.Senders,
	metric metrics.Metrics,
	logs log.Logger,
) *NotificationUsecase {
	return &NotificationUsecase{
		repo:    repo,
		senders: senders,
		metric:  metric,
		logs:    logger.NewHelper(logs, "ts", log.DefaultTimestamp, "scope", "biz-notification"),
	}
}

// CreateNotification creates a Notification, and returns saved Notification.
func (uc *NotificationUsecase) CreateNotification(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
	defer uc.metric.NewTiming().Send(metricCreateNotificationTimings)
	notification, err := uc.repo.Save(ctx, n)
	if err != nil {
		uc.metric.Increment(metricCreateNotificationFailure)
		uc.logs.WithContext(ctx).Errorf("failed to create notification with payload [%s]: %v", n.Payload.String(), err)
	} else {
		uc.metric.Increment(metricCreateNotificationSuccess)
		uc.logs.WithContext(ctx).Infof("successfully created notification with payload [%s]", n.Payload.String())
	}
	return notification, err
}

func (uc *NotificationUsecase) CountOfPendingNotifications(ctx context.Context) (int, error) {
	defer uc.metric.NewTiming().Send(metricCountOfPendingNotificationsTimings)
	cnt, err := uc.repo.CountWaitingNotifications(ctx)
	if err != nil {
		uc.metric.Increment(metricCountOfPendingNotificationsFailure)
		uc.logs.WithContext(ctx).Errorf("failed to count of pending notifications: %v", err)
	} else {
		uc.metric.Increment(metricCountOfPendingNotificationsSuccess)
		uc.logs.WithContext(ctx).Infof("successfully got count of pending notification: %d", cnt)
	}
	return cnt, err
}

// ProcessNotifications concurrency-safe notification processing
func (uc *NotificationUsecase) ProcessNotifications(ctx context.Context, limit int) (int64, int64, error) {
	defer uc.metric.NewTiming().Send(metricProcessNotificationsTimings)
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
				uc.logs.WithContext(ctx).Warnf(
					`unsuccessful attempt to send notification with id %d: %v`,
					notification.ID,
					err,
				)

				notification.Status = schema.StatusRetry
				notification.Retries++
				notification.RetryAt = pointer.ToTime(time.Now().Add(RetryInterval))

				timeFrom := notification.PlannedAt
				if notification.RetryAt != nil {
					timeFrom = *notification.RetryAt
				}

				live := timeFrom.Sub(notification.PlannedAt)
				timeToLive := time.Duration(notification.TTL) * time.Second

				if live > timeToLive {
					uc.logs.WithContext(ctx).Errorf(`failed to send notification with id %d: %v`, notification.ID, err)
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

	err := uc.repo.Transaction(ctx, transactionOptions, transaction)
	if err != nil {
		uc.metric.Increment(metricProcessNotificationsFailure)
		uc.logs.WithContext(ctx).Errorf("failed to process notifications: %v", err)
	} else {
		uc.metric.Increment(metricProcessNotificationsSuccess)
		uc.logs.WithContext(ctx).Infof("successfully processed notifications: found %d, processed %d", found, processed)
	}
	return found, processed, err
}

func (uc *NotificationUsecase) SendNotification(ctx context.Context, dto *NotificationInDTO) error {
	defer uc.metric.NewTiming().Send(metricSendNotificationTimings)

	processors := map[v1.Type]NotificationProcessor{
		v1.Type_plain: uc.ProcessPlainNotification,
		v1.Type_email: uc.ProcessEmailNotification,
	}

	var err error
	processor, ok := processors[dto.SendType]

	if !ok {
		err = errors.Newf(
			http.StatusInternalServerError,
			`UNKNOWN_NOTIFICATION_TYPE`,
			`failed to send notification: unknown type '%s'`,
			dto.SendType.String(),
		)
	} else {
		err = processor(ctx, dto.Payload)
	}

	if err != nil {
		uc.metric.Increment(metricSendNotificationFailure)
		uc.logs.WithContext(ctx).Errorf("failed to send notification: %v", err)
	} else {
		uc.metric.Increment(metricSendNotificationSuccess)
		uc.logs.WithContext(ctx).Info("successfully sent notification")
	}

	return err
}

func (uc *NotificationUsecase) SendNotificationAndSaveToRepo(ctx context.Context, dto *NotificationInDTO) (
	*NotificationOutDTO,
	error,
) {
	defer uc.metric.NewTiming().Send(metricSendNotificationAndSaveToRepoTimings)
	var err error
	result := &NotificationOutDTO{
		ID:   0,
		Sent: false,
	}
	plannedAt := time.Now()
	err = uc.SendNotification(ctx, dto)
	if err == nil {
		result.Sent = true

		model := transformNotificationInDTOToModel(
			dto, func(notification *ent.Notification) {
				notification.Status = schema.StatusSent
				notification.PlannedAt = plannedAt
				notification.SentAt = pointer.ToTime(time.Now())
			},
		)

		var notification *ent.Notification
		notification, err = uc.repo.Save(ctx, model)
		if notification != nil && notification.ID != 0 {
			result.ID = int64(notification.ID)
		}
	}

	if err != nil {
		uc.metric.Increment(metricSendNotificationAndSaveToRepoFailure)
		uc.logs.WithContext(ctx).Errorf("failed to send notification and save to repo: %v", err)
	} else {
		uc.metric.Increment(metricSendNotificationAndSaveToRepoSuccess)
		uc.logs.WithContext(ctx).Info("successfully sent notification and saved to repo")
	}

	return result, err
}

func (uc *NotificationUsecase) EnqueueNotification(ctx context.Context, dto *NotificationInDTO) (
	*NotificationOutDTO,
	error,
) {
	defer uc.metric.NewTiming().Send(metricEnqueueNotificationTimings)

	result := &NotificationOutDTO{
		ID:   0,
		Sent: false,
	}
	model := transformNotificationInDTOToModel(
		dto, func(notification *ent.Notification) {
			notification.Status = schema.StatusPending
		},
	)
	notification, err := uc.repo.Save(ctx, model)
	if notification != nil && notification.ID != 0 {
		result.ID = int64(notification.ID)
	}
	if err != nil {
		uc.metric.Increment(metricEnqueueNotificationFailure)
		uc.logs.WithContext(ctx).Errorf("failed to enqueue notification: %v", err)
	} else {
		uc.metric.Increment(metricEnqueueNotificationSuccess)
		uc.logs.WithContext(ctx).Infof("successfully enqueue notification with id %d", result.ID)
	}
	return result, err
}

type NotificationProcessor func(context.Context, *schema.Payload) error

func (uc *NotificationUsecase) ProcessEmailNotification(ctx context.Context, payload *schema.Payload) error {
	defer uc.metric.NewTiming().Send(metricProcessEmailNotificationTimings)
	var err error
	defer func() {
		if err != nil {
			uc.metric.Increment(metricProcessEmailNotificationFailure)
			uc.logs.WithContext(ctx).Errorf("failed to process email notification: %v", err)
		} else {
			uc.metric.Increment(metricProcessEmailNotificationSuccess)
			uc.logs.WithContext(ctx).Info("successfully process email notification")
		}
	}()
	payloadEmail, err := payload.ToPayloadEmail()
	if err != nil {
		return err
	}
	if err = payloadEmail.Validate(); err != nil {
		return err
	}
	send := uc.senders.EmailSender.SendText
	if isTrue(payloadEmail.IsHTML) {
		send = uc.senders.EmailSender.SendHTML
	}
	err = send(ctx, []string{payloadEmail.To}, payloadEmail.Subject, payloadEmail.Body)
	return err
}

func (uc *NotificationUsecase) ProcessPlainNotification(ctx context.Context, payload *schema.Payload) error {
	defer uc.metric.NewTiming().Send(metricProcessPlainNotificationTimings)
	var err error
	defer func() {
		if err != nil {
			uc.metric.Increment(metricProcessPlainNotificationFailure)
			uc.logs.WithContext(ctx).Errorf("failed to process plain notification: %v", err)
		} else {
			uc.metric.Increment(metricProcessPlainNotificationSuccess)
			uc.logs.WithContext(ctx).Info("successfully process plain notification")
		}
	}()
	var payloadPlain *schema.PayloadPlain
	payloadPlain, err = payload.ToPayloadPlain()
	if err != nil {
		return err
	}
	if err = payloadPlain.Validate(); err != nil {
		return err
	}
	err = uc.senders.PlainSender.Send(ctx, payloadPlain.Message)
	return err
}

func isTrue(bool string) bool {
	lowed := strings.ToLower(bool)
	variants := []string{"1", "true", "yes", "y", "t"}
	for _, variant := range variants {
		if variant == lowed {
			return true
		}
	}
	return false
}

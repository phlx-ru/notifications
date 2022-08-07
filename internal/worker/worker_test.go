package worker

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	"notifications/ent"
	"notifications/ent/schema"
	"notifications/internal/biz"
	"notifications/internal/senders"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

//go:generate mockgen -source ./${GOFILE} -destination ./notification_repo_mock_test.go -package ${GOPACKAGE}

type NotificationRepo interface {
	biz.NotificationRepo
}

type EmailSender interface {
	senders.EmailSender
}

type PlainSender interface {
	senders.PlainSender
}

func TestWorker_Run(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	loggerInstance := log.With(log.DefaultLogger, "ts", log.DefaultTimestamp)
	logger := log.NewFilter(loggerInstance, log.FilterLevel(log.LevelFatal))

	transaction := func(ctx context.Context, _ *sql.TxOptions, actions ...func(context.Context) error) error {
		for _, action := range actions {
			if err := action(ctx); err != nil {
				return err
			}
		}
		return nil
	}

	update := func(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
		rand.Seed(time.Now().Unix())
		n.ID = rand.Int()
		return n, nil
	}

	testCases := []struct {
		name             string
		notificationRepo func() NotificationRepo
		plainSender      func() PlainSender
		emailSender      func() EmailSender
		expected         error
	}{
		{
			name: "basic",
			notificationRepo: func() NotificationRepo {
				notificationRepoMock := NewMockNotificationRepo(ctrl)
				notificationRepoMock.EXPECT().CountWaitingNotifications(gomock.Any()).Return(0, nil)
				notificationRepoMock.EXPECT().Transaction(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				notificationRepoMock.EXPECT().Update(gomock.Any(), gomock.Any()).Times(0)
				notificationRepoMock.EXPECT().ListWaitingNotificationsWithLock(gomock.Any(), gomock.Any()).Times(0)
				return notificationRepoMock
			},
			plainSender: func() PlainSender {
				plainSender := NewMockPlainSender(ctrl)
				plainSender.EXPECT().Send(gomock.Any()).Times(0)
				return plainSender
			},
			emailSender: func() EmailSender {
				emailSender := NewMockEmailSender(ctrl)
				emailSender.EXPECT().SendText(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				emailSender.EXPECT().SendHTML(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				return emailSender
			},
		},
		{
			name: "a-thousand",
			notificationRepo: func() NotificationRepo {
				notificationRepoMock := NewMockNotificationRepo(ctrl)
				notificationRepoMock.EXPECT().
					CountWaitingNotifications(gomock.Any()).
					Return(1000, nil).
					Times(1)

				notificationRepoMock.EXPECT().
					Transaction(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(transaction).
					Times(10)

				notificationRepoMock.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					DoAndReturn(update).
					Times(100)

				notificationRepoMock.EXPECT().
					ListWaitingNotificationsWithLock(gomock.Any(), gomock.Any()).
					DoAndReturn(
						func(_ context.Context, _ int) ([]*ent.Notification, error) {
							return makePlainNotifications(10, "test message")
						},
					).
					Times(10)
				return notificationRepoMock
			},
			plainSender: func() PlainSender {
				plainSender := NewMockPlainSender(ctrl)
				plainSender.EXPECT().Send(gomock.Any()).Times(100)
				return plainSender
			},
			emailSender: func() EmailSender {
				emailSender := NewMockEmailSender(ctrl)
				emailSender.EXPECT().SendText(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				emailSender.EXPECT().SendHTML(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				return emailSender
			},
		},
		{
			name: "retries",
			notificationRepo: func() NotificationRepo {
				notificationRepoMock := NewMockNotificationRepo(ctrl)
				notificationRepoMock.EXPECT().
					CountWaitingNotifications(gomock.Any()).
					Return(10, nil).
					Times(1)

				notificationRepoMock.EXPECT().
					Transaction(gomock.Any(), gomock.Any(), gomock.Any()).
					DoAndReturn(transaction).
					Times(1)

				notificationRepoMock.EXPECT().
					Update(gomock.Any(), gomock.Any()).
					DoAndReturn(
						func(ctx context.Context, n *ent.Notification) (*ent.Notification, error) {
							require.Equal(t, 1, n.Retries)
							require.Equal(t, schema.StatusRetry, n.Status)

							plannedAtDuration := n.PlannedAt.Sub(time.Now())
							// Milliseconds subtract because of parallel run of tests
							moreThanRetryInterval := plannedAtDuration > biz.RetryInterval-50*time.Millisecond
							lessThanRetryIntervalWithFewSeconds := plannedAtDuration < biz.RetryInterval+3*time.Second
							isBetween := moreThanRetryInterval && lessThanRetryIntervalWithFewSeconds

							require.True(t, isBetween)

							return n, nil
						},
					).
					Times(10)

				notificationRepoMock.EXPECT().
					ListWaitingNotificationsWithLock(gomock.Any(), gomock.Any()).
					DoAndReturn(
						func(_ context.Context, _ int) ([]*ent.Notification, error) {
							return makePlainNotifications(10, "test message")
						},
					).
					Times(1)
				return notificationRepoMock
			},
			plainSender: func() PlainSender {
				plainSender := NewMockPlainSender(ctrl)
				plainSender.EXPECT().Send(gomock.Any()).Return(errors.New("test for failed send")).Times(10)
				return plainSender
			},
			emailSender: func() EmailSender {
				emailSender := NewMockEmailSender(ctrl)
				emailSender.EXPECT().SendText(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				emailSender.EXPECT().SendHTML(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
				return emailSender
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				notificationsRepoMock := testCase.notificationRepo()

				sendersMock := &senders.Senders{
					PlainSender: testCase.plainSender(),
					EmailSender: testCase.emailSender(),
				}

				usecase := biz.NewNotificationUsecase(notificationsRepoMock, sendersMock, logger)

				worker := New(usecase, logger, RunOnceOption())

				err := worker.Run(ctx)
				require.Nil(t, err)
			},
		)
	}
}

func makePlainNotifications(length int, message string) ([]*ent.Notification, error) {
	payload, err := schema.PayloadFromProto(map[string]string{"message": message})
	if err != nil {
		return nil, err
	}

	now := time.Now()
	notifications := []*ent.Notification{}
	for i := 1; i <= length; i++ {
		n := &ent.Notification{
			ID:        8080000 + i,
			Type:      schema.TypePlain,
			Payload:   *payload,
			TTL:       100,
			Status:    schema.StatusPending,
			CreatedAt: now,
			UpdatedAt: now,
			PlannedAt: now,
		}
		notifications = append(notifications, n)
	}

	return notifications, nil
}

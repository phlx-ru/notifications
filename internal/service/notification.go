package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"notifications/ent"
	"notifications/ent/schema"
	"notifications/internal/biz"
	"notifications/internal/senders"

	v1 "notifications/api/notification/v1"

	"github.com/go-kratos/kratos/v2/log"
)

type NotificationService struct {
	v1.UnimplementedNotificationServer

	usecase *biz.NotificationUsecase

	senders *senders.Senders

	logger *log.Helper
}

func NewNotificationService(u *biz.NotificationUsecase, s *senders.Senders, l log.Logger) *NotificationService {
	return &NotificationService{
		usecase: u,
		senders: s,
		logger:  log.NewHelper(l),
	}
}

func (s *NotificationService) CreatingTest(ctx context.Context, req *v1.CreatingTestRequest) (
	*v1.CreatingTestReply,
	error,
) {
	n := &ent.Notification{
		SenderID: 0,
		Type:     schema.TypeEmail.String(),
		Payload: (&schema.PayloadEmail{
			To:      `phlx@ya.ru`,
			Subject: `Test email notification`,
			Body:    fmt.Sprintf(`%s — your personal auth code, keep it simple. Message: %s`, code(), req.Message),
		}).MustToPayload(),
		TTL:       300,
		Status:    schema.StatusDraft.String(),
		PlannedAt: time.Now(),
	}

	res, err := s.usecase.CreateNotification(ctx, n)
	if err != nil {
		return nil, err
	}
	return &v1.CreatingTestReply{Result: fmt.Sprintf(`ok with id %d`, res.ID)}, nil
}

func (s *NotificationService) Enqueue(ctx context.Context, req *v1.SendRequest) (*v1.EnqueueResponse, error) {
	return nil, nil
}

func (s *NotificationService) Send(ctx context.Context, req *v1.SendRequest) (*v1.SendResponse, error) {
	payload, err := schema.PayloadFromProto(req.Payload)
	if err != nil {
		return nil, err
	}

	id, err := s.process(
		ctx,
		&processData{
			sendType: req.Type,
			senderID: 0, // TODO
			payload:  payload,
			ttl:      int(req.Ttl),
		},
	)

	s.logger.Infof("notification was sent and has id %d", id)

	return &v1.SendResponse{}, nil
}

func code() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprintf("%d", rand.Int()*10000%10000)
}

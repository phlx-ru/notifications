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

// CreatingTest TODO REMOVE
func (s *NotificationService) CreatingTest(ctx context.Context, req *v1.CreatingTestRequest) (
	*v1.CreatingTestReply,
	error,
) {
	code := func() string {
		rand.Seed(time.Now().Unix())
		return fmt.Sprintf("%d", rand.Int()*10000%10000) //nolint
	}
	n := &ent.Notification{
		SenderID: 0,
		Type:     schema.TypeEmail,
		Payload: (&schema.PayloadEmail{
			To:      `phlx@ya.ru`,
			Subject: `Test email notification`,
			Body:    fmt.Sprintf(`%s — your personal auth code, keep it simple. Message: %s`, code(), req.Message),
		}).MustToPayload(),
		TTL:       300,
		Status:    schema.StatusDraft,
		PlannedAt: time.Now(),
	}

	res, err := s.usecase.CreateNotification(ctx, n)
	if err != nil {
		return nil, err
	}
	return &v1.CreatingTestReply{Result: fmt.Sprintf(`ok with id %d`, res.ID)}, nil
}

func (s *NotificationService) Enqueue(ctx context.Context, req *v1.SendRequest) (*v1.EnqueueResponse, error) {
	payload, err := schema.PayloadFromProto(req.Payload)
	if err != nil {
		return nil, err
	}

	in := &biz.NotificationInDTO{
		SendType: req.Type,
		SenderID: req.SenderId,
		Payload:  payload,
		TTL:      int(req.Ttl),
	}

	response := &v1.EnqueueResponse{}
	result, err := s.usecase.EnqueueNotification(ctx, in)
	if result != nil {
		response.Id = result.ID
	}
	if err == nil {
		s.logger.Infof("notification %d was sent successfully", response.Id)
	} else {
		s.logger.Errorf("notification %d was failed to send: %v", response.Id, err)
	}

	return response, err
}

func (s *NotificationService) Send(ctx context.Context, req *v1.SendRequest) (*v1.SendResponse, error) {
	payload, err := schema.PayloadFromProto(req.Payload)
	if err != nil {
		return nil, err
	}

	in := &biz.NotificationInDTO{
		SendType: req.Type,
		SenderID: req.SenderId,
		Payload:  payload,
		TTL:      int(req.Ttl),
	}

	response := &v1.SendResponse{}
	result, err := s.usecase.SendNotificationAndSaveToRepo(ctx, in)
	if result != nil {
		response.Id = result.ID
		response.Sent = result.Sent
	}
	if err == nil {
		s.logger.Infof("notification %d was sent successfully", response.Id)
	} else {
		s.logger.Errorf("notification %d was failed to send: %v", response.Id, err)
	}

	return response, err
}

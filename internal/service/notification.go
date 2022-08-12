package service

import (
	"context"

	"notifications/ent/schema"
	"notifications/internal/biz"
	"notifications/internal/senders"

	v1 "notifications/api/notification/v1"

	"github.com/AlekSi/pointer"
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

	if req.PlannedAt != nil {
		in.PlannedAt = pointer.ToTime(req.PlannedAt.AsTime())
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

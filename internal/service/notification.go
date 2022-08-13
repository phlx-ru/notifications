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

var (
	StatusesSchemaToProtoMap = map[schema.NotificationStatus]v1.Status{
		schema.StatusDraft:   v1.Status_draft,
		schema.StatusPending: v1.Status_pending,
		schema.StatusSent:    v1.Status_sent,
		schema.StatusRetry:   v1.Status_retry,
		schema.StatusFail:    v1.Status_fail,
	}

	TypesProtoToSchemaMap = map[v1.Type]schema.NotificationType{
		v1.Type_plain:    schema.TypePlain,
		v1.Type_email:    schema.TypeEmail,
		v1.Type_sms:      schema.TypeSMS,
		v1.Type_push:     schema.TypePush,
		v1.Type_telegram: schema.TypeTelegram,
		v1.Type_whatsapp: schema.TypeWhatsApp,
	}
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

func (s *NotificationService) Check(ctx context.Context, req *v1.CheckRequest) (*v1.CheckResponse, error) {
	if req.Id < 0 {
		return nil, v1.ErrorInvalidRequest(`validation failed: id=%d is incorrect`, req.Id)
	}
	if req.Id == 0 {
		return nil, v1.ErrorInvalidRequest(`validation failed: id was not set`)
	}
	status, err := s.usecase.CheckStatus(ctx, req.Id)
	if err == biz.ErrNotificationNotFound {
		return nil, v1.ErrorNotificationNotFound(`notification with id %d was not found`, req.Id)
	}
	if err != nil {
		return nil, v1.ErrorInternalError(`check status failed: %v`, err)
	}

	return &v1.CheckResponse{
		Status: StatusesSchemaToProtoMap[*status],
	}, nil
}

func (s *NotificationService) Enqueue(ctx context.Context, req *v1.SendRequest) (*v1.EnqueueResponse, error) {
	payload, err := schema.PayloadFromProto(req.Payload)
	if err != nil {
		return nil, v1.ErrorInternalError(`payload conversion failed: %v`, err)
	}

	notificationType, ok := TypesProtoToSchemaMap[req.Type]
	if !ok {
		return nil, v1.ErrorInvalidRequest(`validation failed: type %s is unknown`, req.Type.String())
	}
	err = payload.Validate(notificationType)
	if err != nil {
		return nil, v1.ErrorInvalidRequest(`validation failed: %v`, err)
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
	if err != nil {
		s.logger.Errorf(`notification %d was failed to send: %v`, response.Id, err)
		return nil, v1.ErrorInternalError(`enqueue notification failed: %v`, err)
	}
	s.logger.Infof("notification %d was sent successfully", response.Id)
	return response, nil
}

func (s *NotificationService) Send(ctx context.Context, req *v1.SendRequest) (*v1.SendResponse, error) {
	payload, err := schema.PayloadFromProto(req.Payload)
	if err != nil {
		return nil, v1.ErrorInternalError(`payload conversion failed: %v`, err)
	}

	notificationType, ok := TypesProtoToSchemaMap[req.Type]
	if !ok {
		return nil, v1.ErrorInvalidRequest(`validation failed: type %s is unknown`, req.Type.String())
	}
	err = payload.Validate(notificationType)
	if err != nil {
		return nil, v1.ErrorInvalidRequest(`validation failed: %v`, err)
	}

	in := &biz.NotificationInDTO{
		SendType: req.Type,
		SenderID: req.SenderId,
		Payload:  payload,
		TTL:      int(req.Ttl),
	}

	result, err := s.usecase.SendNotification(ctx, in)
	if err != nil {
		s.logger.Errorf(`notification was failed to send: %v`, err)
		return nil, v1.ErrorInternalError(`send notification failed: %v`, err)
	}
	s.logger.Infof("notification %d was sent successfully", result.ID)
	return &v1.SendResponse{
		Id:   result.ID,
		Sent: result.Sent,
	}, nil
}

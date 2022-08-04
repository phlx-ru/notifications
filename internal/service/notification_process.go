package service

import (
	"context"
	"net/http"
	"time"

	"github.com/go-kratos/kratos/v2/errors"

	v1 "notifications/api/notification/v1"
	"notifications/ent"
	"notifications/ent/schema"
)

type processor func(ctx context.Context, pd *processData) (int64, error)

type processData struct {
	sendType v1.SendRequest_Type
	senderID int64
	payload  *schema.Payload
	ttl      int
}

func (s *NotificationService) process(ctx context.Context, pd *processData) (int64, error) {
	processors := map[v1.SendRequest_Type]processor{
		v1.SendRequest_email: s.processEmail,
	}

	processor, ok := processors[pd.sendType]
	if !ok {
		return 0, errors.Newf(
			http.StatusInternalServerError,
			`NOTIFICATION_PROCESSOR_UNKNOWN`,
			`failed to send notification: unknown type '%s'`,
			pd.sendType.String(),
		)
	}

	return processor(ctx, pd)
}

func (s *NotificationService) createSentNotification(ctx context.Context, pd *processData) (*ent.Notification, error) {
	now := time.Now()
	return s.usecase.CreateNotification(
		ctx, &ent.Notification{
			SenderID:  int(pd.senderID),
			Type:      pd.sendType.String(),
			Payload:   *pd.payload,
			TTL:       pd.ttl,
			Status:    schema.StatusSent.String(),
			PlannedAt: now,
			SentAt:    &now,
		},
	)
}

func (s *NotificationService) processEmail(ctx context.Context, pd *processData) (int64, error) {
	pe, err := pd.payload.ToPayloadEmail()
	if err != nil {
		return 0, err
	}
	if err := pe.Validate(); err != nil {
		return 0, err
	}
	if pe.IsHTML {
		err = s.senders.EmailSender.SendHTML([]string{pe.To}, pe.Subject, pe.Body)
	} else {
		err = s.senders.EmailSender.SendText([]string{pe.To}, pe.Subject, pe.Body)
	}
	if err != nil {
		return 0, err
	}
	n, err := s.createSentNotification(ctx, pd)
	if err != nil {
		return 0, err
	}
	return int64(n.ID), nil
}

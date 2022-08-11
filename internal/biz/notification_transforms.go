package biz

import (
	v1 "notifications/api/notification/v1"
	"notifications/ent"
	"notifications/ent/schema"
)

func transformNotificationModelToInDTO(notification *ent.Notification) *NotificationInDTO {
	return &NotificationInDTO{
		SendType:  v1.Type(v1.Type_value[notification.Type.String()]),
		SenderID:  int64(notification.SenderID),
		Payload:   &notification.Payload,
		TTL:       notification.TTL,
		PlannedAt: &notification.PlannedAt,
	}
}

func transformNotificationInDTOToModel(
	dto *NotificationInDTO,
	withFields ...func(*ent.Notification),
) *ent.Notification {
	notification := &ent.Notification{
		SenderID: int(dto.SenderID),
		Type:     schema.NotificationType(dto.SendType.String()),
		Payload:  *dto.Payload,
		TTL:      dto.TTL,
	}
	if dto.PlannedAt != nil {
		notification.PlannedAt = *dto.PlannedAt
	}
	for _, withField := range withFields {
		withField(notification)
	}
	return notification
}

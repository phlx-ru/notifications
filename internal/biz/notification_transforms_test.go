package biz

import (
	"testing"
	"time"

	v1 "notifications/api/notification/v1"
	"notifications/ent"
	"notifications/ent/schema"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"
)

func Test_transformNotificationInDTOToModel(t *testing.T) {
	now := time.Now()

	payload := schema.Payload(map[string]string{`test`: `test`})

	testCases := []struct {
		name      string
		dto       *NotificationInDTO
		withField func(*ent.Notification)
		expected  *ent.Notification
	}{
		{
			name: "basic",
			dto: &NotificationInDTO{
				SendType: v1.Type(v1.Type_value[schema.TypeEmail.String()]),
				Payload:  &payload,
			},
			withField: func(notification *ent.Notification) {
				notification.Status = schema.StatusSent
				notification.PlannedAt = now
				notification.SentAt = pointer.ToTime(now.Add(1 * time.Minute))
			},
			expected: &ent.Notification{
				Type:      schema.TypeEmail,
				Status:    schema.StatusSent,
				Payload:   payload,
				PlannedAt: now,
				SentAt:    pointer.ToTime(now.Add(1 * time.Minute)),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(
			testCase.name, func(t *testing.T) {
				actual := transformNotificationInDTOToModel(testCase.dto, testCase.withField)
				require.Equal(t, testCase.expected, actual)
			},
		)
	}
}

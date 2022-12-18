package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1 "notifications/api/notification/v1"
	"notifications/ent"
	"notifications/ent/schema"

	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"
)

type AbstractJSON map[string]any

func TestV1Check(t *testing.T) {
	var (
		err          error
		notification *ent.Notification
		expect       *httpexpect.Expect
		server       *httptest.Server
	)

	t.Run(
		`prerequisites`, func(t *testing.T) {
			ctx := context.Background()

			now := time.Now()

			model := &ent.Notification{
				Type: schema.TypePlain,
				Payload: (schema.PayloadPlain{
					Message: `hello from api tests!`,
				}).MustToPayload(),
				TTL:       600,
				Status:    schema.StatusSent,
				PlannedAt: now,
				SentAt:    &now,
			}

			notification, err = notificationRepo.Create(ctx, model)
			require.NoError(t, err)

			server = httptest.NewServer(httpServer)

			expect = httpexpect.New(t, server.URL)
		},
	)

	defer server.Close()

	t.Run(
		`existed_id_passed`, func(t *testing.T) {
			expect.POST(`/v1/check`).
				WithHeader(`Authorization`, `Bearer `+jwtToken).
				WithJSON(AbstractJSON{`id`: notification.ID}).
				Expect().
				Status(http.StatusOK).
				JSON().
				Equal(AbstractJSON{`status`: schema.StatusSent})
		},
	)

	t.Run(
		`not_existed_id_passed`, func(t *testing.T) {
			expect.POST(`/v1/check`).
				WithHeader(`Authorization`, `Bearer `+jwtToken).
				WithJSON(AbstractJSON{`id`: notification.ID + 1}).
				Expect().
				Status(http.StatusNotFound).
				JSON().
				Object().
				ContainsMap(
					AbstractJSON{
						`code`:   http.StatusNotFound,
						`reason`: v1.ErrorReason_NOTIFICATION_NOT_FOUND.String(),
					},
				)
		},
	)

	t.Run(
		`negative_id_passed`, func(t *testing.T) {
			expect.POST(`/v1/check`).
				WithHeader(`Authorization`, `Bearer `+jwtToken).
				WithJSON(AbstractJSON{`id`: -4}).
				Expect().
				Status(http.StatusBadRequest).
				JSON().
				Object().
				ContainsMap(
					AbstractJSON{
						`code`:   http.StatusBadRequest,
						`reason`: v1.ErrorReason_INVALID_REQUEST.String(),
					},
				)
		},
	)

	t.Run(
		`id_not_passed`, func(t *testing.T) {
			expect.POST(`/v1/check`).
				WithHeader(`Authorization`, `Bearer `+jwtToken).
				WithJSON(AbstractJSON{`unknown`: `nevermind`}).
				Expect().
				Status(http.StatusBadRequest).
				JSON().
				Object().
				ContainsMap(
					AbstractJSON{
						`code`:   http.StatusBadRequest,
						`reason`: v1.ErrorReason_INVALID_REQUEST.String(),
					},
				)
		},
	)

	t.Run(
		`invalid_id_passed`, func(t *testing.T) {
			expect.POST(`/v1/check`).
				WithHeader(`Authorization`, `Bearer `+jwtToken).
				WithJSON(AbstractJSON{`id`: `incorrect`}).
				Expect().
				Status(http.StatusBadRequest).
				JSON().
				Object().
				ContainsMap(
					AbstractJSON{
						`code`:   http.StatusBadRequest,
						`reason`: `CODEC`,
					},
				)
		},
	)
}

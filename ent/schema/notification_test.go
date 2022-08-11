package schema

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPayload(t *testing.T) {
	payloadEmail := PayloadEmail{
		To:      "somebody@wastold.me",
		Subject: "The woooorld",
		Body:    "Is gonna roll me",
	}

	payloadEmailBytes, err := json.Marshal(payloadEmail)
	require.NoError(t, err)

	payload := payloadEmail.MustToPayload()
	require.JSONEq(t, string(payloadEmailBytes), payload.String())

	payloadEmailCasted, err := payload.ToPayloadEmail()
	require.NoError(t, err)

	require.EqualValues(t, payloadEmail, *payloadEmailCasted)

	ruinedPayload := Payload(map[string]string{"existential": "none"})

	ruinedPayloadEmail, err := ruinedPayload.ToPayloadEmail()
	require.NoError(t, err)
	require.NotNil(t, ruinedPayloadEmail)
	require.Error(t, ruinedPayloadEmail.Validate())
}

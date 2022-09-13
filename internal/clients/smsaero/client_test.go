package smsaero

import (
	"testing"

	"github.com/stretchr/testify/require"
	"notifications/internal/pkg/template"
)

func TestInterpolate(t *testing.T) {
	actual := template.MustInterpolate(
		URLTemplate, map[string]any{
			"email":    "user@company.example",
			"apiKey":   "abcdef",
			"host":     host,
			"sendPath": sendPath,
			"number":   "79009009090",
			"text":     "Тестовое сообщение",
			"sign":     sign,
		},
	)

	expected := `https://user@company.example:abcdef@gate.smsaero.ru/v2/sms/send?number=79009009090&text=%D0%A2%D0%B5%D1%81%D1%82%D0%BE%D0%B2%D0%BE%D0%B5+%D1%81%D0%BE%D0%BE%D0%B1%D1%89%D0%B5%D0%BD%D0%B8%D0%B5&sign=SMS+Aero`

	require.Equal(t, expected, actual)
}

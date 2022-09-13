package sms

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsCompatibleWithGSM7(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		expected bool
	}{
		{
			name:     `latin`,
			message:  `Hello, world!`,
			expected: true,
		},
		{
			name:     `cyrillic`,
			message:  `Привет, мир!`,
			expected: false,
		},
		{
			name:     `latin-with-cyrillic`,
			message:  `What about Ъ symbol?`,
			expected: false,
		},
		{
			name:     `special-compatible`,
			message:  "@ £ $ ¥ ¤ § Æ Ξ \n('\"\"') # ¿Π0ΓΛ0ΨΞHΔ?",
			expected: true,
		},
	}
	for _, testCase := range testCases {
		t.Run(
			testCase.name,
			func(t *testing.T) {
				actual := IsCompatibleWithGSM7(testCase.message)
				require.Equal(t, testCase.expected, actual)
			},
		)
	}
}

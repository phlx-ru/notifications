package strings

import "notifications/internal/pkg/slices"

var (
	TrueValues  = []string{"true", "t", "1", "yes", "y"}
	FalseValues = []string{"false", "f", "0", "no", "n"}
)

func IsBool(s string) bool {
	return IsTrue(s) || IsFalse(s)
}

func IsTrue(s string) bool {
	return slices.Includes(s, TrueValues)
}

func IsFalse(s string) bool {
	return slices.Includes(s, FalseValues)
}

package sms

// see https://www.twilio.com/docs/glossary/what-is-gsm-7-character-encoding

const (
	limitGSM7    = 160
	limitUCS2    = 70
	LimitOverall = 918 // LimitOverall calculated from min(6 * 153, 14 * 67)
)

var (
	GSM7 = []rune{
		'@', '£', '$', '¥', 'è', 'é', 'ù', 'ì', 'ò', 'Ç', '\n', 'Ø', 'ø', '\r', 'Å', 'å',
		'Δ', '_', 'Φ', 'Γ', 'Λ', 'Ω', 'Π', 'Ψ', 'Σ', 'Θ', 'Ξ', '\x1B', 'Æ', 'æ', 'ß', 'É',
		' ', '!', '"', '#', '¤', '%', '&', '\'', '(', ')', '*', '+', ',', '-', '.', '/',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', ':', ';', '<', '=', '>', '?',
		'¡', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O',
		'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', 'Ä', 'Ö', 'Ñ', 'Ü', '§',
		'¿', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
		'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'ä', 'ö', 'ñ', 'ü', 'à',
	}
	hashmapGSM7 = map[rune]bool{}
)

// IsCompatibleWithGSM7 returns true if message possible encode with GSM-7
func IsCompatibleWithGSM7(message string) bool {
	if len(hashmapGSM7) == 0 {
		for _, r := range GSM7 {
			hashmapGSM7[r] = true
		}
	}
	for _, r := range message {
		includes := hashmapGSM7[r]
		if !includes {
			return false
		}
	}
	return true
}

// IsExceedsLimit returns true if message exceeds limit for GSM-7 or UCS-2 encodings
func IsExceedsLimit(message string) bool {
	runes := []rune(message)
	if IsCompatibleWithGSM7(message) {
		return len(runes) > limitGSM7
	}
	return len(runes) > limitUCS2
}

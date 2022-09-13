package senders

type Senders struct {
	PlainSender    PlainSender
	EmailSender    EmailSender
	TelegramSender TelegramSender
	SMSAeroSender  SMSAeroSender
}

func NewSenders(
	plainSender PlainSender,
	emailSender EmailSender,
	telegramSender TelegramSender,
	smsAeroSender SMSAeroSender,
) *Senders {
	return &Senders{
		EmailSender:    emailSender,
		PlainSender:    plainSender,
		TelegramSender: telegramSender,
		SMSAeroSender:  smsAeroSender,
	}
}

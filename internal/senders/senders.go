package senders

type Senders struct {
	EmailSender    EmailSender
	PlainSender    PlainSender
	TelegramSender TelegramSender
}

func NewSenders(plainSender *Plain, emailSender *Email, telegramSender *Telegram) *Senders {
	return &Senders{
		EmailSender:    emailSender,
		PlainSender:    plainSender,
		TelegramSender: telegramSender,
	}
}

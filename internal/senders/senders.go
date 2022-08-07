package senders

type Senders struct {
	EmailSender EmailSender
	PlainSender PlainSender
}

func NewSenders(plainSender *Plain, emailSender *Email) *Senders {
	return &Senders{
		EmailSender: emailSender,
		PlainSender: plainSender,
	}
}

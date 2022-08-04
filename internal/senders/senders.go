package senders

type Senders struct {
	EmailSender *Email
}

func NewSenders(emailSender *Email) *Senders {
	return &Senders{
		EmailSender: emailSender,
	}
}

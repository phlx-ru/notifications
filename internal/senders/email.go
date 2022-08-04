package senders

import (
	"fmt"
	"net/smtp"
	"net/url"

	"github.com/jordan-wright/email"
)

type Email struct {
	From              string
	RelayAddress      string
	RelayAuthUsername string
	RelayAuthPassword string
	RelayAuthHost     string
}

func (e *Email) SendText(to []string, subject, body string) error {
	return e.send(&email.Email{
		To:      to,
		From:    e.From,
		Subject: subject,
		Text:    []byte(body),
	})
}

func (e *Email) SendHTML(to []string, subject, body string) error {
	return e.send(&email.Email{
		To:      to,
		From:    e.From,
		Subject: subject,
		HTML:    []byte(body),
	})
}

func (e *Email) send(mail *email.Email) error {
	auth := smtp.PlainAuth("", e.RelayAuthUsername, e.RelayAuthPassword, e.RelayAuthHost)
	return mail.Send(e.RelayAddress, auth)
}

func NewEmail(from, addr, username, password string) (*Email, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to parse smtp relay address '%s': %w", addr, err)
	}
	if u.Scheme == "" {
		return nil, fmt.Errorf("smtp relay address '%s' is incorrect", addr)
	}
	return &Email{
		From:              from,
		RelayAddress:      addr,
		RelayAuthUsername: username,
		RelayAuthPassword: password,
		RelayAuthHost:     u.Scheme,
	}, nil
}

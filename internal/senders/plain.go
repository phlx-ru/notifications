package senders

import (
	"os"

	"github.com/go-kratos/kratos/v2/log"
)

type PlainSender interface {
	Send(message string)
}

type LogWriter interface {
	Infof(format string, a ...interface{})
}

type Plain struct {
	writer LogWriter
}

func FromPath(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
}

func NewPlain(file *os.File) *Plain {
	logs := log.NewHelper(
		log.With(
			log.NewStdLogger(file),
			"ts", log.DefaultTimestamp,
			"sender", "plain",
		),
	)
	return &Plain{writer: logs}
}

func (p *Plain) Send(message string) {
	p.writer.Infof(`received message [%s]`, message)
}

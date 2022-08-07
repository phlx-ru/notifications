package utils

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var (
	// ErrTermSig returns when OS caught one of signals SIGINT, SIGTERM
	ErrTermSig = errors.New("termination signal caught")
)

// SignalTrap use for intercept OS signal and shutdown process, which use orchestration
// tools (like golang.org/x/sync/errgroup)
type SignalTrap chan os.Signal

// TermSignalTrap returns SignalTrap which caught OS signals
func TermSignalTrap() SignalTrap {
	trap := SignalTrap(make(chan os.Signal, 1))

	signal.Notify(trap, syscall.SIGINT, os.Interrupt, syscall.SIGTERM)

	return trap
}

// Wait is waiting for OS signal and return ErrTermSig if shutdown signal was caught
func (t SignalTrap) Wait(ctx context.Context) error {
	select {
	case <-t:
		return ErrTermSig
	case <-ctx.Done():
		return ctx.Err()
	}
}

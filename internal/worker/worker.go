package worker

import (
	"context"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"notifications/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

const (
	notificationsLimit = 10 // limit of notifications processing at one time

	maxConcurrentWorkers = 10

	sleepDuration = time.Second
)

type Worker struct {
	usecase *biz.NotificationUsecase

	logger *log.Helper

	runOnce bool
}

type Option func(w *Worker)

func RunOnceOption() Option {
	return func(w *Worker) {
		w.runOnce = true
	}
}

func New(u *biz.NotificationUsecase, l log.Logger, options ...Option) *Worker {
	w := &Worker{
		usecase: u,
		logger:  log.NewHelper(l),
	}
	for _, option := range options {
		option(w)
	}
	return w
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			w.logger.Info(`worker get done signal`)
			return nil
		default:
		}
		count, err := w.usecase.CountOfPendingNotifications(ctx)
		if err != nil {
			w.logger.Errorf(`failed to count waiting notifications: %v`, err)
			return err
		}
		if count == 0 {
			if w.runOnce {
				return nil
			}
			w.logger.Infof("primary count = 0, sleeping for %d seconds", int(sleepDuration.Seconds()))
			time.Sleep(sleepDuration)
			continue
		}
		goroutines := int(
			math.Min(
				float64(maxConcurrentWorkers),
				math.Ceil(float64(count)/float64(notificationsLimit)),
			),
		)
		found := int64(0)
		processed := int64(0)
		wg := sync.WaitGroup{}
		for ; goroutines != 0; goroutines-- {
			wg.Add(1)
			go func() {
				defer wg.Done()

				currentFound, currentProcessed, err := w.usecase.ProcessNotifications(ctx, notificationsLimit)
				if err != nil {
					w.logger.Warnf(`error run once notification process: %v`, err)
				}
				atomic.AddInt64(&found, currentFound)
				atomic.AddInt64(&processed, currentProcessed)
			}()
		}
		wg.Wait()

		w.logger.Infof(
			"process iteration complete: primary count = %d, found = %d, processed = %d",
			count,
			found,
			processed,
		)
		if w.runOnce {
			return nil
		}
	}
}

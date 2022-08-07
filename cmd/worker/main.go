package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"path"
	"syscall"
	"time"

	"notifications/internal/biz"
	"notifications/internal/conf"
	"notifications/internal/senders"
	"notifications/internal/utils"
	"notifications/internal/worker"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/joho/godotenv"
	"github.com/vrecan/death/v3"
	"gopkg.in/alexcesaro/statsd.v2"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = `notifications_worker`
	// Version is the version of the compiled software.
	Version = `0.0.1`
	// flagconf is the config flag.
	flagconf string
	// dotenv is loaded from config path .env file
	dotenv string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&dotenv, "dotenv", ".env", ".env file, eg: -dotenv .env.local")
}

func newWorker(u *biz.NotificationUsecase, l log.Logger) *worker.Worker {
	return worker.New(u, l)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

	dead := death.NewDeath(syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		dead.WaitForDeathWithFunc(cancelFunc)
	}()

	envPath := path.Join(flagconf, dotenv)
	err = godotenv.Overload(envPath)
	if err != nil {
		return err
	}

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(conf.EnvDecoder),
	)
	defer func() {
		_ = c.Close()
	}()

	if err = c.Load(); err != nil {
		return err
	}

	var bc conf.Bootstrap
	if err = c.Scan(&bc); err != nil {
		return err
	}

	loggerInstance := log.With(
		log.DefaultLogger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", id,
		"service.name", Name,
		"service.version", Version,
		"trace.id", tracing.TraceID(),
		"span.id", tracing.SpanID(),
	)
	logLevel := log.ParseLevel(bc.Log.Level)
	logger := log.NewFilter(loggerInstance, log.FilterLevel(logLevel))
	log.SetLogger(logger)
	logs := log.NewHelper(logger)

	metrics, err := statsd.New(
		statsd.Address(bc.Metrics.Address),
		statsd.ErrorHandler(
			func(err error) {
				logs.Warnf(`failed to send metrics: %v`, err)
			},
		),
		statsd.Mute(bc.Metrics.Mute),
		statsd.Prefix(Name),
		statsd.Tags("id", id, "name", Name, "version", Version),
		statsd.FlushPeriod(time.Second),
	)
	if metrics == nil {
		if err != nil {
			return err
		}
		return errors.New("metrics client is undefined")
	}
	if err != nil && !bc.Metrics.Mute {
		return err
	}
	defer metrics.Close()
	metrics.Increment("starts.count")

	es := bc.Senders.GetEmail()

	emailSender, err := senders.NewEmail(es.GetFrom(), es.GetAddress(), es.GetUsername(), es.GetPassword())
	if err != nil {
		return err
	}

	plainFilePath := bc.Senders.Plain.GetFile()
	plainFile, err := senders.FromPath(plainFilePath)
	if err != nil {
		return err
	}
	plainSender := senders.NewPlain(plainFile)

	sendersSet := senders.NewSenders(plainSender, emailSender)

	wrkr, cleanup, err := wireWorker(bc.Data, sendersSet, logger)
	if err != nil {
		log.Errorf("failed to wire worker: %v", err)
		return nil
	}
	defer cleanup()

	logs.Info("worker start")
	err = wrkr.Run(ctx)
	if err != nil && err != utils.ErrTermSig && err != context.Canceled {
		logs.Fatalf("worker failed: %v", err)

	}
	logs.Info("worker ends successfully")

	return nil
}

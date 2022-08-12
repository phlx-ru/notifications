package main

import (
	"context"
	"flag"
	"os"
	"path"
	"syscall"

	"notifications/internal/biz"
	"notifications/internal/clients/telegram"
	"notifications/internal/conf"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"
	"notifications/internal/pkg/runtime"
	"notifications/internal/pkg/transport"
	"notifications/internal/senders"
	"notifications/internal/worker"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/joho/godotenv"
	"github.com/vrecan/death/v3"
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
	flag.StringVar(&flagconf, "conf", "./configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&dotenv, "dotenv", ".env", ".env file, eg: -dotenv .env")
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

	ctx, cancelFunc := context.WithCancel(context.Background())
	dead := death.NewDeath(syscall.SIGINT, os.Interrupt, syscall.SIGTERM)
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

	logs := logger.New(id, Name, Version, bc.Log.Level)
	logHelper := logger.NewHelper(logs, "scope", "worker")
	dead.SetLogger(logHelper)

	metric, err := metrics.New(bc.Metrics.Address, Name, bc.Metrics.Mute)
	if err != nil {
		return err
	}
	defer metric.Close()
	metric.Increment("starts.count")

	database, cleanup, err := wireData(bc.Data, logs)
	if err != nil {
		return err
	}
	defer cleanup()

	go database.CollectDatabaseMetrics(ctx, metric, id)
	go runtime.CollectGoMetrics(ctx, metric, id)

	es := bc.Senders.GetEmail()
	emailSender, err := senders.NewEmail(es.From, es.Address, es.Username, es.Password, metric, logs)
	if err != nil {
		return err
	}

	plainFilePath := bc.Senders.Plain.GetFile()
	plainFile, err := senders.FromPath(plainFilePath)
	if err != nil {
		return err
	}
	plainSender := senders.NewPlain(plainFile, metric, logs)

	httpClient := transport.NewHTTPClient()
	telegramClient := telegram.New(bc.Senders.Telegram.BotToken, httpClient, metric, logs)
	telegramSender := senders.NewTelegram(telegramClient, metric, logs)

	sendersSet := senders.NewSenders(plainSender, emailSender, telegramSender)

	wrkr, err := wireWorker(database, sendersSet, metric, logs)
	if err != nil {
		log.Errorf("failed to wire worker: %v", err)
		return nil
	}

	logHelper.Info("worker start")
	err = wrkr.Run(ctx)
	if err != nil && err != context.Canceled {
		logHelper.Fatalf("worker failed: %v", err)

	}
	logHelper.Info("worker ends successfully")

	return nil
}

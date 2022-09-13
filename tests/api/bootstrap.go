package api

import (
	"context"
	"path"

	"notifications/internal/auth"
	"notifications/internal/biz"
	"notifications/internal/clients/smsaero"
	"notifications/internal/clients/telegram"
	"notifications/internal/conf"
	"notifications/internal/data"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"
	"notifications/internal/pkg/runtime"
	"notifications/internal/pkg/transport"
	"notifications/internal/senders"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
)

const (
	flagconf = `../../configs`
	dotenv   = `.env.testing`
	id       = `api`
	name     = `api-tests`
	version  = `1.1.1`
)

var (
	bc               conf.Bootstrap
	database         data.Database
	logs             *log.Filter
	metric           metrics.Metrics
	sendersSet       *senders.Senders
	notificationRepo biz.NotificationRepo
	httpServer       *http.Server
	jwtToken         string
)

func bootstrap() (func(), error) {
	var err error

	ctx := context.Background()

	envPath := path.Join(flagconf, dotenv)
	err = godotenv.Overload(envPath)
	if err != nil {
		return nil, err
	}

	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
		config.WithDecoder(conf.EnvDecoder),
	)

	if err = c.Load(); err != nil {
		return nil, err
	}

	if err = c.Scan(&bc); err != nil {
		return nil, err
	}

	jwtToken = auth.MakeJWT(bc.Auth.Jwt.Secret)

	logs = logger.New(id, name, version, bc.Log.Level)

	metric, err = metrics.New(bc.Metrics.Address, name, bc.Metrics.Mute)
	if err != nil {
		return nil, err
	}

	var databaseCleanup func()
	database, databaseCleanup, err = wireData(bc.Data, logs)
	if err != nil {
		return nil, err
	}

	go database.CollectDatabaseMetrics(ctx, metric)
	go runtime.CollectGoMetrics(ctx, metric)

	es := bc.Senders.GetEmail()
	emailSender, err := senders.NewEmail(
		es.GetFrom(),
		es.GetAddress(),
		es.GetUsername(),
		es.GetPassword(),
		metric,
		logs,
	)
	if err != nil {
		return nil, err
	}

	plainFilePath := bc.Senders.Plain.GetFile()
	plainFile, err := senders.FromPath(plainFilePath)
	if err != nil {
		return nil, err
	}
	plainSender := senders.NewPlain(plainFile, metric, logs)

	httpClient := transport.NewHTTPClient()
	telegramClient := telegram.New(bc.Senders.Telegram.BotToken, httpClient, metric, logs)
	telegramSender := senders.NewTelegram(telegramClient, metric, logs)

	aero := bc.Senders.Sms.Aero
	smsAeroClient := smsaero.New(aero.Email, aero.ApiKey, httpClient, metric, logs)
	smsAeroSender := senders.NewSMSAero(smsAeroClient, metric, logs)

	sendersSet = senders.NewSenders(plainSender, emailSender, telegramSender, smsAeroSender)

	notificationRepo = wireNotificationRepo(database, logs, metric)

	httpServer = wireHTTPServer(database, bc.Server, bc.Auth, sendersSet, metric, logs)

	cleanup := func() {
		_ = c.Close()
		metric.Close()
		databaseCleanup()
	}

	return cleanup, nil
}

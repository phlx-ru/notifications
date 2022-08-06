package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"path"
	"time"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/joho/godotenv"
	"gopkg.in/alexcesaro/statsd.v2"

	"notifications/internal/conf"
	"notifications/internal/data"
	"notifications/internal/senders"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = `notifications`
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

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	flag.Parse()

	var err error

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
	defer c.Close()

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
	logs := log.NewHelper(logger)

	logs.Info("app started")

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

	database, cleanup, err := wireData(bc.Data, logger)
	if err != nil {
		return err
	}
	defer cleanup()

	err = data.Prepare(context.Background(), database, bc.Data.Database.Migrate)
	if err != nil {
		return err
	}

	es := bc.Senders.GetEmail()

	emailSender, err := senders.NewEmail(es.GetFrom(), es.GetAddress(), es.GetUsername(), es.GetPassword())
	if err != nil {
		return err
	}

	sendersSet := senders.NewSenders(emailSender)

	app, err := wireApp(database, bc.Server, sendersSet, logger)
	if err != nil {
		return err
	}

	// start and wait for stop signal
	if err = app.Run(); err != nil {
		return err
	}

	logs.Info("app terminates")

	return nil
}

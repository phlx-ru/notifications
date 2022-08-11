package main

import (
	"context"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"notifications/ent"
	"notifications/ent/schema"
	"notifications/internal/pkg/logger"
	"notifications/internal/pkg/metrics"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/joho/godotenv"

	"notifications/internal/conf"
	"notifications/internal/data"
)

//////////////////////////////////////
//                                  //
// HELP SCRIPT FOR TESTING PURPOSES //
//                                  //
//////////////////////////////////////

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name = `notifications_server`
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
	flag.StringVar(&dotenv, "dotenv", ".env.local", ".env file, eg: -dotenv .env.local")
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

	level := "warn"
	logs := logger.New(id, Name, Version, level)

	metric, err := metrics.New(bc.Metrics.Address, Name, true)
	if err != nil {
		return err
	}
	defer metric.Close()

	database, cleanup, err := wireData(bc.Data, logs)
	if err != nil {
		return err
	}
	defer cleanup()

	ctx := context.Background()
	return seed(ctx, database)
}

func seed(ctx context.Context, database *data.Data) error {
	return database.Seed(
		ctx, func(ctx context.Context, ent *ent.Client) error {
			tx, err := ent.Tx(ctx)
			if err != nil {
				return fmt.Errorf("failed to start transaction: %v", err)
			}

			rand.Seed(time.Now().UnixNano())
			//nolint:gosec // G404: Use of weak random number generator (math/rand instead of crypto/rand)
			serial := rand.Int63()

			for i := 0; i < 1000; i++ {
				_, err := ent.Notification.Create().
					SetSenderID(0).
					SetStatus(schema.StatusPending).
					SetTTL(40).
					SetType(schema.TypePlain).
					SetPayload(
						map[string]string{
							"message": fmt.Sprintf("[serial=%d] notification number %d", serial, i),
						},
					).
					Save(ctx)

				if err != nil {
					rbErr := tx.Rollback()
					if rbErr != nil {
						return fmt.Errorf("failed to rollback with [%v] after [%v]", rbErr, err)
					}
					return fmt.Errorf("failed to make notification: %v", err)
				}
			}

			return tx.Commit()
		},
	)
}

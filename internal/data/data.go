package data

import (
	"context"
	"fmt"

	"notifications/ent"
	"notifications/internal/conf"

	"entgo.io/ent/dialect/sql/schema"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/lib/pq"
)

// ProviderRepoSet is data providers.
var ProviderRepoSet = wire.NewSet(NewGreeterRepo, NewNotificationRepo)

var ProviderDataSet = wire.NewSet(NewData)

// Data .
type Data struct {
	ent    *ent.Client
	logger *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "ent/data/logger-job"))

	client, err := ent.Open(
		c.Database.Driver,
		c.Database.Source,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed opening connection to db: %w", err)
	}

	cleanup := func() {
		logHelper.Info("cleaning database from ent client")
		if client != nil {
			err := client.Close()
			if err != nil {
				logHelper.Errorf(`failed to close database client: %w`, err)
			}
		}
	}
	return &Data{
		ent:    client,
		logger: logHelper,
	}, cleanup, nil
}

// MigrateSoft only creates and updates schema entities
func (d *Data) MigrateSoft(ctx context.Context) error {
	err := d.ent.Schema.Create(ctx, schema.WithForeignKeys(false))
	if err != nil {
		d.logger.WithContext(ctx).Errorf(`failed to soft migrate database schema: %w`, err)
		return err
	}
	return nil
}

// MigrateHard does same as MigrateSoft, but also drop columns and indices
func (d *Data) MigrateHard(ctx context.Context) error {
	err := d.ent.Schema.Create(ctx, schema.WithDropIndex(true), schema.WithDropColumn(true))
	if err != nil {
		d.logger.WithContext(ctx).Errorf(`failed to hard migrate database schema: %w`, err)
		return err
	}
	return nil
}

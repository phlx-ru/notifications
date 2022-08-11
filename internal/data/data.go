package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"notifications/ent"
	"notifications/internal/conf"
	"notifications/internal/pkg/metrics"

	entDialectSQL "entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/schema"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	_ "github.com/lib/pq" // postgres driver for Go's database/sql package
)

const (
	maxOpenConnections = 32
	maxIdleConnections = 30
	maxConnLifetime    = 5 * time.Minute
	sendStatsEvery     = time.Second
)

// ProviderRepoSet is data providers.
var ProviderRepoSet = wire.NewSet(NewNotificationRepo)

var ProviderDataSet = wire.NewSet(NewData)

// Data .
type Data struct {
	db     *sql.DB
	ent    *ent.Client
	logger *log.Helper
}

// NewData .
func NewData(c *conf.Data, logger log.Logger) (*Data, func(), error) {
	logHelper := log.NewHelper(log.With(logger, "module", "ent/data/logger-job"))

	drv, err := entDialectSQL.Open(c.Database.Driver, c.Database.Source)
	if err != nil {
		return nil, nil, err
	}
	// Get the underlying sql.DB object of the driver.
	db := drv.DB()
	db.SetMaxIdleConns(maxIdleConnections)
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetConnMaxLifetime(maxConnLifetime)
	options := []ent.Option{
		ent.Driver(drv),
	}
	if c.Database.Debug {
		options = append(options, ent.Debug())
	}
	client := ent.NewClient(options...)

	cleanup := func() {
		logHelper.Info("closing database client from cleanup() function")
		if client != nil {
			err := client.Close()
			if err != nil {
				logHelper.Errorf(`failed to close database client: %v`, err)
			}
		}
	}
	return &Data{
		db:     db,
		ent:    client,
		logger: logHelper,
	}, cleanup, nil
}

// MigrateSoft only creates and updates schema entities
func (d *Data) MigrateSoft(ctx context.Context) error {
	err := d.ent.Schema.Create(ctx, schema.WithForeignKeys(false))
	if err != nil {
		d.logger.WithContext(ctx).Errorf(`failed to soft migrate database schema: %v`, err)
		return err
	}
	return nil
}

// MigrateHard does same as MigrateSoft, but also drop columns and indices
func (d *Data) MigrateHard(ctx context.Context) error {
	err := d.ent.Schema.Create(ctx, schema.WithDropIndex(true), schema.WithDropColumn(true))
	if err != nil {
		d.logger.WithContext(ctx).Errorf(`failed to hard migrate database schema: %v`, err)
		return err
	}
	return nil
}

func (d *Data) Prepare(ctx context.Context, m conf.Data_Database_Migrate) error {
	if m == conf.Data_Database_none {
		return nil
	}
	if m == conf.Data_Database_soft {
		d.logger.WithContext(ctx).Info("preparing database: running soft migrate")
		return d.MigrateSoft(ctx)
	}
	if m == conf.Data_Database_hard {
		d.logger.WithContext(ctx).Info("preparing database: running hard migrate")
		return d.MigrateHard(ctx)
	}
	return nil
}

func (d *Data) CollectDatabaseMetrics(ctx context.Context, metric metrics.Metrics, id string) {
	metricPrefix := fmt.Sprintf(`postgres.%s.connections.`, id)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		stats := d.db.Stats()

		// The number of established connections both in use and idle.
		metric.Gauge(metricPrefix+`open`, stats.OpenConnections)

		// The number of connections currently in use.
		metric.Gauge(metricPrefix+`used`, stats.InUse)

		// The number of idle connections.
		metric.Gauge(metricPrefix+`idle`, stats.Idle)

		// The total number of connections waited for.
		metric.Gauge(metricPrefix+`wait`, stats.WaitCount)

		// The total time blocked waiting for a new connection.
		// metric.Gauge(metricPrefix+`wait_duration`, stats.WaitDuration) // TODO Duration or count ms?

		// The total number of connections closed due to SetMaxIdleConns.
		metric.Gauge(metricPrefix+`max_idle_closed`, stats.MaxIdleClosed)

		// The total number of connections closed due to SetConnMaxIdleTime.
		metric.Gauge(metricPrefix+`max_idle_time_closed`, stats.MaxIdleTimeClosed)

		// The total number of connections closed due to SetConnMaxLifetime.
		metric.Gauge(metricPrefix+`max_lifetime_closed`, stats.MaxLifetimeClosed)

		time.Sleep(sendStatsEvery)
	}
}

// Seed everything you need by passing the seeding func
func (d *Data) Seed(ctx context.Context, seeding func(context.Context, *ent.Client) error) error {
	return seeding(ctx, d.ent)
}

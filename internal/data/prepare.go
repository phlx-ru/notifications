package data

import (
	"context"

	"notifications/internal/conf"
)

func Prepare(ctx context.Context, d *Data, m conf.Data_Database_Migrate) error {
	if m == conf.Data_Database_none {
		return nil
	}
	if m == conf.Data_Database_soft {
		d.logger.Info("preparing database: running soft migrate")
		return d.MigrateSoft(ctx)
	}
	if m == conf.Data_Database_hard {
		d.logger.Info("preparing database: running hard migrate")
		return d.MigrateHard(ctx)
	}
	return nil
}

// func Stats // TODO MAKE DATABASE STATS

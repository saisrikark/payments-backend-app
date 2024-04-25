package migrate

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

func Run(ctx context.Context, db *bun.DB) error {

	db.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS bun_migrations (
		id SERIAL PRIMARY KEY,
		name TEXT,
		group_id INT,
		migrated_at TIMESTAMP
	);`)

	migrator := migrate.NewMigrator(db, Migrations)
	_, err := migrator.Migrate(ctx)
	if err != nil {
		return fmt.Errorf("unable to migrate [%s]", err.Error())
	}

	return nil
}

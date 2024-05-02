package migrate

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {

		var err error

		_, err = db.ExecContext(ctx, `
			ALTER TABLE transaction ADD COLUMN IF NOT EXISTS balance DOUBLE PRECISION;
		`)
		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		return nil
	})
}

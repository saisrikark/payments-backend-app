package migrate

import (
	"context"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {

		var err error

		_, err = db.ExecContext(ctx, `
		CREATE TABLE account (
			id serial PRIMARY KEY NOT NULL,
			document_id VARCHAR,
			unique(document_id)
			);
		`)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, `
		CREATE TABLE operation_type (
			id serial PRIMARY KEY NOT NULL,
			description varchar
			);
		`)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, `
		CREATE TABLE transaction (
			id serial PRIMARY KEY NOT NULL,
			account_id integer references account (id) NOT NULL,
			operations_type_id integer references operation_type,
			amount DOUBLE PRECISION
			);
		`)
		if err != nil {
			return err
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		return nil
	})
}

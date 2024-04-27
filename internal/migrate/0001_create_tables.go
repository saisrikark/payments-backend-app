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
			document_number VARCHAR,
			unique(document_number)
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
		INSERT INTO operation_type (id, description) VALUES (1, 'Normal Purchase');
		INSERT INTO operation_type (id, description) VALUES (2, 'Purchase with installments');
		INSERT INTO operation_type (id, description) VALUES (3, 'Withdrawal');
		INSERT INTO operation_type (id, description) VALUES (4, 'Credit Voucher');
		`)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, `
		CREATE TABLE transaction (
			id serial PRIMARY KEY NOT NULL,
			account_id integer references account (id) NOT NULL,
			operation_type_id integer references operation_type (id) NOT NULL,
			amount DOUBLE PRECISION,
			event_date TIMESTAMP NOT NULL
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

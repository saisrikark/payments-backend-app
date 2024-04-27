package models

import (
	"context"
	"database/sql"
	"errors"
	"payments-backend-app/pkg/models"
	"strings"

	"github.com/uptrace/bun"
)

type accountsService struct {
	db *bun.DB
}

func NewAccountsService(db *bun.DB) *accountsService {
	return &accountsService{
		db: db,
	}
}

func (as *accountsService) Create(ctx context.Context, account models.Account) (models.Account, error) {

	raccount := models.Account{}

	err := as.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewInsert().Model(&account).Exec(ctx)
		if err != nil {
			return err
		}

		if err := tx.NewSelect().Model(&raccount).Where("document_number LIKE ?", account.DocumentNumber).Scan(ctx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate key value violates unique constraint"):
			err = models.DuplicateRecordErr
		case errors.Is(err, sql.ErrNoRows):
			err = models.NoRecordErr
		}
	}

	return raccount, err
}

func (as *accountsService) GetForID(ctx context.Context, accountID int64) (models.Account, error) {

	raccount := models.Account{}

	err := as.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		if err := tx.NewSelect().Model(&raccount).Where("id = ?", accountID).Scan(ctx); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = models.NoRecordErr
		}
	}

	return raccount, err
}

func (as *accountsService) DeleteForID(ctx context.Context, accountID int64) error {

	err := as.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		_, err := tx.NewDelete().Model(&models.Account{}).Where("id = ?", accountID).Exec(ctx)

		return err
	})

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = models.NoRecordErr
		}
	}

	return err
}

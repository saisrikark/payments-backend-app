package models

import (
	"context"
	"database/sql"
	"errors"
	"payments-backend-app/pkg/models"
	"time"

	"github.com/uptrace/bun"
)

type transactionService struct {
	db *bun.DB
}

func NewTransactionService(db *bun.DB) *transactionService {
	return &transactionService{
		db: db,
	}
}

func (ts *transactionService) Create(ctx context.Context, transaction models.Transaction) (models.TransactionStatus, error) {

	transactionStatus := models.TransactionStatus{}
	rtransaction := models.Transaction{}

	err := ts.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		if err := tx.NewSelect().Model(&models.Account{}).Where("id = ?", transaction.AccountID).Scan(ctx); err != nil {
			return err
		}

		transaction.EventDate = time.Now()

		_, err := tx.NewInsert().Model(&transaction).Exec(ctx)
		if err != nil {
			return err
		}

		if err := tx.NewSelect().Model(&rtransaction).Where("event_date = ?", transaction.EventDate).Scan(ctx); err != nil {
			return err
		}

		transactionStatus.TransactionID = rtransaction.ID
		transactionStatus.AccountID = rtransaction.AccountID

		return nil
	})

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			err = models.NoRecordErr
		}
	}

	return transactionStatus, err
}

func (ts *transactionService) GetForID(ctx context.Context, transactionID int64) (models.Transaction, error) {

	rtransaction := models.Transaction{}

	err := ts.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		if err := tx.NewSelect().Model(&rtransaction).Where("id = ?", transactionID).Scan(ctx); err != nil {
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

	return rtransaction, err
}

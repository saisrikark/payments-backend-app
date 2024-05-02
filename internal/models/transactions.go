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

	unresolvedTransactions := []models.Transaction{}

	err := ts.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {

		if err := tx.NewSelect().Model(&models.Account{}).Where("id = ?", transaction.AccountID).For("UPDATE").Scan(ctx); err != nil {
			return err
		}

		currBalance := transaction.Amount

		if transaction.Amount > 0 {

			// query fetches all transactions with balance < 0
			count, err := tx.NewSelect().
				Model(&unresolvedTransactions).
				Where("balance < 0").
				Where("account_id = ?", transaction.AccountID).
				OrderExpr("event_date ASC").ScanAndCount(ctx)
			if err != nil {
				if !errors.Is(err, sql.ErrNoRows) {
					return err
				}
			}

			if count > 0 {
				// for each transaction, see if we complete the balance and update in db
				for _, unresolvedTransaction := range unresolvedTransactions {
					transactionRemainingBalance := 0.0

					if currBalance > 0 {
						// transaction is resolved set to 0
						if unresolvedTransaction.Balance+currBalance > 0 {
							transactionRemainingBalance = 0.0
							currBalance = currBalance + unresolvedTransaction.Balance
						} else {
							transactionRemainingBalance = unresolvedTransaction.Balance + currBalance
							currBalance = 0
						}

					} else {
						break
					}

					// push to db
					_, err := tx.NewUpdate().Model(&unresolvedTransaction).
						Set("balance = ?", transactionRemainingBalance).
						Where("id = ?", unresolvedTransaction.ID).
						Exec(ctx)
					if err != nil {
						return err
					}
				}
			}

		}

		transaction.EventDate = time.Now()

		transaction.Balance = currBalance
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

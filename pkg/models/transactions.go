package models

import "context"

type TransactionService interface {
	Create(ctx context.Context, transaction Transaction) (TransactionStatus, error)
	GetForID(ctx context.Context, transactionID int64) (Transaction, error)
}

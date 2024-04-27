package models

import "context"

type TransactionService interface {
	Create(ctx context.Context, transaction Transaction) (TransactionStatus, error)
}

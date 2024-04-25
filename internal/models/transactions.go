package models

import (
	"context"
	"payments-backend-app/pkg/models"

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

func (ts *transactionService) Create(ctx context.Context, transaction models.Transaction) (models.Transaction, error) {
	return models.Transaction{}, nil
}

package models

import (
	"context"
	"payments-backend-app/pkg/models"

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
	return models.Account{}, nil
}

func (ts *accountsService) GetForID(ctx context.Context, accountID int64) (models.Account, error) {
	return models.Account{}, nil
}

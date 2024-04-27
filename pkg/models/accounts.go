package models

import "context"

type AccountsService interface {
	Create(ctx context.Context, account Account) (Account, error)
	GetForID(ctx context.Context, accountID int64) (Account, error)
	DeleteForID(ctx context.Context, accountID int64) error
}

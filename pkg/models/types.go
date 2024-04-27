package models

import (
	"time"

	"github.com/uptrace/bun"
)

type Account struct {
	bun.BaseModel `bun:"table:account,alias:a"`

	AccountID      int64  `json:"account_id" bun:"id,autoincrement"`
	DocumentNumber string `json:"document_number" bun:"document_number"`
}

type OperationType struct {
	ID          int64  `json:"id"`
	Description string `json:"description"`
}

type OperationTypeID int

const (
	NormalPurchase OperationTypeID = iota + 1
	PurchaseWithInstallments
	Withdrawal
	CreditVoucher
)

func IsSupportedType(operationType int) bool {
	switch {
	case operationType > 4 || operationType < 1:
		return false
	default:
		return true
	}
}

func IsCredit(operationType int) bool {
	switch {
	case operationType == int(CreditVoucher):
		return true
	default:
		return false
	}
}

type Transaction struct {
	bun.BaseModel `bun:"table:transaction,alias:t"`

	ID              int64     `json:"id" bun:"id,autoincrement"`
	AccountID       int64     `json:"account_id" bun:"account_id"`
	OperationTypeID int64     `json:"operation_type_id" bun:"operation_type_id"`
	Amount          float64   `json:"amount" bun:"amount"`
	EventDate       time.Time `json:"event_date" bun:"event_date"`
}

type TransactionStatus struct {
}

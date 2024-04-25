package models

import "time"

type Account struct {
	DocumentID string `json:"document_id"`
}

type OperationType struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
}

type Transaction struct {
	ID              int64     `json:"id"`
	AccountID       int64     `json:"account_id"`
	OperationTypeID int       `json:"operation_type_id"`
	Amount          float64   `json:"amount"`
	EventDate       time.Time `json:"event_date"`
}

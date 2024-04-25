package server

// write an unmarshaller
type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountID       int64 `json:"account_id"`
	OperationTypeID int64 `json:"operation_type_id"`
}

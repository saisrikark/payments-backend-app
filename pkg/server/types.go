package server

// write an unmarshaller
type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

type CreateAccountResponse struct {
	AccountID int64 `json:"account_id"`
}

type CreateTransactionRequest struct {
	AccountID       int64 `json:"account_id"`
	OperationTypeID int64 `json:"operation_type_id"`
}

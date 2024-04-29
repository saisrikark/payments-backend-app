package server

import (
	"encoding/json"
	"fmt"
	"payments-backend-app/pkg/models"
	"strconv"
	"strings"
)

type CreateAccountRequest struct {
	DocumentNumber string `json:"document_number"`
}

func (c *CreateAccountRequest) UnmarshalJSON(data []byte) error {

	var createAccountRequest struct {
		DocumentNumber string `json:"document_number"`
	}

	if err := json.Unmarshal(data, &createAccountRequest); err != nil {
		return err
	}

	documentNumber := createAccountRequest.DocumentNumber

	switch {
	case documentNumber != strings.TrimSpace(documentNumber):
		return fmt.Errorf("document number has trailing spaces")
	case documentNumber != strings.TrimLeft(documentNumber, "0"):
		return fmt.Errorf("document number has 0's in the beginning")
	case len(documentNumber) > 15:
		return fmt.Errorf("document number length must be no greater than 15")
	case len(documentNumber) == 0:
		return fmt.Errorf("empty document number not allowed")
	}

	_, err := strconv.Atoi(documentNumber)
	if err != nil {
		return fmt.Errorf("invalid number")
	}

	c.DocumentNumber = documentNumber
	return nil
}

type GetAccountResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

type CreateTransactionRequest struct {
	AccountID       int64   `json:"account_id"`
	OperationTypeID int64   `json:"operation_type_id"`
	Amount          float64 `json:"amount"`
}

func (c *CreateTransactionRequest) UnmarshalJSON(data []byte) error {

	var createTransactionRequest struct {
		AccountID       int64   `json:"account_id"`
		OperationTypeID int64   `json:"operation_type_id"`
		Amount          float64 `json:"amount"`
	}

	if err := json.Unmarshal(data, &createTransactionRequest); err != nil {
		return err
	}

	if !models.IsSupportedType(int(createTransactionRequest.OperationTypeID)) {
		return fmt.Errorf("unsupported operation type")
	}

	amountS := fmt.Sprintf("%f", createTransactionRequest.Amount)
	decimal := strings.Trim(strings.Split(amountS, ".")[1], "0")

	switch {
	case len(decimal) > 2:
		return fmt.Errorf("amount must be capped to 2 decimal places")
	}

	c.AccountID = createTransactionRequest.AccountID
	c.OperationTypeID = createTransactionRequest.OperationTypeID

	if models.IsCredit(int(c.OperationTypeID)) {
		c.Amount = createTransactionRequest.Amount
	} else {
		c.Amount = -createTransactionRequest.Amount
	}

	return nil
}

type CreateTransactionResponse struct {
	TransactionID int64 `json:"transaction_id"`
	AccountID     int64 `json:"account_id"`
}

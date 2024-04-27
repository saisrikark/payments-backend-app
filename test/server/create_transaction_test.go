package server

import (
	"context"
	"fmt"
	"net/http"
	"payments-backend-app/pkg/models"
	"payments-backend-app/pkg/server"
	"payments-backend-app/test/testutils"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	ctx := context.Background()
	testServer := testutils.NewTestServer(t)
	testServer.Start(ctx)
	defer testServer.Stop(ctx)
	testServer.WaitForRunnning(t)

	t.Run("Create transaction", func(t *testing.T) {

		accountID := testutils.GenerateRandomNumberInt(10)

		if err := testServer.AccountsService.DeleteForID(ctx, int64(accountID)); err != nil {
			t.Errorf("unable to delete existing id [%s]", err)
		}

		account, err := testServer.AccountsService.Create(ctx, models.Account{
			DocumentNumber: testutils.GenerateRandomNumber(10),
		})
		if err != nil {
			t.Errorf("unable to create account [%s]", err)
		}

		for _, i := range []int{1, 2, 3, 4} {

			t.Run(fmt.Sprintf("operation type %d", i), func(t *testing.T) {

				req := &server.CreateTransactionRequest{
					AccountID:       account.AccountID,
					OperationTypeID: int64(i),
					Amount:          123.11,
				}

				status, _, err := testServer.CallCreateTransaction(req)
				if err != nil {
					t.Errorf("error creating the transaction [%s]", err)
				}

				if status != http.StatusCreated {
					t.Errorf("expected status %d got %d", http.StatusCreated, status)
				}
			})

		}

	})

	t.Run("Bad requests", func(t *testing.T) {

		account, err := testServer.AccountsService.Create(ctx, models.Account{
			DocumentNumber: testutils.GenerateRandomNumber(10),
		})
		if err != nil {
			t.Errorf("unable to create account [%s]", err)
		}

		t.Run("Incorrect operation ID", func(t *testing.T) {

			req := &server.CreateTransactionRequest{
				AccountID:       account.AccountID,
				OperationTypeID: 5,
				Amount:          123.11,
			}

			status, _, err := testServer.CallCreateTransaction(req)
			if err != nil {
				t.Errorf("error creating the transaction [%s]", err)
			}

			if status != http.StatusBadRequest {
				t.Errorf("expected status %d got %d", http.StatusBadRequest, status)
			}

		})

		t.Run("Amount decimal places greater than 2", func(t *testing.T) {

			req := &server.CreateTransactionRequest{
				AccountID:       account.AccountID,
				OperationTypeID: 4,
				Amount:          123.111,
			}

			status, _, err := testServer.CallCreateTransaction(req)
			if err != nil {
				t.Errorf("error creating the transaction [%s]", err)
			}

			if status != http.StatusBadRequest {
				t.Errorf("expected status %d got %d", http.StatusBadRequest, status)
			}

		})

	})

	t.Run("Fail to create transaction for non existent account", func(t *testing.T) {

		req := &server.CreateTransactionRequest{
			AccountID:       int64(testutils.GenerateRandomNumberInt(10)),
			OperationTypeID: 1,
			Amount:          123.11,
		}

		status, _, err := testServer.CallCreateTransaction(req)
		if err != nil {
			t.Errorf("error creating the transaction [%s]", err)
		}

		if status != http.StatusNotFound {
			t.Errorf("expected status %d got %d", http.StatusNotFound, status)
		}

	})
}
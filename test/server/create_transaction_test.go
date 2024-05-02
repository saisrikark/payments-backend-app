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

				status, resp, err := testServer.CallCreateTransaction(req)
				if err != nil {
					t.Errorf("error creating the transaction [%s]", err)
				}

				if status != http.StatusCreated {
					t.Errorf("expected status %d got %d", http.StatusCreated, status)
				}

				if resp != nil {
					if resp.TransactionID == 0 {
						t.Errorf("0 transaction id")
					} else {
						transaction, err := testServer.TransactionService.GetForID(ctx, resp.TransactionID)
						if err != nil {
							t.Errorf("unable to fetch created transaction")
						}

						if resp.TransactionID != transaction.ID {
							t.Errorf("transaction ids not matching resp %d database %d", resp.TransactionID, transaction.ID)
						}

						switch transaction.OperationTypeID {
						case 1, 2, 3:
							if (-req.Amount) != (transaction.Amount) {
								t.Errorf("expected -%f got %f ", req.Amount, transaction.Amount)
							}
						case 4:
							if (req.Amount) != (transaction.Amount) {
								t.Errorf("expected %f got %f ", req.Amount, transaction.Amount)
							}
						}
					}
				} else {
					t.Errorf("empty response body")
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

	// tests to handle remaining balance

	t.Run("Remaining balance", func(t *testing.T) {

		type TestData struct {
			OperationTypeID  int
			Amount           float64
			ExpectedBalances float64
		}

		testDatas := make([]TestData, 0)

		// testDatas = append(testDatas, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           50,
		// 	ExpectedBalances: -50.0,
		// }, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           23.5,
		// 	ExpectedBalances: -23.5,
		// }, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           18.7,
		// 	ExpectedBalances: -18.7,
		// })

		// testDatas = append(testDatas, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           50,
		// 	ExpectedBalances: 0.0,
		// }, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           23.5,
		// 	ExpectedBalances: -13.5,
		// }, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           18.7,
		// 	ExpectedBalances: -18.7,
		// }, TestData{
		// 	OperationTypeID:  4,
		// 	Amount:           60,
		// 	ExpectedBalances: 0.0,
		// })

		// testDatas = append(testDatas, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           50,
		// 	ExpectedBalances: 0.0,
		// }, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           23.5,
		// 	ExpectedBalances: 0.0,
		// }, TestData{
		// 	OperationTypeID:  1,
		// 	Amount:           18.7,
		// 	ExpectedBalances: 0.0,
		// }, TestData{
		// 	OperationTypeID:  4,
		// 	Amount:           60,
		// 	ExpectedBalances: 0.0,
		// }, TestData{
		// 	OperationTypeID:  4,
		// 	Amount:           100,
		// 	ExpectedBalances: 67.8,
		// })

		account, err := testServer.AccountsService.Create(ctx, models.Account{
			DocumentNumber: testutils.GenerateRandomNumber(10),
		})
		if err != nil {
			t.Errorf("unable to create account [%s]", err)
		}

		respIDs := make([]int64, 0)

		for _, testData := range testDatas {

			req := &server.CreateTransactionRequest{
				AccountID:       account.AccountID,
				OperationTypeID: int64(testData.OperationTypeID),
				Amount:          testData.Amount,
			}

			status, resp, err := testServer.CallCreateTransaction(req)
			if err != nil {
				t.Errorf("error creating the transaction [%s]", err)
			}

			if resp != nil {
				respIDs = append(respIDs, resp.TransactionID)
			} else {
				t.Errorf("empty response body")
			}

			if status != http.StatusCreated {
				t.Errorf("expected status %d got %d", http.StatusCreated, status)
			}

		}

		for i, id := range respIDs {
			transaction, err := testServer.TransactionService.GetForID(ctx, id)
			if err != nil {
				t.Errorf("unable to fetch created transaction")
			}

			if transaction.Balance != testDatas[i].ExpectedBalances {
				t.Errorf("operation type %d amount %f expected balance %f got %f",
					testDatas[i].OperationTypeID,
					testDatas[i].Amount,
					testDatas[i].ExpectedBalances,
					transaction.Balance)
			}
		}

	})
}

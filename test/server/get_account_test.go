package server

import (
	"context"
	"net/http"
	"payments-backend-app/pkg/models"
	"payments-backend-app/test/testutils"
	"testing"
)

func TestGetAccount(t *testing.T) {
	ctx := context.Background()
	testServer := testutils.NewTestServer(t)
	testServer.Start(ctx)
	defer testServer.Stop(ctx)
	testServer.WaitForRunnning(t)

	accountsService := testServer.AccountsService

	t.Run("Fetch account", func(t *testing.T) {

		documentNumber := testutils.GenerateRandomNumber(10)

		account, err := accountsService.Create(ctx, models.Account{
			DocumentNumber: documentNumber,
		})
		if err != nil {
			t.Errorf("unable to create account [%s]", err.Error())
		}

		status, raccount, err := testServer.CallGetAccount(int(account.AccountID))
		switch {
		case err != nil:
			t.Errorf("unable to fetch account from http request [%s]", err.Error())
		case status != http.StatusOK:
			t.Errorf("expected status %d got %d", http.StatusOK, status)
		case raccount == nil:
			t.Errorf("empty account response")
		case raccount.AccountID != account.AccountID:
			t.Errorf("expected account id %d got %d", account.AccountID, raccount.AccountID)
		case raccount.DocumentNumber != account.DocumentNumber:
			t.Errorf("expected document number %s got %s", account.DocumentNumber, raccount.DocumentNumber)
		}

	})

	t.Run("Fetch with an ID that does not exist", func(t *testing.T) {

		missingID := int64(testutils.GenerateRandomNumberInt(10))

		if err := accountsService.DeleteForID(ctx, missingID); err != nil {
			t.Errorf("unable to delete existing id [%s]", err)
		}

		status, account, _ := testServer.CallGetAccount(int(missingID))
		switch {
		case status != http.StatusNotFound:
			t.Errorf("expected status %d got %d", http.StatusNotFound, status)
		case account != nil:
			t.Errorf("unexpected response %v", *account)
		}
	})

}

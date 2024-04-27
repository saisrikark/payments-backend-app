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

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()
	testServer := testutils.NewTestServer(t)
	testServer.Start(ctx)
	defer testServer.Stop(ctx)
	testServer.WaitForRunnning(t)

	type TestData struct {
		description        string
		req                *server.CreateAccountRequest
		expectedStatusCode int
	}

	tests := []TestData{
		{
			description: "Expect successful account creation",
			req: &server.CreateAccountRequest{
				DocumentNumber: testutils.GenerateRandomNumber(10),
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			description:        "Bad Request empty json",
			req:                &server.CreateAccountRequest{},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description: "Bad Request empty document number",
			req: &server.CreateAccountRequest{
				DocumentNumber: "",
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			description:        "Bad Request empty body",
			req:                nil,
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	checkResponse := func(req server.CreateAccountRequest, resp server.GetAccountResponse) error {
		switch {
		case resp.AccountID == 0:
			return fmt.Errorf("response account id is 0")
		case resp.DocumentNumber != req.DocumentNumber:
			return fmt.Errorf("unexpected document id got %s expected %s", resp.DocumentNumber, req.DocumentNumber)
		default:
			return nil
		}
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {

			status, resp, err := testServer.CallCreateAccount(test.req)
			if err != nil {
				t.Errorf("create request failed [%s]", err.Error())
			}

			if status != test.expectedStatusCode {
				t.Errorf("unexpected status code expected %d got %d", http.StatusCreated, status)
			}

			if test.req != nil && resp != nil {
				if err := checkResponse(*test.req, *resp); err != nil {
					t.Errorf("fail response check [%s]", err.Error())
				}

				dbaccount, err := testServer.AccountsService.GetForID(ctx, resp.AccountID)
				if err != nil {
					t.Errorf("not able to fetch account [%s]", err.Error())
				}

				if dbaccount.AccountID != resp.AccountID {
					t.Errorf("account id's do not match %d %d", dbaccount.AccountID, resp.AccountID)
				}
			}

		})
	}

	t.Run("Duplicate account creation should fail", func(t *testing.T) {

		req := &server.CreateAccountRequest{
			DocumentNumber: testutils.GenerateRandomNumber(10),
		}

		_, err := testServer.AccountsService.Create(ctx, models.Account{DocumentNumber: req.DocumentNumber})
		if err != nil {
			t.Errorf("could not create account [%s]", err.Error())
		}

		status, _, err := testServer.CallCreateAccount(req)
		if err != nil {
			t.Errorf("create request failed [%s]", err.Error())
		}

		if status != http.StatusConflict {
			t.Errorf("unexpected status code expected %d got %d", http.StatusConflict, status)
		}
	})

}

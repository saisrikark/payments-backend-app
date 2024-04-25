package server

import (
	"context"
	"payments-backend-app/pkg/server"
	"payments-backend-app/test/testutils"
	"testing"
)

var (
	serverPort = 8080
)

func TestCreateAccount(t *testing.T) {
	ctx := context.Background()

	testServer := testutils.NewTestServer(t, testutils.WithPort(serverPort))
	testServer.Start(ctx)
	defer testServer.Stop(ctx)
	testServer.WaitForRunnning(t)

	status, err := testutils.CallCreateAccount(serverPort, server.CreateAccountRequest{})
	if err != nil {
		t.Errorf("create request failed [%s]", err.Error())
	}

	if status != 201 {
		t.Errorf("unexpected status code expected %d got %d", 201, status)
	}

}

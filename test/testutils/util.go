package testutils

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"

	"payments-backend-app/builder"
	"payments-backend-app/pkg/models"
	"payments-backend-app/pkg/server"
)

type TestApp struct {
	port               int
	accountsService    models.AccountsService
	transactionService models.TransactionService
	runner             builder.Runner
	runningLock        *sync.Mutex
}

type TestDatabase struct {
	Instance testcontainers.Container
}

type Option func(*TestApp)

func WithPort(port int) Option {
	return func(ta *TestApp) {
		ta.port = port
	}
}

func WithAccountsService(accountsService models.AccountsService) Option {
	return func(ta *TestApp) {
		ta.accountsService = accountsService
	}
}

func WithTransactionService(transactionService models.TransactionService) Option {
	return func(ta *TestApp) {
		ta.transactionService = transactionService
	}
}

func NewTestServer(t *testing.T, opts ...Option) *TestApp {

	testApp := &TestApp{}
	for _, opt := range opts {
		opt(testApp)
	}

	testApp.runningLock = &sync.Mutex{}

	paymentsAppBuilder := builder.
		NewPaymentsAppBuilder().
		WithPaymentsServerAddr(fmt.Sprintf(":%d", testApp.port)).
		WithAccountsService(testApp.accountsService).
		WithTransactionService(testApp.transactionService).
		WithDatabaseAddr("localhost:5432").
		WithDatabaseName("payments-db").
		WithDatabaseUser("payments-user").
		WithDatabasePassword("payments-password").
		UseInsecureDatabaseConnection()

	paymentsAppRunner, err := paymentsAppBuilder.Build()
	require.NoError(t, err)

	testApp.runner = paymentsAppRunner

	return testApp
}

func (ta *TestApp) Start(ctx context.Context) error {
	go func() {
		ta.runner.Start(ctx)
	}()
	return nil
}

func TryLivenessRequest(port int) int {
	url := fmt.Sprintf("http://localhost:%d%s", port, server.LivenessExtension)

	resp, err := http.Get(url)
	if err != nil {
		return 0
	}

	return resp.StatusCode
}

func (ta *TestApp) WaitForRunnning(t *testing.T) bool {

	isRunning := false

	for i := 0; i < 5; i++ {
		if status := TryLivenessRequest(ta.port); status == 200 {
			isRunning = true
			break
		}
		time.Sleep(time.Second * 1)
	}

	if !isRunning {
		t.FailNow()
	}

	return isRunning
}

func (ta *TestApp) Stop(ctx context.Context) error {
	return ta.runner.Stop(ctx)
}

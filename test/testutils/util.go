package testutils

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/testcontainers/testcontainers-go"

	"payments-backend-app/builder"
	"payments-backend-app/pkg/models"
)

var (
	addr = ":8080"
)

func GetBaseUrl() string {
	return "http://localhost" + addr
}

type TestApp struct {
	baseUrl            string
	AccountsService    models.AccountsService
	TransactionService models.TransactionService
	runner             builder.Runner
}

type TestDatabase struct {
	Instance testcontainers.Container
}

type Option func(*TestApp)

func WithAccountsService(accountsService models.AccountsService) Option {
	return func(ta *TestApp) {
		ta.AccountsService = accountsService
	}
}

func WithTransactionService(transactionService models.TransactionService) Option {
	return func(ta *TestApp) {
		ta.TransactionService = transactionService
	}
}

func NewTestServer(t *testing.T, opts ...Option) *TestApp {

	testApp := &TestApp{}
	for _, opt := range opts {
		opt(testApp)
	}

	envConfig := builder.GetEnvConfig()

	paymentsAppBuilder := builder.
		NewPaymentsAppBuilder().
		WithPaymentsServerAddr(envConfig.PaymentsAppAddr).
		WithAccountsService(testApp.AccountsService).
		WithTransactionService(testApp.TransactionService).
		WithDatabaseAddr(envConfig.DatabaseAddr).
		WithDatabaseName(envConfig.DatabaseName).
		WithDatabaseUser(envConfig.DatabaseUser).
		WithDatabasePassword(envConfig.DatabasePassword)

	if envConfig.UseInsecureDatabase {
		paymentsAppBuilder = paymentsAppBuilder.UseInsecureDatabaseConnection()
	}

	paymentsAppRunner, err := paymentsAppBuilder.Build()
	require.NoError(t, err)

	testApp.AccountsService = paymentsAppBuilder.AccountsService
	testApp.TransactionService = paymentsAppBuilder.TransactionService

	testApp.baseUrl = "http://localhost" + envConfig.PaymentsAppAddr
	testApp.runner = paymentsAppRunner

	return testApp
}

func (ta *TestApp) Start(ctx context.Context) error {
	go func() {
		if err := ta.runner.Start(ctx); err != nil {
			panic(err)
		}
	}()
	return nil
}

func (ta *TestApp) TryLivenessRequest() int {
	resp, err := http.Get(ta.baseUrl + "/liveness")
	if err != nil {
		return 0
	}
	return resp.StatusCode
}

func (ta *TestApp) WaitForRunnning(t *testing.T) bool {

	isRunning := false

	for i := 0; i < 5; i++ {
		if status := ta.TryLivenessRequest(); status == 200 {
			isRunning = true
			break
		}
		time.Sleep(time.Second * 1)
	}

	if !isRunning {
		t.Errorf("server not up and running")
	}

	return isRunning
}

func (ta *TestApp) Stop(ctx context.Context) error {
	return ta.runner.Stop(ctx)
}

func GenerateRandomNumber(length int) string {
	return fmt.Sprintf("%d", GenerateRandomNumberInt(length))
}

func GenerateRandomNumberInt(length int) int {
	min := pow(10, length-1)
	max := pow(10, length) - 1
	return rand.Intn(max-min+1) + min
}

func pow(base, exp int) int {
	result := 1
	for exp > 0 {
		if exp%2 == 1 {
			result *= base
		}
		base *= base
		exp /= 2
	}
	return result
}

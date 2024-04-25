package builder

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"payments-backend-app/internal/models"
	"payments-backend-app/pkg/server"

	"github.com/julienschmidt/httprouter"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Builder interface {
	Build() (Runner, error)
}

type Runner interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Option func(*PaymentsAppBuilder)

type PaymentsAppBuilder struct {
	// database config
	databaseAddr                  string
	databaseName                  string
	databaseUser                  string
	databasePassword              string
	useInsecureDatabaseConnection bool

	// payments server config
	paymentsServerAddr string

	// utils
	logger *slog.Logger
}

func NewPaymentsAppBuilder() *PaymentsAppBuilder {
	pab := &PaymentsAppBuilder{}

	return pab
}

func (pab *PaymentsAppBuilder) WithDatabaseAddr(databaseAddr string) *PaymentsAppBuilder {
	pab.databaseAddr = databaseAddr
	return pab
}

func (pab *PaymentsAppBuilder) WithDatabaseName(databaseName string) *PaymentsAppBuilder {
	pab.databaseName = databaseName
	return pab
}

func (pab *PaymentsAppBuilder) WithDatabaseUser(databaseUser string) *PaymentsAppBuilder {
	pab.databaseUser = databaseUser
	return pab
}

func (pab *PaymentsAppBuilder) WithDatabasePassword(databasePassword string) *PaymentsAppBuilder {
	pab.databasePassword = databasePassword
	return pab
}

func (pab *PaymentsAppBuilder) UseInsecureDatabaseConnection() *PaymentsAppBuilder {
	pab.useInsecureDatabaseConnection = true
	return pab
}

func (pab *PaymentsAppBuilder) WithPaymentsServerAddr(addr string) *PaymentsAppBuilder {
	pab.paymentsServerAddr = addr
	return pab
}

func (pab *PaymentsAppBuilder) WithLogger(logger *slog.Logger) *PaymentsAppBuilder {
	pab.logger = logger
	return pab
}

func (pab *PaymentsAppBuilder) Build() (Runner, error) {

	par := &paymentsAppRunner{}

	sqldb := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithAddr(pab.databaseAddr),
			pgdriver.WithDatabase(pab.databaseName),
			pgdriver.WithUser(pab.databaseUser),
			pgdriver.WithPassword(pab.databasePassword),
			pgdriver.WithInsecure(pab.useInsecureDatabaseConnection)))

	db := bun.NewDB(sqldb, pgdialect.New())

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to connect to database %s", err.Error())
	}

	accountsService := models.NewAccountsService(db)
	transactionService := models.NewTransactionService(db)

	pah := server.NewPaymentsAppHandler(accountsService, transactionService, server.WithLogger(pab.logger))

	router := httprouter.New()
	router.PanicHandler = pah.PanicHandler

	router.GET("/liveness", pah.Liveness)
	router.GET("/readiness", pah.Readiness)
	router.POST("/accounts", pah.CreateAccount)
	router.GET("/accounts/:accountId", pah.GetAccount)
	router.POST("/transactions", pah.CreateTransaction)

	server := &http.Server{
		Addr:    pab.paymentsServerAddr,
		Handler: router,
	}

	par.server = server

	return par, nil
}

type paymentsAppRunner struct {
	server *http.Server
}

func (par *paymentsAppRunner) Start(_ context.Context) error {

	if err := par.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (par *paymentsAppRunner) Stop(_ context.Context) error {
	return par.server.Shutdown(context.Background())
}

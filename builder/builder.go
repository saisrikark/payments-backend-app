package builder

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"payments-backend-app/internal/migrate"
	imodels "payments-backend-app/internal/models"
	"payments-backend-app/pkg/models"

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
	db                            *bun.DB
	disableDatabase               bool
	databaseAddr                  string
	databaseName                  string
	databaseUser                  string
	databasePassword              string
	useInsecureDatabaseConnection bool

	// services
	as models.AccountsService
	ts models.TransactionService

	// payments server config
	paymentsServerAddr string

	// utils
	logger *slog.Logger
}

func NewPaymentsAppBuilder() *PaymentsAppBuilder {
	pab := &PaymentsAppBuilder{}

	return pab
}

func (pab *PaymentsAppBuilder) WithDatabase(db *bun.DB) *PaymentsAppBuilder {
	pab.db = db
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

func (pab *PaymentsAppBuilder) WithAccountsService(as models.AccountsService) *PaymentsAppBuilder {
	pab.as = as
	return pab
}

func (pab *PaymentsAppBuilder) WithTransactionService(as models.TransactionService) *PaymentsAppBuilder {
	pab.ts = as
	return pab
}

func (pab *PaymentsAppBuilder) DisableDatabase() *PaymentsAppBuilder {
	pab.disableDatabase = true
	return pab
}

func (pab *PaymentsAppBuilder) Build() (Runner, error) {

	par := &paymentsAppRunner{}

	if pab.db != nil {
		par.db = pab.db
	}

	if !pab.disableDatabase {
		sqldb := sql.OpenDB(
			pgdriver.NewConnector(
				pgdriver.WithAddr(pab.databaseAddr),
				pgdriver.WithDatabase(pab.databaseName),
				pgdriver.WithUser(pab.databaseUser),
				pgdriver.WithPassword(pab.databasePassword),
				pgdriver.WithInsecure(pab.useInsecureDatabaseConnection)))

		if par.db == nil {
			par.db = bun.NewDB(sqldb, pgdialect.New())
		}

		if err := par.db.Ping(); err != nil {
			return nil, fmt.Errorf("unable to connect to database %s", err.Error())
		}

		if err := migrate.Run(context.Background(), par.db); err != nil {
			return nil, fmt.Errorf("unable to migrate [%s]", err.Error())
		}
	}

	if pab.as == nil {
		pab.as = imodels.NewAccountsService(par.db)
	}

	if pab.ts == nil {
		pab.ts = imodels.NewTransactionService(par.db)
	}

	pah := server.NewPaymentsAppHandler(pab.as, pab.ts, server.WithLogger(pab.logger))

	router := httprouter.New()
	router.PanicHandler = pah.PanicHandler

	router.GET(server.LivenessExtension, pah.Liveness)
	router.GET(server.ReadinessExtension, pah.Readiness)
	router.POST(server.CreateAccountExtension, pah.CreateAccount)
	router.GET(server.GetAccountExtension, pah.GetAccount)
	router.POST(server.CreateTransactionExtension, pah.CreateTransaction)

	server := &http.Server{
		Addr:    pab.paymentsServerAddr,
		Handler: router,
	}

	par.server = server

	return par, nil
}

type paymentsAppRunner struct {
	db     *bun.DB
	server *http.Server
}

func (par *paymentsAppRunner) Start(_ context.Context) error {

	if err := par.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("unable to start server [%s]", err.Error())
	}

	return nil
}

func (par *paymentsAppRunner) Stop(_ context.Context) error {
	return par.server.Shutdown(context.Background())
}

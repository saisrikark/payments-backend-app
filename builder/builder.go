package builder

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	imodels "payments-backend-app/internal/models"
	"payments-backend-app/pkg/models"

	"payments-backend-app/pkg/server"

	"github.com/julienschmidt/httprouter"
	"github.com/uptrace/bun"
)

type Runner interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type Option func(*PaymentsAppBuilder)

type PaymentsAppBuilder struct {
	isBuilt bool

	// database config
	db                            *bun.DB
	disableDatabase               bool
	databaseAddr                  string
	databaseName                  string
	databaseUser                  string
	databasePassword              string
	useInsecureDatabaseConnection bool

	// services
	AccountsService    models.AccountsService
	TransactionService models.TransactionService

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
	pab.AccountsService = as
	return pab
}

func (pab *PaymentsAppBuilder) WithTransactionService(as models.TransactionService) *PaymentsAppBuilder {
	pab.TransactionService = as
	return pab
}

func (pab *PaymentsAppBuilder) DisableDatabase() *PaymentsAppBuilder {
	pab.disableDatabase = true
	return pab
}

func (pab *PaymentsAppBuilder) GetAccountsService() (models.AccountsService, error) {
	if !pab.isBuilt {
		return nil, fmt.Errorf("not built")
	}
	return pab.AccountsService, nil
}

func (pab *PaymentsAppBuilder) GetTransactionService() (models.TransactionService, error) {
	if !pab.isBuilt {
		return nil, fmt.Errorf("not built")
	}
	return pab.TransactionService, nil
}

func (pab *PaymentsAppBuilder) Build() (Runner, error) {

	par := &paymentsAppRunner{}

	if pab.db != nil {
		par.db = pab.db
	}

	if !pab.disableDatabase {
		var err error
		par.db, err = NewDatabase(
			pab.databaseAddr,
			pab.databaseName,
			pab.databaseUser,
			pab.databasePassword,
			pab.useInsecureDatabaseConnection)
		if err != nil {
			return nil, err
		}
	}

	if pab.AccountsService == nil {
		pab.AccountsService = imodels.NewAccountsService(par.db)
	}

	if pab.TransactionService == nil {
		pab.TransactionService = imodels.NewTransactionService(par.db)
	}

	pah := server.NewPaymentsAppHandler(pab.AccountsService, pab.TransactionService, server.WithLogger(pab.logger))

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
	pab.isBuilt = true

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

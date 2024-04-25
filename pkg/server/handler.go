package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"payments-backend-app/pkg/models"
	"sync/atomic"

	"github.com/julienschmidt/httprouter"
)

type ExtensionDetail struct {
	Method string
	httprouter.Handle
}

var (
	LivenessExtension          = "/liveness"
	ReadinessExtension         = "/readiness"
	CreateAccountExtension     = "/accounts"
	GetAccountExtension        = "/accounts/:accountId"
	CreateTransactionExtension = "/transactions"
)

type paymentsAppHandler struct {
	panicCount         atomic.Int64
	accountsService    models.AccountsService
	transactionService models.TransactionService
	logger             *slog.Logger
}

// Option for the payments app server
type Option func(*paymentsAppHandler)

// NewPaymentsAppServer is used to create a payments app server
func NewPaymentsAppHandler(
	accountsService models.AccountsService,
	transactionService models.TransactionService,
	opts ...Option) *paymentsAppHandler {

	pah := &paymentsAppHandler{
		panicCount:         atomic.Int64{},
		accountsService:    accountsService,
		transactionService: transactionService,
	}

	for _, opt := range opts {
		opt(pah)
	}

	if pah.logger == nil {
		handlerOptions := &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}
		pah.logger = slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
	}

	return pah
}

func WithLogger(logger *slog.Logger) Option {
	return func(pas *paymentsAppHandler) {
		pas.logger = logger
	}
}

// PanicHandler is used to recover when there is a crash serving a request
func (pah *paymentsAppHandler) PanicHandler(w http.ResponseWriter, r *http.Request, i interface{}) {
	pah.panicCount.Add(1)
	ctx := context.Background()
	pah.logger.ErrorContext(ctx, "Panic")
	recover()
}

// Liveness returns whether the service is up and running
func (pah *paymentsAppHandler) Liveness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	pah.logger.DebugContext(ctx, "Called")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}

// Readiness returns whether the service is ready to serve request
func (pah *paymentsAppHandler) Readiness(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	pah.logger.DebugContext(ctx, "Called")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK\n")
}

// CreateAccount is used to create an account given a document number
func (pah *paymentsAppHandler) CreateAccount(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	pah.logger.DebugContext(ctx, "Called")

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", string("{}"))
}

func (pah *paymentsAppHandler) GetAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ctx := context.Background()
	pah.logger.DebugContext(ctx, "Called")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "{}")
}

func (pah *paymentsAppHandler) CreateTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()
	pah.logger.DebugContext(ctx, "Called")

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "{}")
}

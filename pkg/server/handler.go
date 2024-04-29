package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"payments-backend-app/pkg/models"
	"strconv"
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

	ba, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := CreateAccountRequest{}
	if err := json.Unmarshal(ba, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ba, _ := json.Marshal(map[string]string{"msg": err.Error()})
		fmt.Fprintf(w, "%s", string(ba))
		return
	}

	account, err := pah.accountsService.Create(ctx, models.Account{DocumentNumber: req.DocumentNumber})
	if err != nil {
		switch {
		case errors.Is(err, models.DuplicateRecordErr):
			w.WriteHeader(http.StatusConflict)
			ba, _ := json.Marshal(map[string]string{"msg": err.Error()})
			fmt.Fprintf(w, "%s", string(ba))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	ba, err = json.Marshal(account)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pah.logger.ErrorContext(ctx, "marshalling error", "err", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", string(ba))
}

// GetAccount fetches an account for the provided account id
func (pah *paymentsAppHandler) GetAccount(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	ctx := context.Background()
	accountIdS := params.ByName("accountId")

	accountId, err := strconv.Atoi(accountIdS)
	if err != nil {
		pah.logger.DebugContext(ctx, "unable to parse account id", "accountIdS", accountIdS, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	account, err := pah.accountsService.GetForID(ctx, int64(accountId))
	if err != nil {
		switch {
		case errors.Is(err, models.NoRecordErr):
			w.WriteHeader(http.StatusNotFound)
			ba, _ := json.Marshal(map[string]string{"msg": err.Error()})
			fmt.Fprintf(w, "%s", string(ba))
		default:
			pah.logger.ErrorContext(ctx, "unable to fetch account for id", "accountID", accountId, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := GetAccountResponse{
		AccountID:      account.AccountID,
		DocumentNumber: account.DocumentNumber,
	}

	ba, err := json.Marshal(resp)
	if err != nil {
		pah.logger.ErrorContext(ctx, "unable to marshal account", "account", account, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", string(ba))
}

// CreateTransaction creates a transaction given an account id, operation type and amount
func (pah *paymentsAppHandler) CreateTransaction(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := context.Background()

	ba, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	req := CreateTransactionRequest{}
	if err := json.Unmarshal(ba, &req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ba, _ := json.Marshal(map[string]string{"msg": err.Error()})
		fmt.Fprintf(w, "%s", string(ba))
		return
	}

	transactionStatus, err := pah.transactionService.Create(ctx, models.Transaction{
		AccountID:       req.AccountID,
		OperationTypeID: req.OperationTypeID,
		Amount:          req.Amount,
	})
	if err != nil {
		switch {
		case errors.Is(err, models.NoRecordErr):
			w.WriteHeader(http.StatusNotFound)
			ba, _ := json.Marshal(map[string]string{"msg": err.Error()})
			fmt.Fprintf(w, "%s", string(ba))
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	resp := CreateTransactionResponse{
		TransactionID: transactionStatus.TransactionID,
		AccountID:     transactionStatus.AccountID,
	}

	ba, err = json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		pah.logger.ErrorContext(ctx, "marshalling error", "err", err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", string(ba))
}

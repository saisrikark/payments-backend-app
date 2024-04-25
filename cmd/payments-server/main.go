package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"payments-backend-app/builder"

	"github.com/spf13/viper"
)

var (
	DATABASE_ADDR_ENV          = "DATABASE_ADDR"
	DATABASE_NAME_ENV          = "DATABASE_NAME"
	DATABASE_USER_ENV          = "DATABASE_USER"
	DATABASE_PASSWORD_ENV      = "DATABASE_PASSWORD"
	DATABASE_WITH_INSECURE_ENV = "DATABASE_WITH_INSECURE"
	PAYMENTS_APP_ADDR_ENV      = "PAYMENTS_APP_ADDR"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	// set default env variables
	viper.SetDefault(DATABASE_ADDR_ENV, "localhost:5432")
	viper.SetDefault(DATABASE_NAME_ENV, "payments-db")
	viper.SetDefault(DATABASE_USER_ENV, "payments-user")
	viper.SetDefault(DATABASE_PASSWORD_ENV, "payments-password")
	viper.SetDefault(DATABASE_WITH_INSECURE_ENV, "true")
	viper.SetDefault(PAYMENTS_APP_ADDR_ENV, ":8080")

	// bind env variables
	viper.BindEnv(DATABASE_ADDR_ENV)
	viper.BindEnv(DATABASE_NAME_ENV)
	viper.BindEnv(DATABASE_USER_ENV)
	viper.BindEnv(DATABASE_PASSWORD_ENV)
	viper.BindEnv(DATABASE_WITH_INSECURE_ENV)
	viper.BindEnv(PAYMENTS_APP_ADDR_ENV)

	// fetch config from env variables
	databaseAddr := viper.GetString(DATABASE_ADDR_ENV)
	databaseName := viper.GetString(DATABASE_NAME_ENV)
	databaseUser := viper.GetString(DATABASE_USER_ENV)
	databasePassword := viper.GetString(DATABASE_PASSWORD_ENV)
	useInsecureDatabase := viper.GetBool(DATABASE_WITH_INSECURE_ENV)
	paymentsAppAddr := viper.GetString(PAYMENTS_APP_ADDR_ENV)

	logger.InfoContext(ctx, "configuration",
		"databaseAddr", databaseAddr,
		"databaseName", databaseName,
		"databaseUser", databaseUser,
		"databasePassword", databasePassword,
		"useInsecureDatabase", useInsecureDatabase,
		"paymentsAppAddr", paymentsAppAddr)

	// build the runner
	paymentsAppBuilder := builder.
		NewPaymentsAppBuilder().
		WithDatabaseAddr(databaseAddr).
		WithDatabaseName(databaseName).
		WithDatabaseUser(databaseUser).
		WithDatabasePassword(databasePassword).
		UseInsecureDatabaseConnection().
		WithPaymentsServerAddr(paymentsAppAddr)

	paymentsAppRunner, err := paymentsAppBuilder.Build()
	if err != nil {
		log.Fatalf("unable to build [%s]", err.Error())
	}

	if err := paymentsAppRunner.Start(ctx); err != nil {
		log.Fatalf("unable to start [%s]", err.Error())
	}
}

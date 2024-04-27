package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"payments-backend-app/builder"
)

func main() {
	ctx := context.Background()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	envConfig := builder.GetEnvConfig()

	logger.InfoContext(ctx, "configuration",
		"databaseAddr", envConfig.DatabaseAddr,
		"databaseName", envConfig.DatabaseName,
		"databaseUser", envConfig.DatabaseUser,
		"databasePassword", envConfig.DatabasePassword,
		"useInsecureDatabase", envConfig.UseInsecureDatabase,
		"paymentsAppAddr", envConfig.PaymentsAppAddr)

	// build the runner
	paymentsAppBuilder := builder.
		NewPaymentsAppBuilder().
		WithDatabaseAddr(envConfig.DatabaseAddr).
		WithDatabaseName(envConfig.DatabaseName).
		WithDatabaseUser(envConfig.DatabaseUser).
		WithDatabasePassword(envConfig.DatabasePassword).
		WithPaymentsServerAddr(envConfig.PaymentsAppAddr)

	if envConfig.UseInsecureDatabase {
		paymentsAppBuilder = paymentsAppBuilder.UseInsecureDatabaseConnection()
	}

	paymentsAppRunner, err := paymentsAppBuilder.Build()
	if err != nil {
		log.Fatalf("unable to build [%s]", err.Error())
	}

	if err := paymentsAppRunner.Start(ctx); err != nil {
		log.Fatalf("unable to start [%s]", err.Error())
	}
}

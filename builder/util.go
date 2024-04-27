package builder

import (
	"context"
	"database/sql"
	"fmt"
	"payments-backend-app/internal/migrate"

	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

var (
	DATABASE_ADDR_ENV          = "DATABASE_ADDR"
	DATABASE_NAME_ENV          = "DATABASE_NAME"
	DATABASE_USER_ENV          = "DATABASE_USER"
	DATABASE_PASSWORD_ENV      = "DATABASE_PASSWORD"
	DATABASE_WITH_INSECURE_ENV = "DATABASE_WITH_INSECURE"
	PAYMENTS_APP_ADDR_ENV      = "PAYMENTS_APP_ADDR"
)

type EnvConfig struct {
	DatabaseAddr        string
	DatabaseName        string
	DatabaseUser        string
	DatabasePassword    string
	UseInsecureDatabase bool
	PaymentsAppAddr     string
}

func GetEnvConfig() EnvConfig {

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

	envConfig := EnvConfig{
		DatabaseAddr:        databaseAddr,
		DatabaseName:        databaseName,
		DatabaseUser:        databaseUser,
		DatabasePassword:    databasePassword,
		UseInsecureDatabase: useInsecureDatabase,
		PaymentsAppAddr:     paymentsAppAddr,
	}

	return envConfig
}

func NewDatabase(databaseAddr string, databaseName string, databaseUser string, databasePassword string, insecure bool) (*bun.DB, error) {
	sqldb := sql.OpenDB(
		pgdriver.NewConnector(
			pgdriver.WithAddr(databaseAddr),
			pgdriver.WithDatabase(databaseName),
			pgdriver.WithUser(databaseUser),
			pgdriver.WithPassword(databasePassword),
			pgdriver.WithInsecure(insecure)))

	db := bun.NewDB(sqldb, pgdialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithEnabled(false),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("unable to connect to database %s", err.Error())
	}

	if err := migrate.Run(context.Background(), db); err != nil {
		return nil, fmt.Errorf("unable to migrate [%s]", err.Error())
	}

	return db, nil
}

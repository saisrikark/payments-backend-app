#!/bin/bash

export DATABASE_ADDR="localhost:5432"
export DATABASE_NAME="payments-db"
export DATABASE_USER="payments-user"
export DATABASE_PASSWORD="payments-password"
export DATABASE_WITH_INSECURE="true"
export PAYMENTS_APP_ADDR=":8080"

go test -count=1 ./... -v

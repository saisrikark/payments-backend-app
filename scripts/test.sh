#!/bin/bash

timeout=60
interval=1
count=0

docker compose -f docker-compose-test.yml up &

while ! nc -z localhost 5432 >/dev/null 2>&1; do
    sleep "$interval"
    count=$((count + interval))
    if [ $count -ge $timeout ]; then
        echo "Timeout: PostgreSQL not available after $timeout seconds."
        exit 1
    fi
done

echo "PostgreSQL is available."

go test ./... -v

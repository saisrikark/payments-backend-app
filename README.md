# payments-backend-app
Payments application implemented using Go

## Requirements
- a free port (preferably 8080) to serve the http endpoints
- a postgres database, here I've used one within the docker compose files hosted at port 5432
- have docker compose and docker installed to avoid running the application locally
- have go version version >= 1.21.1 to run tests and build the application locally
- install "make" to run the Makefile

## APIs List

1. **Create Account**
   - **Endpoint**: `http://localhost:8080/accounts`
   - **Example Request**:
     ```bash
        curl -X POST http://localhost:8080/accounts \
        -d '{
                "document_number": "12345"
            }'
     ```
   - **Sample Response**:
     ```json
        {
            "account_id": 4,
            "document_number": "12345"
        }
     ```

2. **Get Account API**
   - **Endpoint**: `http://localhost:8080/accounts/:accountId`
   - **Example Request**:
     ```bash
        curl http://localhost:8080/accounts/4
     ```
   - **Sample Response**:
     ```json
        {
            "account_id": 4,
            "document_number": "12345"
        }
     ```

3. **Create Transaction API**
   - **Endpoint**: `http://localhost:8080/transactions`
   - **Example Request**:
     ```bash
        curl -X POST http://localhost:8080/transactions \
        -d '{
                "account_id": 6,
                "operation_type_id": 4,
                "amount": 111.11
            }'
     ```
   - **Sample Response**:
     ```json
        {}
     ```

```
Please refer to the open api specification under swagger/* for further information
```

## Setup

### Using Docker

```
Use the script run.sh under scripts to setup the go application as well as a postgres database.
This uses the docker-compose.yml file.
Database is exposed at port 5432.

Payments application is exposed at 8080.
Change the ports if required.

To execute the script
chmod +x ./scripts/run.sh (make the script executable)
./scripts/run.sh (to run the script)
```

### Local Setup

```
Try to avoid this, use only as a last resort.
Set/change below environment variables as per your preference

export DATABASE_ADDR="localhost:5432"
export DATABASE_NAME="payments-db"
export DATABASE_USER="payments-user"
export DATABASE_PASSWORD="payments-password"
export DATABASE_WITH_INSECURE="true"
export PAYMENTS_APP_ADDR=":8080"

Install postgres and create the database, a user and give the password based on the environment variables set above.
Start postgres server.

Have go version >= 1.21.1 installed.
Ensure "make" command is installed.

Build the go application
make clean
make build

Application is stored under bin/
Start the application with above environment variables.
./bin/payments-server
```

## Testing

```
Tests spin up an instance of the payments backend server at a specified port.
All tests make http requests at the server.

To run the tests, start a docker instance at the port in docker-compose-test.yml
(change the database port if required)
chmod +x ./scripts/run_db.sh (make the script executable)
./scripts/run_db.sh (to start the database)

Configure the environment variables in test.sh if a change is required.
Remember to change the database port here is you've changed it in run_db.sh.

To execute the script
chmod +x ./scripts/test.sh (make the script executable)
./scripts/test.sh (to run the script)
```
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINPATH=bin

all: test build

build:
	$(GOBUILD) -o $(BINPATH)/payments-server cmd/payments-server/main.go

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -rf $(BINPATH)/*
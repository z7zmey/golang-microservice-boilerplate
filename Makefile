BASEPATH = $(shell pwd)
export PATH := $(BASEPATH)/bin:$(PATH)

# Basic go commands
GOCMD      = go
GOBUILD    = $(GOCMD) build
GOINSTALL  = $(GOCMD) install
GORUN      = $(GOCMD) run
GOCLEAN    = $(GOCMD) clean
GOTEST     = $(GOCMD) test
GOGET      = $(GOCMD) get
GOFMT      = $(GOCMD) fmt
GOGENERATE = $(GOCMD) generate

# Swagger
SWAGGER = swagger

BUILD_DIR = $(BASEPATH)
COVERAGE_DIR  = $(BUILD_DIR)/coverage
SUBCOV_DIR    = $(COVERAGE_DIR)/packages

# Colors
GREEN_COLOR   = "\033[0;32m"
PURPLE_COLOR  = "\033[0;35m"
DEFAULT_COLOR = "\033[m"

all: clean fmt swagger build test lint

help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binary.'
	@echo '    test               Run unit tests.'
	@echo '    lint               Run all linters including vet and gosec and others'
	@echo '    fmt                Run gofmt on package sources.'
	@echo '    coverage           Report code tests coverage.'
	@echo '    build              Compile packages and dependencies.'
	@echo '    version            Print Go version.'
	@echo '    swagger            Generate swagger models and server'
	@echo ''
	@echo 'Targets run by default are: clean fmt swagger build test lint.'
	@echo ''

clean:
	@echo $(GREEN_COLOR)[clean]$(DEFAULT_COLOR)
	@$(GOCLEAN)

lint:
	@echo $(GREEN_COLOR)[lint]$(DEFAULT_COLOR)
	@$(GORUN) ./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint/main.go run \
	--no-config \
	--disable=errcheck \
	--enable=gosec \
	--enable=prealloc \
	./...

test:
	@echo $(GREEN_COLOR)[test]$(DEFAULT_COLOR)
	@$(GOTEST) -race ./...

fmt:
	@echo $(GREEN_COLOR)[format]$(DEFAULT_COLOR)
	@$(GOFMT) ./...

build:
	@echo $(GREEN_COLOR)[build]$(DEFAULT_COLOR)
	@$(GOBUILD) --tags static -o ./bin ./cmd/...

version:
	@echo $(GREEN_COLOR)[version]$(DEFAULT_COLOR)
	@$(GOCMD) version

swagger-build-binary:
ifeq ("$(wildcard ./bin/$(SWAGGER))","")
	@echo $(GREEN_COLOR)[swagger-build-binary]$(DEFAULT_COLOR)
	@$(GOBUILD) -o ./bin/$(SWAGGER) ./vendor/github.com/go-swagger/go-swagger/cmd/swagger
endif

swagger-clean:
	@echo $(GREEN_COLOR)[swagger cleanup]$(DEFAULT_COLOR)
	@rm -rf $(BASEPATH)/internal/transport/rest/server/models
	@rm -rf $(BASEPATH)/internal/transport/rest/server/restapi

# swagger generate server -f ./api/spec.yaml --exclude-main --flag-strategy=pflag --default-scheme=http --target=./internal/server
swagger: swagger-build-binary swagger-clean
	@echo $(GREEN_COLOR)[swagger]$(DEFAULT_COLOR)
	./bin/$(SWAGGER) generate server \
	   -P models.Principal \
	   -f ./api/swagger.yaml \
	   --exclude-main \
	   --default-scheme=http \
	   --target=$(BASEPATH)/internal/transport/rest/server

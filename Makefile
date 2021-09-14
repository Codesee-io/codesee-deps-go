BIN_DIR          ?= ./bin
COVERAGE_PROFILE ?= coverage.out
TEST_FLAGS       ?=

default: build

.PHONY: build
build:
	@echo "---> Building"
	CGO_ENABLED=0 go build -o $(BIN_DIR)/codesee-deps-go -installsuffix cgo ./cmd/deps

.PHONY: html
html:
	@echo "---> Generating HTML coverage report"
	go tool cover -html $(COVERAGE_PROFILE)

.PHONY: install
install:
	@echo "---> Installing dependencies"
	go mod download

.PHONY: test
test:
	@echo "---> Testing"
	go test -race ./pkg/... -coverprofile $(COVERAGE_PROFILE) $(TEST_FLAGS)

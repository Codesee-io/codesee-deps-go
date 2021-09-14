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

.PHONY: release
release: $(BIN_DIR)/goreleaser
	@echo "---> Creating new release"
ifndef TAG
	$(error TAG must be specified)
endif
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN must be specified. Generate one here: https://github.com/settings/tokens/new)
endif
	sed -i "" "s/version-.*-green/version-$(TAG)-green/" README.md
	git add README.md
	git commit -m $(TAG)
	git tag $(TAG)
	git push origin main --tags
	$(BIN_DIR)/goreleaser release --rm-dist

$(BIN_DIR)/goreleaser:
	@echo "---> Installing goreleaser"
	curl -sfL https://install.goreleaser.com/github.com/goreleaser/goreleaser.sh | sh

.PHONY: test
test:
	@echo "---> Testing"
	go test -race ./pkg/... -coverprofile $(COVERAGE_PROFILE) $(TEST_FLAGS)

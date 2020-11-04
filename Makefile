.PHONY: default fix lint test install build

# --- dev --- #

default: fix lint test

fix:
	@ echo ">> fixing source code"
	@ gofmt -s -l -w cmd
	@ go mod tidy
	@ echo ">> done"

lint:
	@ echo ">> running linter"
	@ golangci-lint run --skip-dirs examples
	@ echo ">> done"

test:
	@ echo ">> running tests"
	@ go test --count=1 `go list ./... | grep -v examples`
	@ echo ">> done"


version=snapshot-$(USER)-$(shell git rev-parse --short HEAD)

install:
	@ echo ">> installing cli (dev)"
	@ go install -ldflags "-X main.version=$(version)" ./cmd/mono
	@ echo ">> done"

# --- release --- #

build:
	@ echo ">> building cli binaries"
	@ echo "todo :P"
	@ echo ">> done"


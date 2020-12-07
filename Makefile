.PHONY: default fix lint test install build

# --- dev --- #

default: fix lint test

fix:
	@ echo ">> fixing source code"
	@ gofmt -s -l -w .
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

empty-dev-bucket:
	@ echo ">> emptying dev artifact bucket"
	@ aws s3 rm $(MONO_ARTIFACT_BUCKET) --recursive
	@ echo ">> done"

LOCAL_INSTALL_VERSION=snapshot-$(USER)-$(shell git rev-parse --short HEAD)

install: default
	@ echo ">> installing cli (dev)"
	@ go install -ldflags "-X main.version=$(LOCAL_INSTALL_VERSION)" ./cmd/mono
	@ echo ">> done"

# --- release --- #

build:
	@ echo ">> building cli binaries"
	@ ./release/build.sh
	@ echo ">> done"

DOCS_OUTFILE=tmp.md

docs: install empty-dev-bucket
	@ echo ">> generating docs"
	@ OUTFILE=$(DOCS_OUTFILE) ./release/docs.sh
	@ echo ">> done"


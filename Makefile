VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

.PHONY: build install test clean release

build:
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o bin/appre ./cmd/appre/

install:
	CGO_ENABLED=1 go install -ldflags "$(LDFLAGS)" ./cmd/appre/

test:
	CGO_ENABLED=1 go test ./...

clean:
	rm -rf bin/ dist/

release: clean
	mkdir -p dist
	CGO_ENABLED=1 GOOS=darwin GOARCH=arm64 go build -ldflags "$(LDFLAGS)" -o dist/appre-darwin-arm64 ./cmd/appre/
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -ldflags "$(LDFLAGS)" -o dist/appre-darwin-amd64 ./cmd/appre/

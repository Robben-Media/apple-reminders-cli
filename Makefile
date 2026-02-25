VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)

.PHONY: build install test clean

build:
	CGO_ENABLED=1 go build -ldflags "$(LDFLAGS)" -o bin/appre ./cmd/appre/

install:
	CGO_ENABLED=1 go install -ldflags "$(LDFLAGS)" ./cmd/appre/

test:
	go test ./...

clean:
	rm -rf bin/

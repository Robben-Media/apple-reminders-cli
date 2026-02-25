# AGENTS.md — apple-reminders-cli

## Build
```bash
make build        # builds bin/appre (CGO_ENABLED=1 required)
make install      # installs to $GOPATH/bin
make test         # runs tests
make clean        # removes bin/
```

## Architecture
- `cmd/appre/` — Cobra CLI entry point and command definitions
- `internal/service/` — Business logic wrapping go-eventkit
- `internal/rpc/` — JSON-RPC 2.0 server (stdin/stdout)
- `internal/model/` — Internal data types

## Conventions
- All output is JSON to stdout, errors as JSON to stderr
- CGO_ENABLED=1 always (EventKit requires cgo)
- macOS only (EventKit is Apple-only)

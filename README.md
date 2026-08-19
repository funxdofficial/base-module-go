# base-module-go

Go **service generator CLI** that scaffolds new Go microservices with a consistent tier architecture. This repository is the generator itself — it is **not** a runnable application.

## Install

```bash
go install github.com/funxdofficial/base-module-go@latest
```

Binary will be available as **`base-module-go`** in your `$GOPATH/bin` (or `$HOME/go/bin`).

From source:

```bash
git clone <repo-url>
cd base-module-go
make install
# or
make build && ./bin/base-module-go --help
```

## Generate a service

### Via CLI (after install)

```bash
# REST API (Echo CRUD, no scheduler)
base-module-go --pkg=order-service --type=rest

# Consumer / cron (no HTTP)
base-module-go --pkg=sync-worker --type=consumer
base-module-go --pkg=sync-worker --type=cons   # alias
```

### Via Makefile (from this repo)

```bash
make fun-build PKG=order-service TYPE=rest
make fun-build PKG=sync-worker TYPE=consumer
make fun-build PKG=billing-service TYPE=rest OUTPUT=/tmp/billing-service
```

| Variable | Default | Description |
|----------|---------|-------------|
| `PKG` | *(required)* | Service name / output folder name |
| `TYPE` | `rest` | `rest`, `consumer`, or `cons` |
| `OUTPUT` | `./<PKG>` | Custom output directory |

### More CLI examples

```bash
# Default type is rest
base-module-go --pkg=payment-service

# Custom output directory
base-module-go --pkg=billing-service --type=rest --output=/tmp/billing-service

# Using new subcommand
base-module-go new order-service --type=rest
base-module-go new sync-worker --type=consumer

# From repo without install
go run . --pkg=my-service --type=rest
```

Generated Go module path: `github.com/funxdofficial/<service-name>`.

Service name rules: lowercase alphanumeric with hyphens, at least 3 characters (e.g. `my-service`).

## When to use `rest` vs `consumer`

| Type | Flag | Use when you need… |
|------|------|--------------------|
| **`rest`** (default) | `--type=rest` | HTTP APIs, CRUD endpoints, webhooks, anything served over Echo REST |
| **`consumer`** | `--type=consumer` or `--type=cons` | Background cron jobs, batch sync, periodic cleanup, scheduled reports — no HTTP server |

## CLI reference

| Command | Description |
|---------|-------------|
| `base-module-go --pkg=<name> [--type=rest\|consumer\|cons] [--output=<dir>]` | Generate a service (default action) |
| `base-module-go new <name> [--type=rest\|consumer\|cons] [--output=<dir>]` | Generate a service (explicit subcommand) |
| `make fun-build PKG=<name> [TYPE=rest] [OUTPUT=<dir>]` | Generate via Makefile (builds local binary first) |
| `base-module-go --help` | Show usage |
| `base-module-go --version` | Show version |

## After generate

```bash
cd order-service
cp .env.example .env
go mod tidy
go test ./...
go run .
```

For local development with sibling module repos, add to generated `go.mod`:

```go
replace (
    github.com/funxdofficial/golang-module-scheduler => ../golang-module-scheduler
    github.com/funxdofficial/golang-module-syslog => ../golang-module-syslog
)
```

### REST dummy CRUD

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/items
curl -X POST http://localhost:8080/api/v1/items \
  -H 'Content-Type: application/json' \
  -d '{"name":"Item A","description":"dummy"}'
```

## Generated project structure

Templates live under `funbuild/tpls/rest/` and `funbuild/tpls/consumer/`, embedded at build time.

### REST (`--type=rest`)

```
<service-name>/
├── config/                         # Tier 1 — env-only config
├── model/                          # Tier 1
│   ├── request/                    # Tier 2
│   └── response/                   # Tier 2
├── route/                          # Tier 1 — REST routes
├── service/                        # Tier 1 — HTTP transport
│   ├── controller/                 # Tier 2
│   ├── usecase/                    # Tier 2
│   └── repository/                 # Tier 2
├── module/                         # syslog bootstrap
└── container/                      # DI wiring
```

### Consumer (`--type=consumer`)

```
<service-name>/
├── config/
├── model/
│   ├── request/
│   └── response/
├── service/                        # Tier 1 — cron transport
│   ├── scheduler.go
│   ├── usecase/
│   └── repository/
├── module/
└── container/
```

## Tier architecture

| Tier | Packages | Role |
|------|----------|------|
| **Tier 1** | `config`, `model`, `service`, `module`, `container` | Env config, domain, transport, DI |
| **Tier 1 (REST only)** | `route` | Echo route registration |
| **Tier 2** | `model/request`, `model/response`, `service/repository`, `service/usecase` | DTOs, data, business logic |
| **Tier 2 (REST only)** | `service/controller` | HTTP orchestration |

**REST wiring:** `repository → usecase → controller → route → service`

**Consumer wiring:** `repository → usecase → scheduler → service`

## Interface file naming

| Layer | Interface file | Implementation file |
|-------|----------------|---------------------|
| Repository | `repository/repository.go` | `repository/repository_<entity>.go` |
| Usecase | `usecase/usecase.go` | `usecase/usecase_<entity>.go` |
| Controller (REST) | `controller/controller.go` | `controller/controller_<entity>.go` |

## Generator development

```bash
make tidy
make test
make build                    # outputs bin/base-module-go
make install                  # go install .
make fun-build PKG=demo-service TYPE=rest
make generate-example         # REST demo at /tmp/demo-service
make generate-consumer-example
make integration-test
```

## Reference

- Scheduler: [`github.com/funxdofficial/golang-module-scheduler`](../golang-module-scheduler)
- Syslog: [`github.com/funxdofficial/golang-module-syslog`](../golang-module-syslog)

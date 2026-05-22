# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

kk-scheduler is a job scheduling system that uses cron and gRPC. Users register services (gRPC servers implementing `KKScheduleTriggerServer`), add jobs with cron specs, and the scheduler triggers jobs at their scheduled times via gRPC calls. It also provides a Web UI for managing services, jobs, task executions, and logs.

## Repository Layout

- `go/` — Go backend (module `github.com/cruvie/kk-scheduler/go`)
- `ui/` — Nuxt 4 frontend (SPA, SSR disabled)
- `deploy-docker/` — Docker deployment files

All backend commands must be run from the `go/` directory.

## Commands

### Backend (from `go/` directory)
```bash
# Format code
just fmt

# Run linter
just lint

# Generate protobuf code (requires buf)
just proto-gen

# Check for dead code
just deadcode

# Check for vulnerabilities
just govulncheck

# Run with race detector
just race_detector

# Run nilaway analysis
just nilaway

# Auto-fix deprecated APIs
just fix

# Update all dependencies
just update-dep

# Clean build cache
just cache-clean

# Run all tests (requires a running server on port 8666)
go test ./...

# Run a single test
go test ./internal/scheduler_test/ -run TestJobList

# Run the server
go run ./internal/main/main.go
```

### Frontend (from `ui/` directory)
```bash
npm install
npm run dev      # Development server on http://localhost:3000
npm run build    # Production build → go/public/
```

## Architecture

### Frontend (Nuxt 4 + Vue 3)

The UI is an SPA (`ssr: false`) using Nuxt 4, Vue 3, and **Nuxt UI v4** (`@nuxt/ui`). It communicates with the backend via **gRPC-Web** using `@connectrpc/connect-web`.

**RPC client setup** (`ui/app/utils/api/`):
- `api.ts` — creates a `GrpcWebTransport` pointing at `http://localhost:8667` (the gRPC-Web HTTP server)
- `client.ts` — creates the typed `KKSchedule` client from the transport
- `interceptor.ts` — adds `JwtAuthKey` and `TraceId` headers to every request

**Generated types** live in `ui/gen/` (protobuf TypeScript via `@bufbuild/protobuf` + `protoc-gen-es`). Each page constructs requests with `create(Schema, data)` from `@bufbuild/protobuf` and calls typed methods on the client.

**Component conventions**: All pages use Nuxt UI v4 components (`UCard`, `UTable`, `UModal`, `USlideover`, `UButton`, `UBadge`, `UForm`). Forms live in `ui/app/components/` (e.g., `JobForm.vue`, `ServiceForm.vue`). The layout in `layouts/default.vue` uses `USidebar` + `UNavigationMenu`.

### Server Startup

`main.go` starts 4 services via `kk_server.KKServer`:
1. **kk-scheduler** — cron scheduler (`scheduler.NewScheduleServer()`)
2. **kk-scheduler-grpc** — gRPC server on `GrpcPort` (default 8666)
3. **kk-scheduler-http** — gRPC-Web HTTP server on `HttpPort` (default 8667), wraps the gRPC server via `grpcweb.WrapServer`
4. **kk-scheduler-web** — static file server on `WebPort` (default 8668), serves `go/public/` (the built UI)

### gRPC Request Flow

Each RPC method follows a strict pattern through these layers:

1. **`api_impl/unary.go`** — `server` struct implements the gRPC service interface. Each method calls `kk_grpc.GrpcHandler(ctx, input, constructor)` which orchestrates:
   - Instantiate the handler via the constructor function
   - Call `CheckInput()` (validates request fields using proto field behaviors)
   - Call `Handler()` which delegates to `Service()` for business logic

2. **`api_handlers/<domain>/`** — Each RPC has 3 files per handler group:
   - `api.go` — Handler struct embedding `*kk_grpc.DefaultApi[InputType]`
   - `check.go` — `CheckInput()` method for request validation
   - `handler.go` — `Handler()` method returning the output protobuf
   - `service.go` — `Service()` method with business logic (reads/writes via `scheduler.GClient`)

3. **`scheduler.GClient`** — Global singleton managing cron jobs and storage. All job/service/task operations go through this client.

4. **`store_driver.StoreDriver`** — Interface for persistence. Methods cover Job CRUD, Service CRUD, and TaskExecution CRUD + log appending.

### Interceptor Chain

gRPC unary calls pass through (in order):
1. `interceptor.UnaryInit` — resolves the method descriptor from proto file descriptor hub, extracts auth requirements + client IP into context
2. `interceptor.UnaryLogging` — structured slog logging
3. `recovery.UnaryServerInterceptor` — panic recovery

The auth interceptor (`interceptor.UnaryAuth`) is configured separately per-deployment and checks `InterceptorAuth` annotations from proto method options (JWT or InternalOnly token).

### Proto Field Validation

`common_go.CheckFields()` reads custom field options (`FieldBehavior`) from proto definitions:
- `REQUIRED` — field must be set
- `UUID7` — field value must be a valid UUIDv7

Authorization requirements per RPC are also declared via proto method options (`InterceptorAuthList`), read by `common_go.MethodDescGetInterceptorAuth()`.

### Trigger Mechanism

When a cron job fires, `triggerFunc()` dials the registered service's gRPC target (using insecure credentials, optionally with `Authority` header for auth tokens), creates a `KKScheduleTriggerClient`, and calls `Trigger(funcName, jobId)`.

### TaskExecutor SDK

`kk_scheduler/task_executor.go` provides a client-side framework for services to build multi-step task workflows:

```go
executor := kk_scheduler.NewTaskExecutor(
    kk_scheduler.WithSchedulerClient(client),
    kk_scheduler.WithJobId(jobId),
)
executor.AddStep("step1", func(ctl *StepCtl) error { ... }, nil)
executor.AddStep("step2", func(ctl *StepCtl) error { ... }, fallbackFunc)
executor.Run(ctx)
```

Each step's handler receives a `*StepCtl` for logging (`ctl.Log(err, msg)`). Steps run sequentially; if a step fails and has a fallback, the fallback runs. The executor auto-creates execution records, appends logs, and updates status (COMPLETED/FAILED) on the scheduler.

### Sentinel Errors

`kk_scheduler/Error.go` defines package-level sentinel errors used across the codebase: `ErrJobNotFount`, `ErrServiceNotFount`, `ErrServiceHasJob`, `ErrStopTask`, and others. All job/service lookups return these typed errors for consistent error handling.

### The `kk_scheduler/` Package

This package serves dual purpose: it contains both **generated protobuf code** (`*.pb.go`, `*_grpc.pb.go`) and **hand-written SDK code** (`task_executor.go`, `task_executor_options.go`, `Error.go`, `Util.go`, `Base.go`). External services import this package to get proto types and the TaskExecutor framework. The proto generation (`just proto-gen`) only regenerates `.pb.go` files — hand-written files are safe.

### Models & Database

`internal/models/` contains GORM model structs (`Job`, `Service`, `TaskExecution`). The `query/` subdirectory contains generated query code (via `gorm.io/gen`). Tables are auto-created on startup via `kk_pg.CreateTables()`.

## Testing

Tests are integration tests that connect to a running gRPC server on port 8666. Start the server first (`go run ./internal/main/main.go`), then run tests with `go test ./...`.

Test files live alongside their packages (`*_test.go`) except for `internal/scheduler_test/` which contains end-to-end tests exercising the full gRPC API via `grpc.NewClient`. These use `testify/assert` and insecure credentials.

The test in `internal/scheduler_test/client_server_test.go` implements a stub `KKScheduleTriggerServer` — use this as a reference when building services that receive scheduler triggers.

## Configuration

`go/config.toml`:
```toml
GrpcPort = 8666
HttpPort = 8667
WebPort = 8668

[Store]
Choose="PG"
[Store.PG.DSN]
Host = "127.0.0.1"
Port = 5432
User = "postgres"
Password = "testpg"
DBName = "kk_scheduler"
Schema = ""
SSLMode = "disable"
TimeZone = "UTC"
```

Set `KK_Schedule` env var for environment mode (handled by `kk_go_kit`).

## Storage Backends

Default is Postgres. Etcd also exists as a legacy option (`store_driver/etcd.go`). To add a new backend:
1. Create `store_xxxx.go` in `internal/store_driver/` implementing `StoreDriver`
2. Register it in `NewStoreDriver()` in `driver.go`
3. Add config fields in `g_config/config.go` and `config.toml`

## Key Dependencies

- `gitee.com/cruvie/kk_go_kit` — Shared library: gRPC utilities (`kk_grpc.GrpcHandler`, `DefaultApi`), logging, server lifecycle (`kk_server`), PostgreSQL helpers (`kk_pg`), JWT
- `github.com/robfig/cron/v3` — Cron scheduling engine
- `github.com/BurntSushi/toml` — TOML config parsing
- `gorm.io/gorm` + `gorm.io/gen` — ORM and query code generator
- `github.com/improbable-eng/grpc-web` — gRPC-Web proxy for browser clients
- `go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc` — OpenTelemetry tracing
- `google.golang.org/grpc` + `google.golang.org/protobuf` — gRPC and protobuf

### Frontend
- Nuxt 4 (SSR disabled, SPA mode)
- `@nuxt/ui` v4 for components
- `@connectrpc/connect-web` for gRPC-Web client communication

## Protobuf

- Proto files: `go/kk_scheduler/*.proto`
- Generated code: `*_grpc.pb.go` (gRPC service stubs), `*.pb.go` (messages)
- Uses buf with `go/buf.gen.yaml`
- Edition 2023 with `API_OPAQUE` level
- Custom proto extensions define field validation rules (`FieldBehavior`) and auth requirements (`InterceptorAuth`) per RPC method

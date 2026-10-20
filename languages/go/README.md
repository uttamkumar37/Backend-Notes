# Learn backend with Go

One reading order. Detailed notes in short sections. One project that grows with you.

```text
go/
├── README.md                 Start here; learning checklist
├── 01-go-basics.md            Language, standard library, Java → Go
├── 02-http-api.md             HTTP, clients, auth, contracts, Gin
├── 03-databases-cache.md      SQL, transactions, pools, Valkey, outbox
├── 04-concurrency.md          Channels, locks, workers, cancellation
├── 05-production.md           Testing, telemetry, security, deployment
├── SOURCES.md                 Official references and research choices
├── practice/                  Two small programs; one Go module
│   ├── basics/
│   └── packages/greeting/
├── project/                   Investor API; one flat Go module
└── check.sh                   Check both modules
```

## Follow this order

| Step | Read | Do before moving on |
| --- | --- | --- |
| 1 | [Go basics](01-go-basics.md) | Change a function, import a package, explain an error |
| 2 | [HTTP APIs](02-http-api.md) | Run the project and complete create → get → update → delete |
| 3 | [Databases and caching](03-databases-cache.md) | Implement persistence, then prove cache invalidation |
| 4 | [Concurrency](04-concurrency.md) | Explain ownership; cancel work; run the race detector |
| 5 | [Production habits](05-production.md) | Test a failure, trace startup, explain deployment and recovery |

Each guide starts with an overview and continues into detailed reference sections. Read one section, try its example, and return for another topic later; completing a whole file is not a daily target.

## Find a topic quickly

| Need | Jump to |
| --- | --- |
| Understand slice/map ownership | [Slices](01-go-basics.md#arrays-slices-and-copying), [maps](01-go-basics.md#maps-and-sets) |
| Understand pointers and interface nil | [Receivers](01-go-basics.md#pointers-receivers-and-mutation), [typed nil](01-go-basics.md#the-typed-nil-trap) |
| Use errors and generics | [Errors](01-go-basics.md#error-design), [generics](01-go-basics.md#generics-and-constraints) |
| Read modules and packages | [Tooling](01-go-basics.md#module-commands-and-tooling), [package design](01-go-basics.md#package-design-and-initialization) |
| Build an HTTP boundary | [Layers](02-http-api.md#handler-service-and-repository), [clients](02-http-api.md#calling-another-http-service) |
| Review API access and limits | [Authorization](02-http-api.md#authentication-and-authorization), [overload](02-http-api.md#rate-limiting-and-overload) |
| Work with PostgreSQL | [Query example](03-databases-cache.md#a-context-aware-query), [transactions](03-databases-cache.md#transaction-lifecycle), [pools](03-databases-cache.md#pool-sizing-and-saturation) |
| Understand cache correctness | [Keys](03-databases-cache.md#cache-keys-and-serialization), [stale data](03-databases-cache.md#stampedes-negative-caching-and-stale-data) |
| Stop concurrent work safely | [Context](04-concurrency.md#context-trees-and-causes), [worker pool](04-concurrency.md#a-bounded-worker-pool) |
| Investigate a hang/race | [Failure types](04-concurrency.md#races-deadlocks-and-leaks) |
| Test and measure | [Tests](05-production.md#table-driven-tests), [benchmarks](05-production.md#benchmarks-and-profiles) |
| Operate a service | [Logs](05-production.md#structured-logging-example), [containers](05-production.md#container-build-and-runtime), [troubleshooting](05-production.md#production-style-troubleshooting) |

**Code available:** two practice programs and an in-memory HTTP project.
Gin, PostgreSQL, Valkey, Vault, containers, and Kubernetes are explained as next steps; their integrations are exercises, not supplied implementations.

The notes also contain complete standalone Go programs and test examples. Save a complete program as a temporary `main.go` and run `go run main.go`. Package/test examples belong in a separate empty scratch module initialized with `go mod init example.com/notes-scratch`; save tests as `something_test.go` and run `go test`. Do not paste several `main` functions into the project package. Database examples explicitly state their runtime prerequisites.

## Start in five minutes

Install Go using the [official instructions](https://go.dev/doc/install). The module files need Go **1.27 or newer**. Choose a supported patched release; check [release history](https://go.dev/doc/devel/release) when installing. One editor and a terminal are enough to begin.

From the repository root:

```bash
cd languages/go/practice
go run ./basics
go run ./packages
go test ./...
```

The programs print `Learning Go for backend engineering.` plus a next-step line, then `Welcome, Asha.`. Tests should report `ok` for three packages. Next, open [the project guide](project/README.md).

## Learn a little, then prove it

For each 45–90 minute session: read one section, change one behavior, run its tests, then explain what failed and why. Next session, rebuild that small change from memory.

**Your checkpoint:** you previously reported the first program complete. Package/module practice remains pending. This restructuring does not mark new skills complete.

- [x] First program — completed by your earlier report
- [ ] Packages and modules — trace the import; change the greeting
- [ ] Values, slices, structs, pointers, interfaces, errors
- [ ] HTTP CRUD and input validation
- [ ] PostgreSQL and transactions
- [ ] Cache hits, invalidation, and fallback
- [ ] Cancellation, synchronization, and shutdown
- [ ] Tests, debugging, and deployment

Mark a topic complete when you can **build it, explain it, and repair a failure** without copying the solution.

## Check everything

From the repository root:

```bash
bash languages/go/check.sh
bash languages/go/check.sh race
```

Checks formatting, vet, tests, builds, and module-cache integrity. Race mode needs a supported C toolchain; on macOS, use Apple Command Line Tools. The script installs nothing and starts no infrastructure. Run `go test ./...` inside a module; the enclosing repository has no root Go module.

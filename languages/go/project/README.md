# Investor API

A small backend to read, run, change, and debug. Uses fictional names and process-local memory. No external services or libraries.

## Run

From the repository root:

```bash
cd languages/go/project
go run .
```

Expect a JSON log with `"msg":"listening"` and address `127.0.0.1:8080`. Leave this terminal running. Ctrl+C stops it.

From a second terminal:

```bash
curl -i --max-time 3 http://127.0.0.1:8080/health/live

curl -i --max-time 3 -X POST http://127.0.0.1:8080/api/v1/investors \
  -H 'Content-Type: application/json' \
  -d '{"name":"Asha","country":"IN"}'

curl -i --max-time 3 http://127.0.0.1:8080/api/v1/investors/1

curl -i --max-time 3 -X PUT http://127.0.0.1:8080/api/v1/investors/1 \
  -H 'Content-Type: application/json' \
  -d '{"name":"Asha Updated","country":"GB"}'

curl -i --max-time 3 'http://127.0.0.1:8080/api/v1/investors?after=0&limit=10'

curl -i --max-time 3 -X DELETE http://127.0.0.1:8080/api/v1/investors/1
```

On a fresh process, expect statuses **200, 201, 200, 200, 200, 204**, in order. Create assigns ID 1 and `Location: /api/v1/investors/1`. Repeated creates get new IDs; restarting clears all data. [Full API contract](../02-http-api.md).

To use another port, run `INVESTOR_ADDR=127.0.0.1:9090 go run .` and change the curl URLs. An explicitly empty value fails; missing uses the default.

## Read the code in order

| File | Follow this flow |
| --- | --- |
| [main.go](main.go) | Address → listener → signal context → server → bounded shutdown |
| [api.go](api.go) | Route → decode/ID checks → store call → status/JSON |
| [store.go](store.go) | Normalize/validate → lock → copy/update data → unlock |
| [api_test.go](api_test.go) | Successful CRUD, bad bodies, pagination, rejected writes |
| [store_test.go](store_test.go) | Validation, snapshots, ordering, concurrent access |
| [main_test.go](main_test.go) | Configuration, occupied port, active-request shutdown |

All source files share `package main`; `go run .` compiles them together. Keep this flat until splitting a package solves a real problem.

## Understand the new syntax

- `newStore()` allocates a map and returns `*Store`. `newHandler(store)` passes that same store into the HTTP layer; there is no dependency-injection container.
- `input, err := normalize(input)` reuses the input parameter and declares a new error in the same scope. A failed check returns before mutation.
- `item, ok := s.items[id]` reads a map value plus a presence flag. Missing data is different from a valid zero-valued field.
- `s.mu.Lock()`/`defer s.mu.Unlock()` protect the whole mutation; read locks allow simultaneous readers.
- `[]Investor` is a slice; `items[:limit]` keeps a prefix. `sort.Slice` receives a comparison function so map iteration order never becomes API ordering.
- `var input *Input` plus `Decode(&input)` lets the decoder leave JSON `null` as nil. A second decode must reach `io.EOF`, proving there is no extra value.
- `errors.Is` recognizes a domain cause; `errors.As(err, &tooLarge)` writes the matched typed error into `tooLarge`.
- `make(chan error, 1)` holds the server result even if the owner is handling cancellation. `go func() { ... }()` starts that work; `select` waits for failure or cancellation.

Imports have focused roles: `fmt`/`errors` describe failures; `strings`/`utf8` validate text; `sync` protects memory; `sort` stabilizes order; `json`/`mime`/`url`/`strconv`/`io` parse the boundary; `net`/`http` handle connections; `slog` records lifecycle events; `os`/`signal`/`syscall` receive process inputs; `context`/`time` bound cleanup.

## Change one thing at a time

1. Add PATCH and its HTTP tests.
2. Replace memory with PostgreSQL behind a small repository interface.
3. Add Gin, request logs, then a measured cache. Follow the numbered notes for concurrency and deployment.

Current limits: no persistence, auth, dependency integration, or request telemetry. Returned pages are bounded, but total stored records and the in-memory sorting work are not. Standard JSON accepts duplicate keys. This is a learning baseline, not a completed production system.

## Test and debug

From this directory:

```bash
GOWORK=off go test -count=1 -v ./...
GOWORK=off go test -race ./...
```

- **Address already in use:** identify your existing listener or choose port 9090; do not terminate unrelated processes.
- **400:** check the ID, required fields, JSON shape, and pagination.
- **415:** add `Content-Type: application/json`.
- **404:** check method/path, then whether the record exists in this process.
- **Data disappeared:** a restart creates a fresh memory store; persistence is the next data exercise.

Explain why the failure occurred, fix it, rerun the smallest relevant test, then run `bash ../check.sh` before moving on.

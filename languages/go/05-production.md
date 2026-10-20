# 5. Production habits

[Concurrency](04-concurrency.md) → **Production** → [Start](README.md)

Goal: make changes you can verify, explain, deploy, and recover.

Use [testing examples](#table-driven-tests), [configuration](#configuration-and-dependency-wiring), [logging](#structured-logging-example), [deployment](#container-build-and-runtime), and [troubleshooting](#production-style-troubleshooting) as separate reading stops.

## Test behavior at the right layer

| Layer | Prove | Example |
| --- | --- | --- |
| Unit | A rule or transformation | Invalid input does not change stored data |
| HTTP | Method, body, status, headers | Create returns 201 and `Location` |
| Integration | Real dependency behavior | A failed transaction leaves no partial write |
| End to end | A small real workflow | Create → restart → get |
| Failure test | Recovery behavior | Slow dependency respects deadline |

Use table-driven tests for cases; use small fakes at meaningful boundaries. A fake database cannot prove SQL, migrations, or transaction behavior. Keep real dependency tests isolated and explicitly runnable.

From `languages/go/project/`:

```bash
GOWORK=off go test -count=1 -cover ./...
GOWORK=off go test -race ./...
```

Coverage shows executed statements, not correctness. Test failure paths, ownership, authorization, cleanup, and concurrent changes. Use benchmarks and CPU/heap profiles to investigate a measured bottleneck before optimizing. [Go diagnostics](https://go.dev/doc/diagnostics).

## Trace startup and configuration

In an unfamiliar repo, find `main` → configuration → constructors → dependency checks → route registration → server → shutdown. Find who closes each resource and what happens if construction fails halfway.

`INVESTOR_ADDR=127.0.0.1:9090 go run .` passes a value from your shell to its child process. `os.LookupEnv` distinguishes missing from explicitly empty. An IDE has its own launch environment; `.env` files are not loaded by Go automatically.

Podman can inject variables into a container; Kubernetes can inject ConfigMap/Secret values or mount files. Each process reads its own environment. Updating a manifest does not rewrite a running process's environment; rollout/restart behavior matters.

## Security and secrets

Validate untrusted input, authorize each record access, parameterize SQL, bound requests, and set timeouts. Return safe errors; keep credentials and personal data out of logs. Browser cookie authentication also needs a deliberate CSRF policy; CORS is not authorization.

Keep production secrets outside source. For Vault, learn identity → least-privilege policy → secret retrieval → renewal/re-authentication → expiry/failure. Not every secret behaves like a renewable lease. Decide what startup and an already-running service should do when credentials expire. Practice only with dummy local credentials. [Vault leases](https://developer.hashicorp.com/vault/docs/concepts/lease).

Scan known vulnerabilities with `govulncheck`, reviewing affected call paths and the Go toolchain as well as dependencies. Pin chosen tool/library versions when introducing them; checksum verification is a different check. [Go security guidance](https://go.dev/doc/security/best-practices).

## Observe and debug

Logs describe events, metrics show trends, and traces connect work across services. The project uses `slog` for startup/shutdown JSON; Zerolog is a later alternative. Add request duration, status, route template, and a bounded request ID. Do not label metrics with raw user IDs or arbitrary URLs. [OpenTelemetry signals](https://opentelemetry.io/docs/concepts/signals/).

When an API fails: reproduce the request → identify its failing boundary → inspect the error/correlation → test one explanation → fix → add a regression test. For latency, separate handler time, pool wait, query time, cache time, and downstream time.

## Services, containers, deployment

Keep one service until independent ownership/scaling gives a concrete reason to split. Reuse HTTP clients, propagate context, close response bodies, and bound time/response sizes. Retry only appropriate failures with a total budget, backoff, jitter, and idempotency. A circuit breaker limits calls during failure; it does not repair a dependency.

On macOS, Podman runs Linux containers through a VM. `localhost` means the caller's own network environment. A host-run API reaches published dependency ports; container-to-container calls use their network names. Build a Linux binary for the container's architecture and run as a nonroot user.

For Kubernetes: a Deployment manages replicas; a Service provides discovery. **Startup** allows initialization, **readiness** controls eligibility for traffic, and **liveness** can trigger restart. Avoid a database-dependent liveness probe that restarts every API during a database outage. Match termination grace time to draining/cleanup. [Probe semantics](https://kubernetes.io/docs/concepts/workloads/pods/probes/).

OKD/OpenShift adds platform concepts such as Routes and security constraints; a local Kubernetes cluster does not provide every OpenShift API. Learn generic deployment first, then the environment-specific controls.

## Table-driven tests

Go tests live in `_test.go` files. A test receives `*testing.T`; failures should explain the input, actual result, and expected result. Subtests give related cases individual names.

Complete test file — save as `normalize_test.go` in a temporary module and run `go test -v`:

```go
package notes

import (
	"strings"
	"testing"
)

func normalizeName(value string) string {
	return strings.TrimSpace(value)
}

func TestNormalizeName(t *testing.T) {
	cases := []struct {
		name, input, want string
	}{
		{"spaces", " Asha ", "Asha"},
		{"empty", "", ""},
		{"unicode", " 界 ", "界"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			if got := normalizeName(test.input); got != test.want {
				t.Errorf("normalizeName(%q) = %q; want %q", test.input, got, test.want)
			}
		})
	}
}

func FuzzNormalizeName(f *testing.F) {
	f.Add(" Asha ")
	f.Add("")
	f.Fuzz(func(t *testing.T, input string) {
		once := normalizeName(input)
		if twice := normalizeName(once); twice != once {
			t.Fatalf("normalization is not idempotent: %q -> %q", once, twice)
		}
	})
}
```

The fuzz test checks a property, not a list of expected names. Normal test runs execute seed cases; active fuzzing needs an explicit command such as `go test -fuzz=FuzzNormalizeName -fuzztime=5s` in that temporary module. Save discovered regressions as part of the test corpus. [Fuzzing tutorial](https://go.dev/doc/tutorial/fuzz).

## Test doubles and HTTP tests

A stub supplies a response; a fake has simplified working behavior; a mock checks expected interactions. Choose a double based on the behavior the test must prove. An implementation-coupled mock can fail harmless refactors while missing a real contract bug.

Use `httptest.NewRecorder` for focused handler tests and `httptest.NewServer` when the actual HTTP client/transport boundary matters. Test cancellation, malformed replies, body limits, and dependency status mapping as well as success. [httptest reference](https://pkg.go.dev/net/http/httptest).

The current project includes both styles. Separate tests of the memory store cannot prove future PostgreSQL transactions or Valkey failover.

## Deterministic tests and integration boundaries

Prefer explicit signals over sleeps, isolated data over shared mutable fixtures, and injected clocks/ID generators when time/randomness affect a rule. `t.TempDir` and `t.Cleanup` make ownership visible; `t.Setenv` changes process state and should not be mixed with parallel tests that rely on that state.

Integration tests need a real compatible dependency, known schema/data, isolated names or databases, bounded startup waits, and cleanup. Avoid a test that silently skips every meaningful assertion because the dependency is missing.

A skipped integration suite is neither a failure nor evidence of passing integration behavior. State its prerequisites and exact command in the project once it exists.

## Benchmarks and profiles

Complete benchmark file — save as `join_test.go` in a temporary module:

```go
package notes

import (
	"strings"
	"testing"
)

func BenchmarkJoin(b *testing.B) {
	parts := []string{"investor", "demo", "42"}
	for b.Loop() {
		if got := strings.Join(parts, ":"); got != "investor:demo:42" {
			b.Fatal(got)
		}
	}
}
```

Run `go test -run '^$' -bench BenchmarkJoin -benchmem` there. `b.Loop` is available in this repository's Go version. Timing varies by machine/load; do not turn one sample into a universal claim. [testing reference](https://pkg.go.dev/testing).

CPU profiles show where execution time is spent; heap profiles distinguish allocation/retention views; block/mutex profiles investigate waiting. Profile a representative workload. An idle service's profile does not explain a loaded service.

Keep profiling access restricted and deliberate. Diagnostic endpoints may reveal internals or consume resources; do not accidentally expose a default profiling mux on the public listener.

## Debugging with Delve

Use a debugger when you need to inspect a value or control flow at the failure point. Set a breakpoint near the boundary, inspect inputs/errors, step through the relevant call, and compare the observed path with the contract.

For a local project, `dlv debug .` builds/starts it under Delve when Delve is installed and compatible. Investigate goroutines/stacks for blocking problems. An optimized release binary may offer a different debugging experience from a debug build.

Do not replace an unfamiliar implementation before locating the failure. First identify whether the fault is configuration, input, logic, synchronization, dependency behavior, or deployment.

## Configuration and dependency wiring

Parse configuration once at startup into typed values. Apply defaults only where they are meaningful, reject malformed/explicitly empty required values, and report the **field name** without leaking secret values.

Construct dependencies in a visible order. If the third constructor fails after two resources were acquired, close the resources already owned. A small `run` function returning an error is easier to reason about than hidden work in several `init` functions.

| Stage | Example responsibility |
| --- | --- |
| Parse | Port, durations, environment names |
| Validate | Nonempty required URL; positive bounded timeout |
| Construct | Pool/client/service/router |
| Verify | Required dependency connectivity/authentication |
| Serve | Start accepting work |
| Drain | Stop admission and finish/cancel work |
| Close | Release clients, pools, background workers |

The project traces `INVESTOR_ADDR`; do the same for every added field. A local `.env` loader, if introduced, is another explicit source with precedence rules, not Go runtime magic.

## Structured logging example

Complete program:

```go
package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	requestLog := logger.With("request_id", "demo-request")
	requestLog.Info("request finished", "route", "/api/v1/investors", "status", 200)
}
```

Output is a JSON event containing message, level, request ID, route, status, and a variable timestamp. Field names/types make events queryable. Context does not automatically add your request fields to every logger; propagate the intended logger/attributes deliberately. [slog reference](https://pkg.go.dev/log/slog).

Zerolog uses an event chain such as `logger.Info().Str("request_id", id).Int("status", 200).Msg("request finished")`. The final message/send operation emits the event. Check level configuration and avoid `Fatal` where process exit would bypass owned cleanup. [Zerolog reference](https://pkg.go.dev/github.com/rs/zerolog).

Log a failure at the layer that has the context to act on it. Avoid logging the same returned error at every layer. Do not serialize an entire request/DTO just because a logger supports arbitrary objects.

## Metrics, tracing, and cardinality

| Signal | Answer it helps with | Example |
| --- | --- | --- |
| Metric | How often/how much/how slow over time? | Request rate and latency histogram |
| Log | What happened in this event? | A classified dependency failure |
| Trace | Where did one distributed operation spend time? | API → database → downstream spans |

Use bounded labels such as route template, method, and status category. A label per user, request ID, or raw URL creates an unbounded series space. Choose base units and meaningful counter/histogram names. [Prometheus naming](https://prometheus.io/docs/practices/naming/).

Propagate trace context across supported transports. A correlation ID is useful but does not itself create parent/child spans or timings. Sampling and telemetry export failures are operational choices; telemetry should not become an uncontrolled failure amplifier.

## Service objectives and capacity

Define success from the user's operation: successful responses, acceptable latency, and data correctness. Separate an objective from an alert threshold and from a load-test result.

Measure rate, errors, latency distribution, saturation, queue lengths, and dependency capacity. Averages conceal tail latency. Compare p50/p95/p99 with enough representative samples and a known measurement window.

Load tests need representative operations/data, a controlled arrival model, and explicit pass criteria. Raising concurrency until a laptop crashes does not identify a useful operating point.

## Secrets, Vault, and rotation

Separate application identity from the secret it retrieves. Grant only required paths/actions. Static KV secrets, dynamic database credentials, and authentication tokens have different expiry/lease behavior. KV versioning does not mean KV secrets automatically renew like dynamic leases. [Vault KV](https://developer.hashicorp.com/vault/docs/secrets/kv).

Design renewal/re-authentication as application-owned work: bounded retries, cancellation, safe state publication, expiry handling, and shutdown. A token string shared across goroutines needs synchronized replacement if it changes.

Never assume a local development root token represents production identity. Test permission denied, revoked credentials, expired tokens, unavailable Vault, and rotation while requests are active. Keep dummy secrets out of logs too, so the exercised behavior resembles the intended policy.

## Secure service boundaries

Validate input and authorize at the use-case/data boundary. Use TLS where required, least-privilege database roles, constrained outbound destinations, and maintained cryptographic/password libraries rather than inventing protocols.

Review dependencies and toolchains, not just application code. `go mod verify` checks integrity of cached modules; `govulncheck` evaluates known vulnerabilities and affected usage. Neither replaces code review or a threat model.

For an unfamiliar endpoint, identify attacker-controlled inputs, identity source, protected data, outbound access, resource consumption, and safe failure behavior. Keep that review concrete rather than adding generic middleware without knowing its purpose.

## Service communication and resilience

Write a downstream contract for timeout, status/error mapping, maximum body, authentication, observability, retries, and idempotency. A missing optional profile and a failed required account lookup should not necessarily produce the same response.

Retry at a deliberate layer. If three nested layers each permit three attempts, the deepest operation can receive 27 attempts for one original request. Use a total time/attempt budget and avoid retrying overload immediately.

A circuit breaker limits calls while a dependency appears unhealthy; half-open probing must also be bounded. Bulkheads isolate capacity. Fallback must preserve business meaning and authorization, not merely produce a 200 response.

## REST, gRPC, and protocol choice

REST-style HTTP/JSON is straightforward to inspect with common tools and suits many external APIs. gRPC defines services/messages through Protocol Buffers and supports unary/streaming calls with generated clients and servers. [gRPC introduction](https://grpc.io/docs/what-is-grpc/introduction/).

Choose from consumer needs, compatibility, operational support, and streaming/contract requirements. A faster serialization benchmark does not establish better end-to-end reliability. Either approach still needs deadlines, authentication, authorization, bounded payloads, observability, and version compatibility.

## Messaging, webhooks, and scheduled jobs

For an at-least-once consumer, duplicate delivery is expected. Acknowledge only after the required effect is safely recorded. Coordinate deduplication with the state change when correctness requires it; a separate best-effort “seen ID” write can leave a crash gap.

Define ordering scope, retry/dead-letter policy, poison-message handling, and replay behavior. Global ordering often limits throughput; per-entity ordering may be enough.

Webhooks need authenticated payloads, replay protection according to the provider contract, bounded request processing, and durable handling if accepted asynchronously. Scheduled work across replicas needs an ownership/duplication policy; a ticker in every replica runs the job multiple times.

## Container build and runtime

Separate build tools from runtime requirements. A Go executable can often run without the compiler installed, but it may still need CA certificates, timezone data, dynamic libraries, or writable directories depending on its imports/build.

For this project's pure-Go Linux ARM64 build, from `languages/go/project/`:

```bash
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -o /tmp/backend-notes-investor .
```

This writes a Linux artifact; it does not run on macOS directly. Choose the architecture for the actual runtime, not automatically the developer laptop. CGO dependencies can change the cross-compilation requirements.

In a container, configure `INVESTOR_ADDR=0.0.0.0:8080` so the process accepts traffic on the container interface. For laptop-only access, publish to host `127.0.0.1`. Binding only container loopback usually prevents traffic from the published port reaching it.

Use an explicit reviewed image version/digest when implementing a Containerfile. Run with an appropriate nonroot identity, bounded resources, and the minimal filesystem permissions needed. Ensure the application receives termination signals. [Podman runtime options](https://docs.podman.io/en/latest/markdown/podman-run.1.html).

## Host, VM, container, and cluster networking

| Caller | Dependency address usually comes from |
| --- | --- |
| Go process on macOS | Host-published loopback port |
| Container in a named network | Dependency's network DNS name and container port |
| Pod in Kubernetes | Service DNS name and Service port |
| Container calling a host process | Provider-specific reachable host address, verified in that environment |

`localhost` always refers to the caller's own network namespace/environment. Verify both directions in a hybrid IDE/container workflow. A working client CLI on macOS does not prove a container can reach the same URL.

## Kubernetes and OKD building blocks

| Object/concept | Responsibility |
| --- | --- |
| Pod | Process/container execution unit with a lifecycle |
| Deployment | Desired replicas and rolling replacement |
| Service | Stable discovery/routing to selected ready endpoints |
| ConfigMap / Secret | Configuration delivery; access controls still matter |
| Namespace / RBAC | Organization and API permissions |
| Requests / limits | Scheduling/resource controls; reason about CPU/memory separately |
| Ingress / OpenShift Route | Platform-specific external routing |
| PersistentVolumeClaim | Request for storage; durability depends on storage implementation |

Secret base64 encoding is not encryption. A service account's access to cluster APIs is different from an end user's access to application records. OpenShift's security constraints can require images to work with an assigned arbitrary UID.

## Probes, rollout, and termination

Startup probes defer the other probes until initialization succeeds. Readiness controls eligibility for service traffic; liveness can restart an unhealthy container. Design each around its purpose rather than copying the same deep dependency check everywhere.

A rollout can temporarily increase replicas and therefore pool connections. Include that overlap in capacity planning. Readiness success is not proof of a completed migration or a safe mixed-version rollout unless you designed those checks.

On termination, traffic removal, grace periods, hooks, signal delivery, and application draining interact. Test the complete path instead of assuming a successful local Ctrl+C test proves cluster behavior. [Pod lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/).

## CI, dependencies, and releases

Use the existing `check.sh` as the local entry point. CI should run modules independently with `GOWORK=off`, verify formatting, vet/test/build, and add race/integration/security checks where supported and relevant.

Upgrade a dependency in a focused change: inspect its release notes, choose a version, update module metadata, run affected tests, and review behavior. Avoid combining an unexplained toolchain/library/image upgrade with a business feature.

Keep artifact identity, configuration, schema compatibility, rollout checks, and rollback instructions together. A container tag that moves does not reliably identify the binary that was tested. Activate/publish workflows or images only when that work is requested.

## Production-style troubleshooting

| Failure | Evidence and next step |
| --- | --- |
| Starts then exits | Exit status, first startup error, configuration provenance |
| Container healthy locally, unreachable remotely | Bind address, publication, DNS, network, proxy |
| CrashLoopBackOff | Previous container logs, termination reason, failing startup path |
| OOMKilled | Memory limit, retained heap, queues/caches, other memory usage |
| Readiness fails | Probe path/port, initialization, required dependencies |
| Image pull failure | Image identity, architecture, registry access/authentication |
| Latency spike | Rate, queue, CPU, pool/lock waits, downstream latency |
| Duplicate job | Retry/ack boundary, idempotency, multiple scheduler owners |
| Credentials suddenly fail | Expiry, rotation publication, policy, target identity |

For each incident, record the observed symptom and scope, relevant timeline, one tested explanation, corrective action, verification, and a prevention task. Preserve evidence before restarting away the useful state when practical.

## Code review questions

Can a reviewer trace the trigger to behavior and failure? Are ownership, cancellation, transaction boundaries, and auth decisions explicit? Does each introduced abstraction solve a current need? Do tests cover the changed contract and an important failure?

For a PR description, state the concrete before/after behavior and relevant verification. Do not claim cache failover, database rollback, or cluster deployment passed because a unit suite passed.

## Practice

1. Add request logging and a test proving request bodies/tokens do not appear in it.
2. Make a downstream test server slow or return 503. Show a bounded failure without duplicate side effects.
3. Containerize the API, then add deployment/probes. Demonstrate a failed readiness check and recovery before adding more infrastructure.

For CI, run the same local checks for both modules before packaging. A release also needs a schema-compatible rollback and a short runbook. Add infrastructure when you can explain the application's behavior first.

Final assessment: open an unfamiliar fictional Go service; trace startup/configuration; identify storage, clients, and deployment; implement one small behavior change; test a failure; and explain how to roll back safely. Repeat topics that still require copying the answer.

## Quick answers

| Question | Answer |
| --- | --- |
| Does 100% coverage prove correct authorization? | No; assertions and meaningful cases determine what was proved. |
| Does a unit fake prove transaction rollback? | No; verify actual storage behavior with an isolated real dependency. |
| Is a decoded token a verified identity? | No; cryptographic and claim validation still matter. |
| Does a ConfigMap update rewrite process environment? | No; launch/rollout behavior matters. |
| Is readiness the same as liveness? | No; traffic eligibility and restart decisions have different purposes. |
| Can a container use its own loopback address for a different container? | No; use the correct network name/address. |
| Can retries increase an outage? | Yes; nested attempts amplify load and latency. |
| Is a rollback only an older image? | No; configuration, schema compatibility, and data changes matter too. |

Check yourself: what does coverage miss? How does a value reach the process environment? Why should a database outage usually not fail liveness?

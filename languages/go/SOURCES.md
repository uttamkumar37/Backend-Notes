# Research and references

Reviewed **20 October 2026**. Prefer the official reference for the concept you are implementing; do not read everything before starting.

The Go track keeps five topic guides, two small practice programs, and one working HTTP baseline. The guides now include short reference sections, complete programs, pitfalls, exercises, and quick answers while retaining the same folder structure. The teaching order is a curriculum choice: language → HTTP → durable data → concurrency → production. Real-dependency integrations remain exercises until implemented and tested.

| Read when needed | Official source |
| --- | --- |
| Install/update Go | [Install](https://go.dev/doc/install), [release history and support policy](https://go.dev/doc/devel/release) |
| Check language semantics | [Specification](https://go.dev/ref/spec) |
| Understand dependencies/layout | [Modules](https://go.dev/ref/mod), [module organization](https://go.dev/doc/modules/layout) |
| Handle error identities | [Error wrapping, Is and As](https://go.dev/blog/go1.13-errors) |
| Understand JSON edge cases | [encoding/json](https://pkg.go.dev/encoding/json) |
| Check HTTP/status/shutdown contracts | [net/http](https://pkg.go.dev/net/http), [RFC 9110](https://www.rfc-editor.org/rfc/rfc9110.html) |
| Adopt the office HTTP framework | [Gin binding and validation](https://gin-gonic.com/en/docs/binding/binding-and-validation/) |
| Build database access | [Go database guide](https://go.dev/doc/database), [transactions](https://go.dev/doc/database/execute-transactions), [pgxpool](https://pkg.go.dev/github.com/jackc/pgx/v5/pgxpool) |
| Protect database invariants | [PostgreSQL isolation](https://www.postgresql.org/docs/current/transaction-iso.html), [indexes](https://www.postgresql.org/docs/current/indexes.html) |
| Understand cache failover limits | [Valkey Sentinel](https://valkey.io/topics/sentinel/) |
| Reason about concurrent code | [Context](https://go.dev/blog/context), [memory model](https://go.dev/ref/mem), [race detector](https://go.dev/doc/articles/race_detector) |
| Diagnose performance | [Go diagnostics](https://go.dev/doc/diagnostics) |
| Manage secrets and vulnerabilities | [Vault leases](https://developer.hashicorp.com/vault/docs/concepts/lease), [Go security](https://go.dev/doc/security/best-practices) |
| Understand observability/deployment | [OpenTelemetry signals](https://opentelemetry.io/docs/concepts/signals/), [Kubernetes probes](https://kubernetes.io/docs/concepts/workloads/pods/probes/) |

Additional references used for the expanded notes:

| Topic | Primary references |
| --- | --- |
| Collections and text | [Slices explained](https://go.dev/blog/slices-intro), [slices API](https://pkg.go.dev/slices), [maps explained](https://go.dev/blog/maps), [maps API](https://pkg.go.dev/maps), [strings and runes](https://go.dev/blog/strings) |
| Type abstractions and loops | [Methods](https://go.dev/tour/methods/1), [interfaces](https://go.dev/tour/methods/9), [generics tutorial](https://go.dev/doc/tutorial/generics), [loop-variable change](https://go.dev/blog/loopvar-preview) |
| Runtime and standard library | [GC guide](https://go.dev/doc/gc-guide), [time](https://pkg.go.dev/time), [io](https://pkg.go.dev/io), [os](https://pkg.go.dev/os), [bufio](https://pkg.go.dev/bufio), [strconv](https://pkg.go.dev/strconv) |
| Build and introspection | [Module versions](https://go.dev/doc/modules/version-numbers), [embed](https://pkg.go.dev/embed), [reflect](https://pkg.go.dev/reflect) |
| Cancellation and coordination | [context API](https://pkg.go.dev/context), [sync API](https://pkg.go.dev/sync), [errgroup API](https://pkg.go.dev/golang.org/x/sync/errgroup) |
| HTTP/security contracts | [Problem details](https://www.rfc-editor.org/rfc/rfc9457.html), [JWT best practices](https://www.rfc-editor.org/rfc/rfc8725.html), [OWASP CSRF prevention](https://cheatsheetseries.owasp.org/cheatsheets/Cross-Site_Request_Forgery_Prevention_Cheat_Sheet.html), [gRPC introduction](https://grpc.io/docs/what-is-grpc/introduction/) |
| SQL implementation | [database/sql](https://pkg.go.dev/database/sql), [Go query handling](https://go.dev/doc/database/querying), [constraints](https://www.postgresql.org/docs/current/ddl-constraints.html), [locking](https://www.postgresql.org/docs/current/explicit-locking.html), [EXPLAIN](https://www.postgresql.org/docs/current/using-explain.html) |
| Cache implementation | [Valkey data types](https://valkey.io/topics/data-types/), [SET options](https://valkey.io/commands/set/), [Go Redis client](https://redis.io/docs/latest/develop/clients/go/) |
| Tests and logging | [testing](https://pkg.go.dev/testing), [httptest](https://pkg.go.dev/net/http/httptest), [fuzz tutorial](https://go.dev/doc/tutorial/fuzz), [slog](https://pkg.go.dev/log/slog), [Zerolog](https://pkg.go.dev/github.com/rs/zerolog) |
| Operations | [Vault KV](https://developer.hashicorp.com/vault/docs/secrets/kv), [Podman run](https://docs.podman.io/en/latest/markdown/podman-run.1.html), [Pod lifecycle](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/), [Prometheus naming](https://prometheus.io/docs/practices/naming/) |

Version note: the host was observed running Go **1.27.1**. Official release history lists **1.27.3**, released 19 October 2026, including a security fix. No SDK installation/update was performed. Modules declare language minimum `go 1.27.0`; that is not a claim that the minimum patch is the recommended deployment version. Recheck supported patches when building/deploying.

Project baseline verified on 20 October 2026 with Go 1.27.1: `bash languages/go/check.sh` and its `race` mode passed formatting, vet, tests, builds, and module verification for both modules. A built binary passed live HTTP health/CRUD/pagination, 400/404/413/415 cases, and SIGTERM shutdown.

After expanding the notes, the fast repository checks passed again. The five guides contain 135 sections. All 18 fenced Go examples were parsed/formatted and checked in temporary modules: 14 standalone programs produced the expected output/fields, two database packages compiled, table tests and five seconds of active fuzzing passed, and the benchmark executed. The four concurrency programs ran with the race detector; worker cancellation/invalid-worker cases also passed. Database examples were not executed against PostgreSQL. All 82 local links/section anchors resolved.

These results verify the supplied code, not learner mastery or the unimplemented PostgreSQL/cache/deployment exercises.

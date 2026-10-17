# Troubleshooting and planned failure drills

For a current setup error, use [Session 0 debugging](notes/00-environment/DEBUGGING.md). The catalog below reserves real office-style scenarios for later hands-on lessons. The application, containers, and cluster named here are **planned**, not implemented or running as part of these notes.

Each drill will deliberately introduce one local fault, record its symptoms, gather evidence, identify the root cause, apply a targeted fix, and verify recovery with a regression check. Use only fictional lab data. Record results in [PROGRESS.md](PROGRESS.md).

## Start with the failing layer

Find the first failed boundary: configuration → process → listener → network/DNS → authentication → application logic → dependency. Compare logs using one request ID and record the elapsed time/status. A 500 is a symptom; it does not identify which dependency failed. Read the first useful cause before changing configuration or increasing timeouts.

Commands below assume the matching phase already exists. Container names, route paths, and the local context follow [COMMANDS.md](COMMANDS.md). `POD_NAME` means a pod selected from `get pods`; do not paste it unchanged. Only query named nonsecret configuration fields and avoid dumping full environments, tokens, or database URLs into logs.

## Startup and HTTP — Phases 2–4

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F01 Missing required environment variable | Process exits before accepting requests | Read first startup validation error; check the specific nonsecret variable using `printenv HTTP_ADDR` | Launch process did not inherit configuration; set the named value in Terminal/IDE/container launch configuration and restart that process | Valid config starts; missing-config automated test still rejects absence |
| F02 Gin endpoint returns 404 | Request reaches server but route is not found | `curl -i --max-time 3 http://127.0.0.1:8080/api/v1/investors`; compare method/prefix with router registration | Wrong method/path, group prefix, or wrong service port; correct the request or registration | Correct route/status and a test for the formerly failing path |
| F03 API returns 500 | Safe error response; request ID links to error log | Capture status/request ID with curl; follow structured error and wrapped cause in local logs | Repository/handler error incorrectly handled or dependency failing; repair identified cause and preserve safe public error | Happy path returns expected data; failure path has correct status and no internal secret details |

`printenv HTTP_ADDR` checks only this shell. It does not inspect an already running IDE/container process. Teach the variable's actual origin in the environment lesson before diagnosing it by name.

## PostgreSQL — Phase 5

These checks distinguish transport from authentication and SQL. The [psql manual](https://www.postgresql.org/docs/current/app-psql.html) documents connection options and command execution.

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F04 Connection refused | Immediate TCP refusal, startup or repository errors | `podman ps -a`, `podman port investor-postgres`, `pg_isready -h 127.0.0.1 -p 5432` | DB stopped, wrong host/port, or absent publication; fix that specific boundary | Same host/role/database used by app passes `SELECT 1`, then CRUD test |
| F05 Wrong database password | PostgreSQL reports password authentication failed | Use the documented `psql ... -W -c 'SELECT 1;'` against the same host/role/database; inspect DB error without logging password | Secret differs from role credential; correct config or reset only the isolated lab role using its admin procedure. Changing initialization env does not rewrite an existing volume's role password | Direct query and application query succeed with updated credential; bad credential still rejected |
| F06 Pool exhaustion (additional drill) | Pool acquisition waits then times out despite reachable DB | Pool statistics, query durations, `pg_stat_activity`, goroutine dump in local debugger | Leaked rows/transaction/connection, long query, or total pools exceed DB budget; fix ownership/queries/capacity before enlarging pool | Repeated concurrent workload has bounded waits, released connections, and passing cleanup tests |

## Valkey/Redis — Phase 6

Run cache and Sentinel diagnostics from a network where the announced addresses are reachable. [Valkey's Sentinel guide](https://valkey.io/topics/sentinel/) explains discovery and failover considerations.

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F07 Cache unavailable | Cache connection error; API should fall back within budget | `valkey-cli -h 127.0.0.1 -p 6379 PING`; inspect cache-error/DB-fallback counters | Cache stopped or wrong endpoint/auth; restore dependency while ensuring code handles transient cache failures | DB fallback returns correct investor during outage; cache traffic resumes after recovery |
| F08 Stale cached data | GET returns pre-update fictional investor | `GET investor:1` and `TTL investor:1` through CLI; compare authoritative DB row after update | Invalidation missing/wrong key or read/write race; invalidate after commit, bound TTL, and document concurrent consistency strategy | Update/read test and expiry test pass; measured DB/query behavior matches chosen consistency contract |
| F09 Sentinel primary failover | Brief write failure or client keeps dialing old primary | `valkey-cli -p 26379 SENTINEL get-master-addr-by-name lab-primary`; `SENTINEL CKQUORUM lab-primary`; client/Sentinel logs | Quorum, master-name mismatch, unreachable announced address, or plain client bypassing Sentinel; repair topology/discovery and use failover-aware client | Stop only the isolated lab primary; observe a new primary, reconnect/read/write recovery, and bounded API behavior during transition |

Each CLI operation after `valkey-cli` is a server command; Sentinel checks go to its port, not the data port. A three-Sentinel single-laptop exercise cannot demonstrate machine-level high availability.

## Vault — Phase 8

The Vault CLI is an optional tool in the isolated dev environment. Use selected metadata fields rather than sharing full token lookup output; [Vault token lookup](https://developer.hashicorp.com/vault/docs/commands/token/lookup) supports field output.

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F10 Authentication/policy failure | Login fails or secret read returns permission denied | Auth-loop error category; `vault token capabilities secret/data/investor-platform` when authenticated; compare configured KV version/path | Invalid auth input, role, policy, or KV v2 data path; correct the specific mapping rather than granting wildcard access | Authorized dummy read succeeds and disallowed path remains denied |
| F11 Expired Vault token | Previously valid reads start failing | In isolated lab, `vault token lookup -field=ttl` and renewal-loop logs; lookup may itself fail after expiry | Nonrenewable/expired token or stopped loop; reauthenticate using the dev identity and repair cancellation/backoff handling | Short-TTL test reaches expiry then recovers; no token values in logs; shutdown ends loop |

Vault development root tokens do not exercise realistic short-TTL renewal. The lesson must create an appropriate limited, short-lived test token to run F11.

## Concurrency and downstream services — Phases 7 and 9

The [Go race detector guide](https://go.dev/doc/articles/race_detector) explains executed-path detection. Deadline propagation and local goroutine inspection will be taught before these drills.

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F12 Downstream returns 503 | Investor request gets partial-failure response or bounded error | Same request ID across Investor/Account/Portfolio logs; call the configured dependency directly with bounded curl | Unready downstream or intentional fault; restore it, classify transient errors, and bound retries to idempotent operations and remaining deadline | Defined response during outage, no duplicate side effect, normal result after recovery |
| F13 Context deadline exceeded | Request stops near configured budget | Record each hop's elapsed/remaining time; reproduce slow fixture; inspect DB pool/query timing | Slow dependency, pool wait, or too-small budget; propagate context and repair slow boundary before changing budgets | Cancellation test terminates work and returns expected status within bounded tolerance |
| F14 Goroutine deadlock | Work hangs; runtime may report all goroutines asleep in simple programs | Pause in Delve and inspect goroutine stacks; run isolated test with `go test -timeout 5s ./...` | Send has no receiver, lock cycle, or incomplete channel closure; repair ownership/order and add cancellation | Formerly hanging case terminates with correct result and bounded test time |
| F15 Data race | Intermittent output or detector warning with conflicting accesses | From affected module: `go test -race -count=1 ./...`; follow both stack traces | Unsynchronized shared memory; use clear ownership, lock, channel, or suitable atomic invariant | Same workload passes race instrumentation and logical correctness checks; explain coverage limits |

Do not deliberately introduce races/deadlocks into the default passing platform. The future lesson will isolate broken examples and provide opt-in reproduction commands.

## Podman networking and processes — Phase 10

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F16 Podman network issue | DNS lookup/refusal between containers | `podman network ls`; inspect selected network membership; `podman port investor-postgres` | Containers on different networks, missing alias, or host port confused with container port; repair network/URL | Probe from the actual caller's runtime reaches target and application request succeeds |
| F17 One service cannot reach another | Host curl works but a service call fails | Compare configured base URL and bind address; caller-side logs; inspect target listener | `localhost` points to caller, wrong DNS alias, or target binds only loopback; use appropriate network DNS/published port and listener | Test all required directions: host→container, container→container, and optional container→host |
| F18 Container exits immediately | `podman ps` omits it; `ps -a` shows exited | `podman ps -a`; `podman logs investor`; `podman inspect --format '{{.State.ExitCode}}' investor` | Startup validation, binary/OS architecture, missing runtime file/CA, permissions, or entrypoint error; fix first recorded cause | Container stays running, passes health check, and cleanly handles shutdown |

## Kubernetes — Phase 11

Use explicit `--context kind-go-lab -n investor-lab` for this local lab. The [Kubernetes pod debugging guide](https://kubernetes.io/docs/tasks/debug/debug-application/debug-running-pod/) covers status, events, and logs.

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F19 CrashLoopBackOff | Pod repeatedly restarts with increasing backoff | `kubectl --context kind-go-lab -n investor-lab describe pod POD_NAME`; `logs POD_NAME --previous` with same context/namespace | Startup crash, invalid config, permissions, probe kills, or OOM; use terminated reason/events to identify cause | Correct rollout becomes available and restart count stops increasing |
| F20 Readiness failure | Running pod absent from ready service endpoints | `describe pod` events; compare probe path/port and app readiness/dependency logs; inspect EndpointSlices | Wrong probe or real unavailable dependency; repair actual cause and allow startup time through an appropriate startup probe | Pod Ready, endpoint present, routed request succeeds; failed dependency produces intended readiness behavior |
| F21 ImagePullBackOff / OOMKilled (additional drills) | Image pull retry, or terminated reason `OOMKilled` | Pod events/image name/architecture; terminated state and measured memory | Missing image/tag/registry access, unimported local image, or insufficient limit/leak; fix image distribution or measured memory cause | Correct image runs on local node architecture and workload stays within measured budget |

CrashLoopBackOff is a restart/backoff state, not the root cause. Failed readiness alone does not restart a container; liveness/startup failures can.

## Git and modules — Phases 1 and 13

| ID / fault to introduce | Symptoms | Diagnostic commands/evidence | Root cause and targeted fix | Recovery verification |
| --- | --- | --- | --- | --- |
| F22 Git merge conflict | Unmerged paths and conflict markers | `git status --short`, `git diff --name-only --diff-filter=U`; read both intended changes | Overlapping edits; resolve semantics in a disposable local branch, format, test, and stage only resolved paths | No unmerged paths/markers; relevant tests and review confirm both intended behaviors |
| F23 Go dependency build failure | Missing module, incompatible API, checksum, or version error | `go version`, `go env GOMOD GOWORK`, first build error, `go list -m all`, `go mod verify` | Wrong working directory/toolchain/version, missing requirement, or altered cache; fix identified requirement/source. Do not disable checksum verification or clear every cache as the first step | Fresh module-scoped build/tests succeed; dependency file changes are understood and reviewed |

The original request's 21 scenario categories are covered here; pool exhaustion and image/OOM failures add related drills, producing 23 catalog entries. Entries become completed exercises only when their phase supplies reproducible setup, automated checks, observed evidence, and recovery steps.

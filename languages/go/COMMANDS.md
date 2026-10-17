# Command reference

Start with the inventory. Later sections are reference commands for phases whose modules/services have been created; they are **not** instructions to launch the entire stack now. Each lesson will explain new command syntax and provide exact expected output for its implemented code.

## Available now — inventory

From any directory:

```bash
bash /Users/uttamkumar/Uttam-workspace/GitHubProject/Backend-Notes/languages/go/scripts/check-environment.sh
```

Expected: macOS, architecture, tool versions, Podman VM table, and optional-tool paths/missing entries. The script deliberately continues when a tool is missing; review every reported failure. See [Session 0](notes/00-environment/README.md) for the meaning of each command.

Inspect this repository's changes from the repository root:

```bash
cd /Users/uttamkumar/Uttam-workspace/GitHubProject/Backend-Notes
git status --short
git diff -- languages/go
```

`cd` changes working directory. `git status --short` reports modified/untracked paths. `git diff -- languages/go` limits the tracked-file diff to this folder; untracked files are listed by status and need opening separately.

## Current lesson — Phase 1 lesson 2

From any directory:

```bash
cd /Users/uttamkumar/Uttam-workspace/GitHubProject/Backend-Notes/languages/go/01-go-fundamentals/02-modules-packages
GOTOOLCHAIN=local GOWORK=off go run .
GOTOOLCHAIN=local GOWORK=off go test -v ./...
```

Expect `Welcome, Asha.` and passing command/library checks. `GOWORK=off` makes this an independent module test. The first test run may download the pinned `go-cmp` helper. The [lesson](01-go-fundamentals/02-modules-packages/README.md) explains the source, dependency files, and commands.

To test Lessons 1 and 2 through the supplied workspace, continue from that module root:

```bash
cd workspace
GOTOOLCHAIN=local GOWORK="$PWD/go.work" go test -v ../... ../../01-first-program/...
```

`$PWD/go.work` explicitly selects this workspace. The two patterns select both member modules' package trees; `./...` from the workspace folder alone does not do that. Use `cd ..` afterward to return to the Lesson 2 root.

## Earlier lesson — Phase 1 lesson 1

From any directory:

```bash
cd /Users/uttamkumar/Uttam-workspace/GitHubProject/Backend-Notes/languages/go/01-go-fundamentals/01-first-program
GOTOOLCHAIN=local go run .
GOTOOLCHAIN=local go test -v ./...
```

The program prints the welcome and next-step lines. Both `TestWelcomeMessageUsesArguments` and `Example` should pass. The [lesson](01-go-fundamentals/01-first-program/README.md) explains every source/test line and command, including building a native executable. No external dependencies are needed.

## Phase 1 onward — from an individual Go module root

**Prerequisite:** select the intended directory containing `go.mod`; the current module is `01-go-fundamentals/02-modules-packages/`. Do not run these from the enclosing Backend-Notes root and expect every nested module to be covered. The reference includes commands introduced in later lessons; today's required commands are explained above and in Lesson 2.

| Command | Purpose | Expected result |
| --- | --- | --- |
| `go fmt ./...` | Format packages recursively in the current module | Names of changed files, or no output; modifies source formatting |
| `go vet ./...` | Inspect suspicious constructs | No diagnostics on success |
| `go test ./...` | Run discovered package tests | `ok` package lines; cached results may be marked |
| `go test -count=1 ./...` | Run again without reusing cached successful results | Fresh test results |
| `go test -race ./...` | Run executed tests with race instrumentation | No race report; requires supported target/toolchain and CGO support |
| `go test -cover ./...` | Show statement coverage of exercised code | Coverage percentage; not a correctness certificate |
| `go build ./...` | Check that all packages/binaries compile | No output on success |
| `go mod tidy` | Reconcile module requirements/checksums | Review changes in `go.mod`/`go.sum` |
| `go mod download` | Populate module download cache | Normally quiet; may use network |
| `go mod verify` | Check downloaded module cache contents against recorded hashes | `all modules verified`; does not scan for vulnerabilities |
| `GOWORK=off go test ./...` | Test module without a surrounding workspace | Exposes dependencies accidentally supplied only by a workspace |

`./...` is Go's package pattern under the current directory, subject to module boundaries. It is not a shell wildcard for all modules. `:=` and Go pointers are not shell syntax; they will be explained in Go lessons, not inferred from these commands. Dependency behavior follows the [Go Modules Reference](https://go.dev/ref/mod), and instrumentation is described in the [race detector guide](https://go.dev/doc/articles/race_detector).

## Phase 3 onward — HTTP inspection

**Prerequisite:** the Investor API from its lesson is running on macOS at port 8080, with the routes below implemented.

```bash
curl -i --max-time 3 http://127.0.0.1:8080/health/live
curl -sS --max-time 3 http://127.0.0.1:8080/api/v1/investors
```

`-i` includes HTTP status/headers. `--max-time 3` bounds curl's total request time in seconds. `-sS` hides the progress meter but preserves errors. Expect status 200 and JSON on the healthy path. Curl normally exits successfully even for HTTP 500; check the status explicitly or use an appropriate HTTP-failure option when scripting. Later authenticated routes need the lesson's dummy identity headers/token.

To display JSON fields when the response is valid JSON and `jq` is installed, pipe the second response to `jq .`; the dot filter displays the parsed value. Keep malformed responses visible while debugging rather than assuming every body is JSON.

## Phase 5 onward — dependencies

**Prerequisite:** the matching local dependency lesson has created the named containers, roles, schema, and ports. A host CLI is optional if the same tool exists in its container.

| Command | What it checks / expected result |
| --- | --- |
| `podman machine list` | Existing VM names and state |
| `podman info` | Engine/API reachable; client version alone cannot establish this |
| `podman ps -a` | Running and stopped containers; inspect the specific lab container |
| `podman port investor-postgres` | Host mapping for planned database container |
| `psql -h 127.0.0.1 -p 5432 -U investor -d investor -W -c 'SELECT 1;'` | Host-to-DB TCP/auth/query; prompts for the dummy lab password; one row |
| `valkey-cli -h 127.0.0.1 -p 6379 PING` | Planned local cache responds `PONG` |
| `valkey-cli -h 127.0.0.1 -p 6379 TTL investor:1` | Positive seconds until expiry; `-1` means no expiry, `-2` means absent |

Here `-h`, `-p`, `-U`, and `-d` select host, port, role, and database for `psql`; `-W` prompts, and `-c` runs SQL. Names/ports are the intended lab convention, not an observed installed database. See [PostgreSQL's psql manual](https://www.postgresql.org/docs/current/app-psql.html).

## Phase 11 onward — explicit local cluster

**Prerequisite:** Phase 11 created Kind cluster `go-lab`, namespace `investor-lab`, and deployment `investor`. Use the explicit local context. Replace the literal `POD_NAME` below with a name from `get pods` before running pod-specific commands.

```bash
kubectl --context kind-go-lab -n investor-lab get pods
kubectl --context kind-go-lab -n investor-lab describe pod POD_NAME
kubectl --context kind-go-lab -n investor-lab logs POD_NAME --previous
kubectl --context kind-go-lab -n investor-lab rollout status deployment/investor
```

`--context` selects the cluster connection; `-n` selects namespace. `describe` includes state/events. `--previous` requests logs from the prior terminated container instance and may fail if none exists. `rollout status` reports whether the deployment has reached its desired state. No local cluster is created by this reference. See Kubernetes' [pod debugging guide](https://kubernetes.io/docs/tasks/debug/debug-application/debug-running-pod/).

## Phase 13 onward — vulnerability checks

The implemented lesson will install/pin `govulncheck` and run `govulncheck ./...` from each module. It examines known vulnerabilities relevant to Go code; it complements tests and dependency verification. See [Go Vulnerability Management](https://go.dev/doc/security/vuln/). No scanner installation or scan is claimed during Session 0.

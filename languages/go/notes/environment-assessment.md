# Observed environment — 17 October 2026

Captured through read-only commands in this workspace. These are machine observations; learner knowledge was not assessed. No tools were installed, no VM was created/started, no application secrets were read, and no office access was used.

**Subsequent checkpoint:** you authorized moving to [Phase 1, lesson 1](../01-go-fundamentals/01-first-program/README.md). The original observations below remain unchanged; authorization does not establish completion of unreported exercises or an SDK update.

## Core findings

| Check | Observed result | Meaning / next step |
| --- | --- | --- |
| macOS | 26.5.2, build 25F84 | Local environment for the initial lessons |
| Architecture | `arm64` | Native Apple Silicon process |
| Physical memory | 17,179,869,184 bytes = 16 GiB | Stage larger infrastructure rather than launching everything now |
| Apple developer tools | `/Library/Developer/CommandLineTools` | Selected tools directory exists |
| Homebrew | 7.0.9; `/opt/homebrew/bin/brew` | Available on `PATH` |
| Go | `go1.27.1 darwin/arm64`; `/opt/homebrew/bin/go` | Native installed SDK; targeted patch update documented in SETUP |
| Go build target | `darwin`, `arm64` | Matches the local host for these checks |
| Git | 2.50.1 (Apple Git-155); `/usr/bin/git` | Sufficient for local Git workflow |
| Podman client | 6.1.3; `/opt/homebrew/bin/podman` | Installed; engine workloads not tested by inventory |
| Podman Machine | `podman-machine-default`, `libkrun`, already running | 6 vCPUs, 8 GiB memory, 100 GiB virtual disk allocation |
| VS Code | 1.142.0, `arm64`; launcher on `PATH` | App present; Go extension/SDK configuration not inspected |
| `curl`, `jq`, `make` | Present under `/usr/bin` | Available; versions/features not exhaustively assessed |

`GOTOOLCHAIN=local` was set only on the inventory's Go commands to inspect the installed SDK; this report makes no claim about a persistent Go setting. The VM disk figure is virtual capacity, not measured free disk space.

## Not found on the current terminal PATH

`gh`, `psql`, `valkey-cli`, `redis-cli`, `dlv`, `golangci-lint`, `kubectl`, `oc`, `kind`, `helm`, `goland`, `vault`, `prometheus`, `grafana`, and `bruno`.

These findings do not prove that corresponding GUI applications, container images, or binaries outside `PATH` are absent. In particular, GoLand may be installed through Toolbox or another directory. Missing future-phase tools do not block learning basic Go.

## Unverified and intentionally separate

- Editor Go extension, selected SDK, and a working debugger breakpoint.
- GitHub/Copilot authentication, IDE licenses, or personal account state.
- Running containers, PostgreSQL/cache availability, VM engine API, or available disk space.
- Any Kubernetes context or corporate infrastructure.
- Learner understanding, exercise completion, and phase completion.

The read-only inventory script ran successfully with no `CHECK FAILED` entries. Its zero exit code indicates that collection completed, not that every tool was found or every future environment is ready.

The official [Go downloads](https://go.dev/dl/) and [Homebrew formula](https://formulae.brew.sh/formula/go) listed Go **1.27.2** as stable when checked on this date. The installed SDK remains **1.27.1**; no update was applied. Review the targeted update instructions in [SETUP.md](../SETUP.md) before the first module.

**Current learning checkpoint:** see [PROGRESS.md](../PROGRESS.md), which tracks learner confirmations and the active lesson. Refer to [Session 0](00-environment/README.md) if a setup gap appears; the original machine observations above are a dated snapshot.

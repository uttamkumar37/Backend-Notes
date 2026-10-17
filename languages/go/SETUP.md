# macOS setup — tools introduced when needed

Start with [Session 0](notes/00-environment/README.md). This file is a staged reference, not a command block to run in full. The current [assessment](notes/environment-assessment.md) found the core tools already installed. No installation or configuration change has been performed for this curriculum.

## Tool locations and verification

| Tool | Where it runs | Needed | Verify / expected result |
| --- | --- | --- | --- |
| Homebrew | Native macOS; normally `/opt/homebrew` on Apple Silicon | Setup | `brew --version`: version line |
| Go SDK/toolchain, `gofmt`, tests | Native macOS; Linux toolchain in container builds later | Phase 1 | `GOTOOLCHAIN=local go version`: installed SDK and `darwin/arm64` |
| VS Code + Go extension | Native macOS | One editor now | Open a `.go` file; SDK recognized, navigation works |
| GoLand | Native macOS; alternative editor | Optional | Open SDK settings; compare terminal SDK; launcher may be absent |
| Git | Native macOS; Apple command-line tools are sufficient initially | Now | `git --version`: version line |
| GitHub CLI (`gh`) | Native macOS | When using GitHub features | `gh --version`: version; local lessons do not require login |
| GitHub Copilot | Extension/plugin in the editor, with account service | Optional | Editor shows enabled state; inspect suggestions before accepting |
| Podman CLI / Podman Machine | CLI on macOS; engine inside Linux VM | Phase 5 minimal; Phase 10 full | `podman machine list`, then `podman info`: VM/engine status |
| PostgreSQL server | Linux container in Podman VM | Phase 5 | `psql ... -c 'SELECT 1;'`: one row |
| `psql` / `pg_isready` | macOS client or inside database container | Phase 5 | `psql --version`: client version; does not prove DB health |
| Valkey server and Sentinel | Linux containers in Podman VM | Phase 6 | CLI `PING`: `PONG`; Sentinel checks later |
| Valkey/Redis CLI | macOS client or cache container | Phase 6 | `valkey-cli --version` or `redis-cli --version` |
| Vault dev server | Isolated local Linux container; dummy data only | Phase 8 | Authenticated dummy secret read succeeds; dev mode is ephemeral |
| Postman or Bruno | Native desktop client; `curl` is sufficient to start | Phase 3 | Send a request to the running local API |
| `curl`, `jq`, `make` | Native macOS; Apple Make is available here | As needed | `curl --version`, `jq --version`, `make --version` |
| Delve (`dlv`) | Native macOS debugger for local Go process | Early debugging | `dlv version`, then stop at an actual breakpoint |
| `golangci-lint` | Native macOS and Linux CI | As checks expand | `golangci-lint version`; configuration must match installed major |
| `kubectl`, Kind, Helm | Native CLIs; Kind nodes run in Linux containers | Phase 11 | Client version checks, then `kubectl get nodes` after cluster creation |
| `oc` | Native OpenShift client where available | Phase 11 comparison | `oc version --client`; this does not create an OpenShift cluster |
| Prometheus, Grafana, trace backend | Optional local containers/pods | Phase 13 | Ready endpoints, scrape success, query a known metric, open a dashboard |

An absent CLI is not proof a GUI application is absent. Installed commands and running services are distinct. All database/cache/HTTP verification targets below refer to the fictional lab introduced later.

## Core setup

### Homebrew and Apple command-line tools

Homebrew is already available at `/opt/homebrew/bin/brew`. If absent on another Mac, follow the [official Homebrew installation instructions](https://docs.brew.sh/Installation), review the installer, and use its displayed shell-setup instructions. Apple Silicon and Intel installations use different prefixes.

```bash
xcode-select -p
brew --prefix
```

`xcode-select -p` prints the active developer-tools directory; it does not install Xcode. `brew --prefix` prints Homebrew's installation root. Expect `/Library/Developer/CommandLineTools` and `/opt/homebrew` here. If command-line tools are missing, Apple's `xcode-select --install` launches their installer; run it only when needed.

### Go SDK and version policy

Observed SDK: **Go 1.27.1**. The [official Go downloads](https://go.dev/dl/) and [Homebrew Go formula](https://formulae.brew.sh/formula/go) listed **1.27.2** as stable when checked on **17 October 2026**. Refresh these sources when installing later; do not infer that the recorded version is permanently current.

For this Homebrew-managed installation, the targeted update is:

```bash
brew update
brew upgrade go
GOTOOLCHAIN=local go version
```

`brew update` refreshes package metadata; `brew upgrade go` changes the installed Go package and may update required dependencies; the final command verifies the installed SDK. These commands are documented here and have **not** been executed. If Go is absent, use `brew install go` instead. An official `darwin-arm64` installer is an alternative; choose one installation method and check `command -v go` rather than stacking conflicting SDK installations.

The SDK bundles the compiler and core tools. Toolchain selection can also follow module configuration; the inventory temporarily uses `GOTOOLCHAIN=local`. Future module and CI versions will be selected together, with exact dependency versions in `go.mod` and checksums in `go.sum`. Avoid teaching that `go.sum` is a lockfile or a list of only currently imported packages.

### Editor and debugger

VS Code is already present. Install the official **Go** extension published by the Go team (`golang.go`) through Extensions when preparing the editor. Open the lesson's folder and use **Go: Locate Configured Go Tools** to inspect the SDK/tools. Tool installation prompts can download `gopls` and Delve; inspect what is being installed. Navigation/testing/debugging are documented in [Go in VS Code](https://code.visualstudio.com/docs/languages/go).

If you prefer GoLand, select the same native SDK in its Go SDK/GOROOT settings; open it from Applications or Toolbox. You only need one editor. A working breakpoint is the real debugger check; a version line alone is insufficient.

Optional native debugger and linter installations, when their lesson needs them:

```bash
brew install delve golangci-lint
dlv version
golangci-lint version
```

Homebrew's [Delve formula](https://formulae.brew.sh/formula/delve) installs `dlv`; its [golangci-lint formula](https://formulae.brew.sh/formula/golangci-lint) installs the linter. Choose releases supporting the selected Go version, and later pin matching CI versions.

### Git, GitHub CLI, and optional Copilot

Apple Git works for local lessons. Later, `brew install gh` installs the [GitHub CLI](https://formulae.brew.sh/formula/gh); `gh --version` verifies installation. GitHub authentication is a separate user action and unnecessary for writing/running local examples.

Copilot is optional. Follow [GitHub's IDE setup instructions](https://docs.github.com/en/copilot/how-tos/copilot-in-your-ide/set-up-copilot/install-copilot-extension) in your chosen editor and sign in yourself if desired. Availability depends on account access; no paid subscription is a prerequisite for this curriculum. Verify it in a disposable fictional example, then review and test the suggestion.

## Phase 3 — API clients and small terminal tools

`curl`, `jq`, and `make` are already present. `curl` makes HTTP requests, `jq` filters JSON, and `make` runs named command recipes. macOS Make can be older than GNU Make tutorials assume; recipes will use compatible features or document an explicit alternative.

Choose Postman or Bruno through its official macOS installer if a GUI helps; `curl` covers the same initial request exercises. No account or cloud workspace is needed for local API tests. Later lessons include exact methods, headers, bodies, expected statuses, and API startup commands before asking you to send requests.

## Phases 5–6 — Podman, PostgreSQL, Valkey

The existing Podman installation works as a client, and its VM was already running during inventory. Do not initialize another machine. For a new installation, Podman's [current installation guide](https://podman.io/docs/installation) recommends its macOS installer; Homebrew is a community alternative. Keep the existing working setup unless there is a concrete reason to change it.

At the container lesson, inspect the machine first:

```bash
podman machine list
podman info
```

If no machine exists, `podman machine init` creates one and `podman machine start` starts it. If an existing machine is stopped, start that machine by name. These are state-changing actions for the later lesson, not inventory steps. [Podman Machine documentation](https://docs.podman.io/en/latest/markdown/podman-machine.1.html) covers VM management.

Use ARM64-compatible PostgreSQL and Valkey images with explicit versions recorded in the lab configuration. Publish development ports on `127.0.0.1`, use named volumes for intentional persistence, and avoid a second native database server competing for port 5432. Do not pull images or start services during Session 0.

Optional host clients, installed only when needed:

```bash
brew install libpq
export PATH="$(brew --prefix libpq)/bin:$PATH"
psql --version
```

`libpq` includes PostgreSQL client tools and is [keg-only](https://formulae.brew.sh/formula/libpq). `$(brew --prefix libpq)` captures its installation path; `/bin` selects executables; `export PATH=...` prepends that directory for this shell and its future children. It preserves the old path and does not start PostgreSQL. Use the container's client instead if you prefer no host client installation.

```bash
brew install valkey
valkey-cli --version
```

The [Valkey formula](https://formulae.brew.sh/formula/valkey) also includes compatible Redis-named tools. Install one tool family; the formula can conflict with a separate Redis installation. The server remains in Podman for our default workflow. Do not run `brew services start valkey` simply to obtain its CLI.

## Phase 8 — local Vault

Run Vault only with dummy secrets in the lab's loopback-bound development container. The later lesson will provide the exact image/version, dev-only token, policy, secret path, and authenticated verification commands together. Separate dev-mode behavior from real sealing, storage, identity, HA, and token renewal. No corporate credentials or Vault instance is needed.

## Phase 10 — Compose and hybrid debugging

`podman compose` delegates to an **external Compose provider**; installing Podman alone need not provide one. See the [Podman Compose command documentation](https://docs.podman.io/en/latest/markdown/podman-compose.1.html). One option is the [podman-compose formula](https://formulae.brew.sh/formula/podman-compose):

```bash
brew install podman-compose
PODMAN_COMPOSE_PROVIDER=podman-compose podman compose version
```

The assignment selects the provider for that command. Verify the chosen provider before using a compose file. The complete environment and its cleanup commands will arrive with the phase; no compose file exists yet.

For hybrid work, Investor runs in GoLand/VS Code on macOS and calls dependencies through published `127.0.0.1` ports. Containerized services use container DNS names on their network. To call a host service from a container, evaluate `host.containers.internal` with the Podman network in use; a service listening only on host loopback may need a deliberately reachable bind address. Verify each direction explicitly. `localhost` always refers to the caller's own network environment.

## Phase 11 — local Kubernetes and OpenShift client

Install these only when reaching deployments:

```bash
brew install kubernetes-cli kind helm
kubectl version --client
kind version
helm version --short
```

The Homebrew formulae are [kubernetes-cli](https://formulae.brew.sh/formula/kubernetes-cli), [Kind](https://formulae.brew.sh/formula/kind), and [Helm](https://formulae.brew.sh/formula/helm). Versions must be chosen together with the Kind node image and supported Kubernetes client/server skew. A latest client alone does not establish compatibility with an older cluster. Helm major versions can affect tutorial commands and charts.

Kind supports a [Podman provider](https://kind.sigs.k8s.io/docs/user/quick-start/) with [rootless prerequisites and troubleshooting](https://kind.sigs.k8s.io/docs/user/rootless). The later cluster lesson will validate the VM's cgroup/resource/logging requirements before using `KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster`. Do not run Linux `systemctl` commands on macOS or replace the Podman configuration to address an unobserved issue.

Where an ARM64 package is available, `brew install openshift-cli` installs [`oc`](https://formulae.brew.sh/formula/openshift-cli); verify with `oc version --client`. Kind covers the required local Kubernetes practice. An OpenShift client does not add Route/SCC APIs to Kind. A full local OpenShift distribution is optional and must be checked for current architecture/resource support before installation.

## Phase 13 — monitoring and reproducible builds

Prometheus, Grafana, and a local trace backend will be added as optional container/cluster profiles, with fictional data, pinned versions, and loopback access. The 16 GiB Mac and existing 8 GiB VM should be treated as a finite shared budget: introduce one stack at a time and measure usage before combining a cluster, databases, Sentinel, and dashboards. No VM resource change is needed now.

Verification will include a known generated metric, a successful Prometheus scrape, a Grafana query/dashboard, and a trace crossing services. Pin the Go toolchain, module dependencies, image versions/digests, and CI tools for the implemented project. Recheck compatibility whenever upgrading; this setup guide alone is not a reproducible build manifest.

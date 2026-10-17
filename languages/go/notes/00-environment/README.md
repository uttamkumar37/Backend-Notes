# Session 0 — Assess your MacBook development environment

**Time:** 45–90 minutes, including review and exercises. **Prerequisites:** macOS Terminal and your existing Java/backend knowledge. **Objective:** explain which tools are available, which architecture they target, and what is actually needed before the first Go lesson.

This was the initial setup lesson. You subsequently authorized starting [Phase 1, lesson 1](../../01-go-fundamentals/01-first-program/README.md). Its exercises remain available for reference; authorizing progression does not imply you completed unreported knowledge checks. The [roadmap](../roadmap.md) describes the full sequence.

## What → why → when

A development environment is the collection of programs that edit, compile, test, run, and inspect your application. The terminal finds commands through `PATH`, an ordered list of directories. A tool can exist on disk without being discoverable through that list.

Checking first prevents unnecessary reinstalls and separates a missing command from a broken application. You repeat these checks when changing laptops, upgrading Go, or investigating a build that works for a colleague but fails locally.

Your Mac uses `arm64`. Go calls macOS `darwin`, so the native Go target is `darwin/arm64`. Linux containers on the same processor use `linux/arm64`; sharing a CPU architecture does not make a macOS executable a Linux executable.

## How — one small command set

Open **Terminal.app** and run the following block. You may paste it as a group; each line is an independent check. None installs software or starts services.

```bash
sw_vers
uname -m
command -v brew go git podman
brew --version
GOTOOLCHAIN=local go version
GOTOOLCHAIN=local go env GOOS GOARCH
git --version
podman --version
podman machine list
```

If a command says `command not found`, retain the message and continue with the next line. Do not install every roadmap tool now. Alternatively, the repository script runs these checks, handles missing commands, and lists optional tools:

```bash
bash /Users/uttamkumar/Uttam-workspace/GitHubProject/Backend-Notes/languages/go/scripts/check-environment.sh
```

`bash` runs the script file. The absolute path identifies this exact copy independently of your working directory. The script uses two helper functions: `run_check` prints a command and reports failures while continuing; `version_check` checks whether a command exists before running it. Its optional-tool loop prints paths only. Review its source before running it to see every operation.

## Explain every command

| Command | What each part means | Expected result and what it proves |
| --- | --- | --- |
| `sw_vers` | macOS system-version utility; no arguments needed | `ProductName`, `ProductVersion`, and `BuildVersion`; identifies the OS |
| `uname -m` | `uname` prints system information; `-m` requests machine architecture | `arm64` in a native Apple Silicon terminal; `x86_64` may indicate a terminal running through Rosetta |
| `command -v brew go git podman` | Shell builtin `command`; `-v` asks how each name resolves; remaining words are command names | Paths such as `/opt/homebrew/bin/go`; missing names have no path, and a nonzero status is possible |
| `brew --version` | Homebrew executable; `--version` requests its version | `Homebrew <version>`; checks the package manager, not all available packages |
| `GOTOOLCHAIN=local go version` | The assignment applies to this command only; `go` invokes the SDK; `version` reports it | `go version go<version> darwin/arm64`; `local` prevents automatic toolchain selection/download for this check |
| `GOTOOLCHAIN=local go env GOOS GOARCH` | `env` is a Go subcommand; `GOOS` and `GOARCH` are the selected target OS and CPU names | Two lines: `darwin`, then `arm64`. These are build targets and may differ if you previously set overrides |
| `git --version` | Version of the source-control command | `git version <version>`; no GitHub login or network access needed |
| `podman --version` | Version of the macOS Podman client | `podman version <version>`; does not prove its Linux engine is reachable |
| `podman machine list` | `machine` manages local Linux VMs; `list` reports existing VMs | Table of names, state, CPU, and memory; a stopped/absent VM is acceptable for Phase 1 |

Go can select another toolchain based on configuration and a module's requirements. The per-command `local` setting keeps the inventory about the installed SDK. It does not change your shell profile or Go configuration. See the [official toolchain explanation](https://go.dev/doc/toolchain).

## Expected output on this machine

The assessment captured on **17 October 2026** included:

```text
ProductName: macOS
ProductVersion: 26.5.2
arm64
Homebrew 7.0.9
go version go1.27.1 darwin/arm64
darwin
arm64
git version 2.50.1 (Apple Git-155)
podman version 6.1.3
```

The machine list showed `podman-machine-default` already running. This is an observation, not an instruction to allocate another VM. Versions and spacing can differ on a later run; interpret the fields rather than matching the whole output. Full findings are in the [environment assessment](../environment-assessment.md).

## Java comparison

| Familiar Java idea | Go/macOS equivalent | Important difference |
| --- | --- | --- |
| JDK with `java` and `javac` | Go SDK with `go` and `gofmt` | `go` provides build, run, test, and dependency commands; no JVM is required to run a compiled Go binary |
| IntelliJ project SDK | GoLand SDK or VS Code's selected Go executable | The editor can select a different SDK from your terminal; compare the path and version |
| Maven/Gradle executable on `PATH` | `go` on `PATH` | Module management is built into Go; module syntax comes in Phase 1 |
| A Linux VM running containers | Podman Machine on macOS | The macOS CLI sends requests to the engine inside a Linux VM |
| A Spring Boot process on your Mac | A Go service launched by your IDE | The process uses its launch environment; it does not receive later Terminal variable changes |

Podman's documentation explains the [macOS client and VM arrangement](https://podman.io/docs/installation). The database and cache will run inside that VM's containers later. Your editor and Go debugger will normally run directly on macOS.

## Hands-on exercise and debugging

Complete [EXERCISES.md](EXERCISES.md), then the temporary missing-`PATH` exercise in [DEBUGGING.md](DEBUGGING.md). Record your interpretation, not just command output. No Go syntax is introduced during this inventory lesson.

Common mistakes are treating a missing launcher as proof an IDE is absent, assuming an installed Podman client proves a running engine, and installing Intel binaries because a tutorial used `amd64`.

In a real microservice repository, the same checks help answer: Which Go version builds this service? Is it running on the host or in a container? Which hostname belongs to that runtime? Those questions come before changing application code.

## Five knowledge checks

1. What do `darwin` and `arm64` tell you? Why would a container use `linux/arm64`?
2. Why can `command -v goland` fail even when GoLand opens from Applications?
3. What does `GOTOOLCHAIN=local` change, and how long does it last here?
4. Why does `podman --version` succeeding not guarantee containers can run?
5. Which tools must work before your first Go lesson, and which can wait?

Answer first, then compare with [INTERVIEW.md](INTERVIEW.md).

## Recap and completion criteria

The terminal locates tools; the Go SDK builds programs; the editor assists development; Podman delegates Linux-container execution to a VM. Checking one layer does not establish the health of the others.

- [ ] Explain the core command output in your own words.
- [ ] Confirm a native Go SDK, Git, and one usable editor; review required updates in [SETUP.md](../../SETUP.md).
- [ ] Complete the temporary `PATH` debugging exercise and recover the original shell behavior.
- [ ] Record confidence and gaps in [PROGRESS.md](../../PROGRESS.md).
- [ ] Confirm the environment is understood and ready before starting Phase 1.

**Next topic, after confirmation:** Phase 1, lesson 1 — Go SDK, a first small module, `package main`, imports, `func main`, variables, and the preserved welcome program, with an automated behavior test.

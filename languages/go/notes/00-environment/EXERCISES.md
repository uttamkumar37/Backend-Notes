# Session 0 exercises

Allow 15–25 minutes after reading the lesson. These exercises inspect your own Mac and need no office access.

## 1. Build an environment map

Run the command block in [README.md](README.md). Write down your macOS version, architecture, Go version, Go build target, Git version, and Podman machine state. Explain what each value proves and one thing it does not prove.

**Verification:** distinguish a macOS Go binary from a Linux container binary, and do not infer engine readiness from the Podman client version.

## 2. Choose your first editor

Open VS Code or GoLand through Applications. Find where it displays/configures the Go SDK. Compare that location with:

```bash
command -v go
```

`command` is a shell builtin and `-v` reports which executable the current shell selects. The editor may use a separately selected SDK. Record any difference; no editor configuration change is required for this exercise.

**Verification:** locate the SDK setting and identify whether terminal and editor agree. A missing `goland` terminal launcher is not a failed installation check.

## 3. Sort tools by when they are needed

Using [SETUP.md](../../SETUP.md), list the tools needed now, at PostgreSQL/cache stages, and at Kubernetes deployment. State which processes run directly on macOS and which run inside Linux containers.

**Verification:** the Phase 1 list does not require Kubernetes, Vault, PostgreSQL, or a running Podman VM.

## 4. Explain before changing anything

Complete the `PATH` debugging drill in [DEBUGGING.md](DEBUGGING.md). Answer the five questions in [INTERVIEW.md](INTERVIEW.md) before reading their explanations.

**Evidence to record:** interpretation of your inventory, chosen editor, drill cause and fix, and any remaining error. Update only your own progress; running an inventory is not completion of Go fundamentals.

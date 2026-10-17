# Session 0 knowledge checks with explanations

Try answering before reading the explanations. These check setup understanding rather than memorized versions.

## 1. What do `darwin/arm64` and `linux/arm64` mean?

The first component identifies the target OS and the second the CPU architecture. Your native Mac uses Darwin; a Linux container uses Linux. An executable also depends on its OS ABI, so the processor architecture alone does not make it portable. In Java terms, distinguish the platform running the JVM from the portable bytecode it interprets or compiles.

## 2. Does a missing `goland` command prove GoLand is not installed?

No. `command -v` searches the shell's `PATH`; a GUI application may be installed without a launcher. Open it and inspect the Go SDK setting. The same distinction applies to Postman, Bruno, and other desktop tools.

## 3. Why use `GOTOOLCHAIN=local go version` here?

It reports the installed SDK without selecting/downloading another toolchain. The assignment immediately before the command is scoped to that process; it does not persist in shell configuration. Module-driven toolchain selection comes later. See the [Go toolchain documentation](https://go.dev/doc/toolchain).

## 4. Where do Podman's containers run on a Mac?

Inside a Linux VM managed by Podman Machine. The macOS client sends commands to that engine. An installed client and reachable engine are separate checks. Later, service connectivity debugging also requires distinguishing host addresses from container-network addresses.

## 5. What is the minimum environment for the first Go lesson?

A working native Go SDK, terminal, Git for this repository, and one editor that can use the SDK. Go includes its testing tools. PostgreSQL, Valkey, Vault, Kubernetes, Helm, and dashboards arrive when their lessons need them. Copilot is optional; you must still explain, test, and review code yourself.

**Pass criterion:** explain all five independently and diagnose the `PATH` exercise. Copying answers does not establish confidence.

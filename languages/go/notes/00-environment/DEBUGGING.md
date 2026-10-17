# Session 0 debugging

## Drill — installed Go suddenly disappears

**Symptom:** a command works in one shell but another says it cannot find it.

**Reproduce safely:** from your normal Terminal, run this command. It changes `PATH` only inside a temporary child shell, leaving the parent untouched.

```bash
/usr/bin/env PATH=/usr/bin:/bin /bin/zsh -f -c 'command -v go || printf "go is outside this temporary PATH\n"'
```

**Explain the syntax:** `/usr/bin/env` launches a process with the supplied `PATH`; `:` separates directories. `/bin/zsh` starts a shell; `-f` skips normal user startup files; `-c` executes the quoted text. `command -v go` searches for Go. `||` runs the fallback `printf` if the search fails. Single quotes preserve the text for the child shell. Go in `/opt/homebrew/bin` is deliberately outside this search path. A separate `/usr/bin/go`, if present, would change the result.

**Expected output for this machine:**

```text
go is outside this temporary PATH
```

**Root cause:** the child shell cannot search the Homebrew directory. Go has not been removed.

**Fix and verify:** return to the original shell, which the command already does, then run:

```bash
command -v go
GOTOOLCHAIN=local go version
```

Expect `/opt/homebrew/bin/go` and a `darwin/arm64` version line. `GOTOOLCHAIN=local` applies only to that invocation. If these still fail, compare your installation with [Homebrew's supported prefix and shell setup](https://docs.brew.sh/Installation) before editing a profile. Do not replace the full `PATH` in your real shell.

## Other setup failures

| Symptom | Diagnose | Likely cause | Fix | Verify |
| --- | --- | --- | --- | --- |
| `go` missing in an IDE terminal only | Compare `command -v go` in both terminals and inspect IDE SDK | Different inherited environment or SDK selection | Reopen IDE after configuring SDK/shell environment | Both terminals report the intended Go installation |
| `uname -m` says `x86_64` | Check whether Terminal opens with Rosetta | Intel process on Apple Silicon | Open native Terminal and choose ARM64 downloads | `uname -m` reports `arm64` |
| Unexpected `go env GOOS GOARCH` | Check whether those two named variables were set in shell/IDE | Leftover cross-compilation configuration | Remove only unintended overrides | Target becomes `darwin` / `arm64` |
| Podman client works, containers fail | `podman machine list` | VM absent/stopped or wrong connection | Follow [SETUP.md](../../SETUP.md) when containers are needed | VM running, then `podman info` succeeds |
| `psql` missing after `libpq` install | `brew --prefix libpq` | Keg-only directory outside `PATH` | Use full executable path or documented temporary path addition | `psql --version` succeeds |
| `goland` absent from `PATH` | Open GoLand from Applications | Optional launcher unconfigured | Use GUI; add launcher later if desired | Editor opens and recognizes SDK |

Record the symptom, evidence identifying its layer, and command proving recovery. Avoid reinstalling everything in response to one missing command.

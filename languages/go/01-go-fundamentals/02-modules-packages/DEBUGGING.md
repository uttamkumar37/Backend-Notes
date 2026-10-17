# Lesson 2 debugging

Introduce one local fault at a time. Commands below name their working directory explicitly. Restore only the line/file you changed; broad Git restore commands can discard your other uncommitted notes. Assistant verification uses temporary copies for these drills.

## Drill 1 — a function becomes unexported

**Working directory:** `02-modules-packages/`.

**Reproduce:** change the library declaration from `func Welcome(name string) string` to `func welcome(name string) string`. Change the call in `main.go` from `greeting.Welcome(investorName)` to `greeting.welcome(investorName)`. Leave the test file alone.

These are the same declarations/calls explained in README, with the initial letter changed. `welcome` still has a string parameter and string result, but it is now unexported.

```bash
GOTOOLCHAIN=local GOWORK=off go run .
```

**Symptom:** compilation fails at the lowercase call, with wording such as `name welcome not exported by package greeting` or `undefined: greeting.welcome`, depending on toolchain diagnostics.

**Root cause:** package visibility. Finding/importing the correct package does not make all its functions callable. Lowercase functions remain usable inside their own package.

**Fix:** restore uppercase `Welcome` in the declaration and call. Do not “fix” it by moving every file into `package main` or adding another `go.mod`.

**Verify:** `go run .` prints `Welcome, Asha.` and `go test -v ./...`, with the same environment prefixes, passes both packages. The external package test also depends on the exported API.

## Drill 2 — a missing checksum under read-only module mode

**Working directory:** `02-modules-packages/`. Keep an editor copy of the two supplied `go.sum` lines before changing the file.

**Reproduce:** remove only the first checksum line, the one for module content without `/go.mod`. Keep the metadata entry. Run:

```bash
GOTOOLCHAIN=local GOWORK=off go test -mod=readonly ./...
```

`-mod=readonly` prevents automatic module-file updates for this build/test. The library test now reports a missing `go.sum` entry for the package supplying `cmp`; exact surrounding diagnostics may vary. The program's runtime uses only local/standard packages, so its package may still pass while the library test setup fails.

**Root cause:** a required test dependency lacks its recorded content checksum. This is not evidence of an absent `Welcome` function or broken SDK. Under writable module mode Go can often regenerate missing entries; the read-only flag makes the missing metadata explicit.

**Fix:** from the module root, run `GOTOOLCHAIN=local GOWORK=off go mod tidy`. It uses cached/public dependency metadata to restore generated checksums. Do not type a guessed hash or disable checksum verification.

**Verify:** compare `go.sum` with the supplied two entries, rerun the read-only test command, and run `go mod verify` with the same prefixes. Both succeed. The [dependency management guide](https://go.dev/doc/modules/managing-dependencies) explains generated dependency files.

## Drill 3 — workspace path points to the wrong directory

**Working directory:** `02-modules-packages/workspace/`.

**Reproduce:** change `../../01-first-program` in `go.work` to `../../missing-lesson`. Then run:

```bash
GOTOOLCHAIN=local GOWORK="$PWD/go.work" go test ../... ../../01-first-program/...
```

**Symptom:** Go fails while loading the workspace and reports that the member's `go.mod` cannot be read.

**Root cause:** a `use` entry points to a nonexistent local module directory. Paths are relative to `go.work`, not relative to your Git repository root or a module's import path. Do not download a dependency to compensate for a misspelled local path.

**Fix:** restore `../../01-first-program` in `go.work`.

**Verify:** the same workspace test command succeeds, then run `cd ..` and `GOTOOLCHAIN=local GOWORK=off go test ./...` to establish independent module behavior.

## Common failures

| Symptom | Diagnose | Cause | Fix | Verify |
| --- | --- | --- | --- | --- |
| `no required module provides package .../greetings` | Compare import with `go.mod` prefix and actual `greeting/` directory | Misspelled local import | Correct the path; do not `go get` a fictional local path | Run/list/tests succeed |
| `found packages ... in ...` | Check package declarations in that directory | Incompatible production package names mixed together | Keep command and library in their own directories; external `_test` package is a supported special case | `go list ./...` reports the two intended packages |
| Editor shows unresolved imports but terminal tests pass | Open the exact lesson module folder or configured workspace | Editor is operating with a different project/module context | Select correct folder/SDK/workspace and refresh its project model | Navigation reaches local `Welcome`; terminal and editor agree |
| `go test ./...` from the workspace folder does not test members | Inspect `pwd`, member paths, and `go env GOWORK` | Pattern describes directories under the workspace folder, not all `use` entries | Use the explicit two member package patterns | Both lessons' test names appear |
| `go mod verify` passes but tests fail | Read the actual failing assertion/compiler diagnostic | Cache integrity and program behavior are different checks | Fix the identified logic/source issue | Relevant behavior test passes |

The [workspace tutorial](https://go.dev/doc/tutorial/workspaces) describes workspace membership; the [go command reference](https://pkg.go.dev/cmd/go) describes package patterns and module-mode flags. Preserve exact diagnostics and working directory when asking for help.

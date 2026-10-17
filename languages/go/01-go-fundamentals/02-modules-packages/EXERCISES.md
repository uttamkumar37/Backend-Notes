# Lesson 2 exercises

Allow 25–35 minutes after reading the lesson. Start in `02-modules-packages/`. Prefix independent Go commands with `GOTOOLCHAIN=local GOWORK=off`, as explained in README. Keep exercise faults isolated and restore the intended working source.

## 1. Map the import before executing

Open `go.mod`, `main.go`, and `greeting/greeting.go`. Write the module path, library subdirectory, complete library import path, declared package name, exported function name, and argument/result values.

Run `go list ./...` with the prefixes above and compare the two printed package paths with your prediction.

**Verify:** you can locate the function without searching a remote repository or confusing `package greeting` with the full import path. Explain why there is only one `go.mod` for these two production packages.

## 2. Observe package test boundaries

Run these from the module root:

```bash
GOTOOLCHAIN=local GOWORK=off go test -v ./greeting
GOTOOLCHAIN=local GOWORK=off go test -v .
GOTOOLCHAIN=local GOWORK=off go test -v ./...
```

`./greeting` selects the library package; `.` selects the command package; `./...` includes both. The command example calls the library, but running its tests does not automatically execute the library's own test functions.

**Verify:** identify `TestWelcomeUsesInvestorName` and `Example` in the appropriate runs. Explain why `greeting_test` must import the package under test and can call only exported names.

## 3. Inspect and reconcile the dependency files

Read `go.mod` and `go.sum`, then run:

```bash
GOTOOLCHAIN=local GOWORK=off go list -m all
GOTOOLCHAIN=local GOWORK=off go mod download
GOTOOLCHAIN=local GOWORK=off go mod tidy
GOTOOLCHAIN=local GOWORK=off go mod verify
```

The first download may use the public network; already-cached modules need not download again. Compare files in your editor before/after. Since the supplied source already declares its dependency, tidy should leave these files unchanged. Do not change the version simply to produce a diff.

**Verify:** explain the difference between package path `github.com/google/go-cmp/cmp` and module path `github.com/google/go-cmp`; distinguish content and `/go.mod` checksum entries. Explain why importing a dependency only from tests does not remove its requirement.

## 4. Work with two modules

From the Lesson 2 root:

```bash
cd workspace
GOTOOLCHAIN=local GOWORK="$PWD/go.work" go env GOWORK
GOTOOLCHAIN=local GOWORK="$PWD/go.work" go test -v ../... ../../01-first-program/...
cd ..
GOTOOLCHAIN=local GOWORK=off go test -v ./...
```

**Verify:** the workspace run includes both lessons; the independent run includes only Lesson 2. Each member still has its own module path. Explain every `..` from the workspace file's location. No workspace file is placed at the enclosing repository root.

## 5. Diagnose a module/package failure

Complete one drill in [DEBUGGING.md](DEBUGGING.md), record its actual error and cause, restore the fault, and rerun the affected tests. Then answer the five questions in [INTERVIEW.md](INTERVIEW.md).

Record your confidence and remaining gaps in [PROGRESS.md](../../PROGRESS.md). Passing supplied tests and independently explaining the import/dependency/workspace relationships are separate completion criteria.

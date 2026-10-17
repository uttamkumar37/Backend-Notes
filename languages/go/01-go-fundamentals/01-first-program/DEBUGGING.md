# Lesson 1 debugging

Run all commands from this lesson's module directory. Introduce one fault at a time and manually undo only your change afterward; avoid repository-wide restore commands that could erase other notes.

## Main drill — correct types, wrong argument order

**Reproduce:** in `welcomeMessage`, change its formatter call to this line:

```go
return fmt.Sprintf("Learning %s for %s.", goal, language)
```

This replaces the same return statement from README. All syntax is already introduced: `return` sends back the string from `fmt.Sprintf`; the first `%s` now consumes `goal`, and the second consumes `language`. Both arguments are still strings, so the types are valid.

**Symptoms and diagnostics:**

```bash
GOTOOLCHAIN=local go run .
GOTOOLCHAIN=local go test -v ./...
```

The program's first line becomes `Learning backend engineering for Go.`. The unit-test failure should include:

```text
welcomeMessage() = "Learning Investor Platform for Java."; want "Learning Java for Investor Platform."
```

The output example also fails because it expects the original complete output. `go test` returns nonzero; line numbers/timings in the diagnostic may vary.

**Root cause:** the string formatting arguments were swapped. Compilation validates the types used here; it cannot infer the intended meaning of two valid strings. Reinstalling Go or adjusting `go.mod` cannot fix this logic error.

**Fix:** put `language` first and `goal` second, as in the supplied implementation. Read the diff to confirm only your intended line changed.

**Recovery verification:**

```bash
GOTOOLCHAIN=local go test -v ./...
GOTOOLCHAIN=local go run .
```

Both checks pass and the two original output lines return. Explain the formatter's inputs before moving on. Assistant verification of this drill uses a temporary copy, leaving the supplied working module intact.

## Additional small failures

| Fault to try / symptom | Diagnose | Root cause | Targeted fix | Verify |
| --- | --- | --- | --- | --- |
| Use `language = "Go"` as the first declaration; undefined `language` | Run `go run .`; read the first compiler diagnostic and its line | `=` assigns an existing name, but none was declared | Restore `language := "Go"` | Program compiles and tests pass |
| Add a second `goal := "API development"` in the same block; no new variables | Read the reported declaration line | Single-name short declaration introduces no new variable | Use `goal = ...` for reassignment, then restore the supplied value after observing output | Compilation succeeds; update expectations only for intentional behavior changes |
| Leave an import unused; import error | Read the import named in the diagnostic | Go rejects unused imports in a file | Remove the unused import; keep `fmt` because this file uses it | `go run .` compiles |
| Run `go run .` from the enclosing `languages/go/` instead of here; missing main module | Run `pwd` and confirm the directory containing this `go.mod` | Current directory is outside the lesson module | `cd` to `01-go-fundamentals/01-first-program` | `go run .` and module-scoped tests succeed |
| Omit `// Output:` from `Example`; example no longer checks the output contract | Compare the function with the supplied test file and the verbose test list | Output-checking examples need the expected-output marker | Restore the marker and both expected lines | `go test -v` executes `Example`, which fails when output is deliberately wrong |
| Edit source but run the old `bin/first-program`; stale output | Compare `go run .` with the existing executable | Built file still contains the previous program | Rebuild using README commands | Both executions match the intended current source |

`pwd` prints the current directory. The `go` commands use the installed SDK through the per-command `GOTOOLCHAIN=local` setting. Example-test behavior is documented in the [testing package reference](https://pkg.go.dev/testing#hdr-Examples).

## Optional editor walkthrough

Open this lesson folder in your chosen editor. If the Go debugger is configured, put a breakpoint at `message := welcomeMessage(language, goal)`. Start debugging the program (not the test), inspect `language` and `goal`, step into the function, then step back to see `message` assigned. Do not expect `message` to have its final value before that call completes.

If the editor requests tool installation, keep the exact message; the terminal exercises work independently. A full Delve session comes later. For now, tracing a small program and reading a test diagnostic is the required debugging practice.

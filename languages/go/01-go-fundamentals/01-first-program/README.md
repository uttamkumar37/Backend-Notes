# Phase 1 · Lesson 1 — Your first Go module and tested program

**Time:** 60–90 minutes. **Current machine:** Apple Silicon Mac, installed Go 1.27.1. **Prerequisites:** [environment assessment](../../notes/00-environment/README.md), a terminal, and one editor. Your instruction to continue authorizes starting Phase 1; it does not mark earlier exercises or this phase's assessment complete.

## 1. Today's objective

Run this small program, trace its execution from `main`, explain every line, and use an automated test to catch incorrect behavior. You will learn `package`, `import`, `func`, string variables, `:=`, function arguments, return values, and a first test. One lesson is enough today; database/framework/infrastructure setup comes later.

By the end you should be able to change the learning goal, predict the result, and distinguish a compiler error from a failed behavior test.

## 2. Prerequisite concepts

A **source file** contains code such as `main.go`. A **compiler** translates it into machine code. A **package** groups Go source files that build together. A **module** gives one or more packages a shared identity and dependency configuration through `go.mod`. Your working directory determines which module/package a terminal command uses.

You already know methods, arguments, strings, and test assertions from Java. Go uses those ideas with different syntax and explicit package/function boundaries. You do not need Spring's application context or a JVM to run this compiled executable.

## 3. What → why → when → how

| Concept | What | Why / when | How it appears here |
| --- | --- | --- | --- |
| SDK | Compiler, runtime support, and development commands | Build/run/test any Go application | Installed `go` command |
| Module | Named collection of packages with dependency metadata | Establish a project boundary | This directory's `go.mod` |
| Package | Source files compiled together | Organize related functions/types | `package main` in both `.go` files |
| Entry point | Function invoked for application execution after initialization | Start a command or service | `func main()` |
| Variable | Named value with a type and scope | Carry data through your program | `language`, `goal`, `message` |
| Function | Named operation with inputs and optional results | Reuse behavior and test it independently | `welcomeMessage` creates a string |
| Test | Automated comparison with expected behavior | Detect a regression after changing code | Argument test and executable example |

The [official getting-started tutorial](https://go.dev/doc/tutorial/getting-started) explains the basic module/program workflow. Here, all files are already supplied: open this lesson rather than running an initializer again.

### This module's files

```text
01-first-program/
├── go.mod
├── main.go
├── main_test.go
├── .gitignore
├── README.md
├── EXERCISES.md
├── DEBUGGING.md
└── INTERVIEW.md
```

`go.mod` contains:

```go
module example.com/go-lab/first-program

go 1.27.0
```

This is **module configuration syntax**, not executable Go source. `module` declares an example module path; it is not a request to contact that website. `go 1.27.0` declares the minimum Go version and language-version baseline. Installed Go 1.27.1 satisfies it. This declaration does not pin every build to patch 1.27.0. See [the module reference](https://go.dev/ref/mod#go-mod-file-go).

There are no third-party dependencies, so no `require` entries or `go.sum` file are needed. Standard packages come with Go. In a brand-new directory you would create a module with `go mod init example.com/go-lab/first-program`; that has already been done for this lesson. Running it here again reports that `go.mod` already exists. Lesson 2 explores module paths, checksums, packages, and workspaces further.

`.gitignore` contains `/bin/`, which tells Git to ignore build outputs in this module's `bin` folder. It does not affect compilation or tests.

## 4. Java / Spring Boot comparison

| Java idea | Here in Go | Where the analogy stops |
| --- | --- | --- |
| `public static void main(String[] args)` | `func main()` in `package main` | Go's entry function has no parameters or return values; arguments are obtained separately when needed |
| Class with static helper methods | Package-level functions | No class or object is required to call these functions |
| `String language = "Go"` or local `var` | `language := "Go"` | The type is inferred as `string`; the variable remains statically typed |
| `String.format(...)` | `fmt.Sprintf(...)` | It returns text; it does not print |
| `System.out.println(...)` | `fmt.Println(...)` | Go's printing function also returns a byte count and possible write error |
| JUnit assertion | `if got != want { t.Errorf(...) }` | Go's built-in test runner discovers named test functions; no assertion library is needed |
| Spring Boot application startup | This program's `main` | We currently have ordinary function calls, with no server, dependency injection container, or framework lifecycle |

## 5. Small runnable example

Open [main.go](main.go). It adapts your preserved [welcome program](../../practice/01-basics/hello_backend.go) by returning the welcome text before printing it. That makes the message independently testable.

```go
package main

import "fmt"

func main() {
	language := "Go"
	goal := "backend engineering"
	message := welcomeMessage(language, goal)

	fmt.Println(message)
	printNextSteps()
}

func welcomeMessage(language string, goal string) string {
	return fmt.Sprintf("Learning %s for %s.", language, goal)
}

func printNextSteps() {
	fmt.Println("Next: practice variables, functions, structs, and errors.")
}
```

## 6. Explain every line and value

Line numbers below refer to `main.go`; blank lines separate ideas and perform no operation.

| Line | Explanation |
| --- | --- |
| 1 | `package` identifies the package for this file. `main` is the special package name used to build an executable command. It is not a class. The file could have another name; the entry point is determined by package/function, not by the filename `main.go`. |
| 3 | `import` makes another package available to this source file. `"fmt"` is the standard formatting package; its quoted import path is not a variable. We need it to build and print strings. |
| 5 | `func` declares a function. `main` is its name, `()` says it accepts no parameters, and `{` begins its body. Nothing follows `()` before `{`, so it declares no return values. |
| 6 | `language` is a new local variable. `:=` declares it and assigns `"Go"`. The double quotes make a string literal; Go infers type `string`. This value originates in source code, not an environment variable or a request. |
| 7 | Another short declaration creates local `goal` with type `string`. Its literal value is `"backend engineering"`. |
| 8 | `welcomeMessage(language, goal)` calls our function. The comma separates two **arguments**: current values of `language` and `goal`. Its returned string initializes a new local `message` using `:=`. The call must finish before `message` is available. |
| 10 | The dot selects exported `Println` from imported package `fmt`. Parentheses call it; `message` is its argument. It writes that text to standard output and adds a newline. No variable is declared here. |
| 11 | Calls our `printNextSteps` function with no arguments. It prints its own line and returns no values. |
| 12 | `}` ends `main`. Once it returns, this simple application's execution finishes. There are no background tasks. |
| 14 | Declares `welcomeMessage`. `language string` and `goal string` are two **parameters**: input names followed by their types. The final `string` after `)` is the single return type. These local parameter names are distinct from `main`'s variables; values are passed into the call. |
| 15 | `fmt.Sprintf` formats a new string. The first `%s` uses `language`; the second uses `goal`, in that order. Spaces, the word `for`, and the final period come from the literal. `return` sends the string to the caller and ends this function invocation. It does not write to the terminal. |
| 16 | Ends `welcomeMessage`. A function may be declared below its caller; it does not execute just because its definition appears in the file. |
| 18 | Declares a no-parameter, no-result function named `printNextSteps`; its initial lowercase letter keeps it unexported outside this package. `welcomeMessage` is also unexported. |
| 19 | Writes the literal next-step text and a newline using `fmt.Println`. The text originates right here. |
| 20 | Ends `printNextSteps`; execution resumes after its call, then `main` returns. |

`{}` enclose function or conditional bodies, `()` declare parameters or call functions, and commas separate arguments. Go code here has no explicit line-ending semicolons; ordinary line breaks are sufficient. Local variables have scope within their enclosing function/block. This is ordinary sequential execution; no goroutines or parallelism are involved.

### Why `:=` here, and when `=` would be used

The three names in `main` are introduced for the first time, so they use `:=`. Later assignment to the same variable uses `=`, for example inserting `goal = "cloud-native services"` **after** the declaration and **before** creating `message`. The name already exists, so that line changes its value instead of declaring another variable. `message` must be recomputed if you want a later goal change reflected in it.

For a single name already declared in the same block, another `goal := ...` fails because it declares no new variable. The fuller multiple-name/shadowing rules come later. `:=` is allowed within functions; package-level declarations use other syntax taught in subsequent lessons. See the [Go specification's short declarations](https://go.dev/ref/spec#Short_variable_declarations).

### Follow execution and returns

```text
go run . compiles the package, then launches the resulting program
  main()
    language = "Go"; goal = "backend engineering"
    welcomeMessage(language, goal)
      returns "Learning Go for backend engineering."
    message receives that returned value
    fmt.Println(message) writes line 1
    printNextSteps() writes line 2
    main returns; process ends
```

The assignments above describe values; they are an execution sketch rather than more Go source. Passing strings gives the called function its own parameter values. There is no pointer receiver, `&` address operation, or `*` dereference in `main.go`. Creating your own pointer-based types is not a prerequisite today.

`welcomeMessage` returns a `string` and cannot return an `error` under its current signature. `fmt.Sprintf` also returns only a string. `fmt.Println` returns a byte count and an `error` if its output write fails; these calls intentionally ignore both for a tiny terminal demonstration. Go does not automatically throw an exception when a function returns an error. This output-only simplification is not the error-handling design for the future service; explicit errors and I/O handling get their own lessons. The [fmt reference](https://pkg.go.dev/fmt) documents formatting and printing behavior.

### Understand the automated tests

Open [main_test.go](main_test.go). Go uses the `_test.go` suffix and function naming/signature conventions to discover tests. Its testing support is built in; see [the official testing tutorial](https://go.dev/doc/tutorial/add-a-test).

```go
package main

import "testing"

func TestWelcomeMessageUsesArguments(t *testing.T) {
	got := welcomeMessage("Java", "Investor Platform")
	want := "Learning Java for Investor Platform."

	if got != want {
		t.Errorf("welcomeMessage() = %q; want %q", got, want)
	}
}

func Example() {
	main()
	// Output:
	// Learning Go for backend engineering.
	// Next: practice variables, functions, structs, and errors.
}
```

| Test-file line | Explanation |
| --- | --- |
| 1 | Same `main` package as the implementation, so tests can call its unexported functions. A test file is not part of the normal application build. |
| 3 | Imports standard package `testing`; it supplies the test state/reporting type. Imports apply to each file, so this file does not inherit the other file's `fmt` import. |
| 5 | `func` declares a runner-discovered test named `TestWelcomeMessageUsesArguments`. It accepts `t`, provided by Go's test runner, and returns nothing. `*testing.T` means a **pointer to a testing.T value**: `*` is part of the type here, not multiplication or a dereference expression. The pointer lets methods report results to the same test state. We do not create it with `&`; the runner supplies it. |
| 6 | `got :=` declares a string from the actual function result. Test values `"Java"` and `"Investor Platform"` deliberately differ from `main`'s defaults, so hardcoding the default welcome cannot pass this check. |
| 7 | `want :=` declares the independently specified expected string. These are test values, not user account data. |
| 9 | `if` runs its block only when `got != want`. `!=` means not equal and produces a boolean. Parentheses around the condition are not required; braces are. |
| 10 | `t.Errorf` is a method call on the runner's test value. It formats a failure and marks this test failed; it does not panic or immediately stop execution. `%q` shows quoted/escaped strings, making missing spaces/newlines visible. `got` and `want` fill those placeholders in order. |
| 11–12 | Close the conditional and the test. If the values match, no failure is reported. |
| 14–15 | `Example` is an executable package example with no parameters/results; it calls our `main` so both output lines are checked. Go supplies a separate test runner entry point; `go test` does not automatically execute application `main`. |
| 16–18 | `//` starts a comment. The exact `Output:` marker lets Go's example runner compare captured standard output with these expected lines. They are not extra print calls. With this marker, incorrect output fails the example. |
| 19 | Ends the example. |

Both `testing.T` usage and executable examples are described in the [testing package reference](https://pkg.go.dev/testing). These are behavior tests: one checks dynamic inputs, the other the user's complete two-line output. They do not claim input validation, HTTP behavior, or production readiness.

## 7. macOS terminal commands

Open Terminal and enter this exact lesson directory:

```bash
cd /Users/uttamkumar/Uttam-workspace/GitHubProject/Backend-Notes/languages/go/01-go-fundamentals/01-first-program
GOTOOLCHAIN=local go version
GOTOOLCHAIN=local go run .
GOTOOLCHAIN=local go test -v ./...
GOTOOLCHAIN=local go vet ./...
```

| Command part | Meaning |
| --- | --- |
| `cd .../01-first-program` | Select the directory containing this module's `go.mod`. It is not the enclosing repository root. |
| `GOTOOLCHAIN=local` | Use the installed toolchain for just that command; no persistent environment change or automatic SDK download. |
| `go version` | Display the installed SDK and native target. |
| `go run .` | Compile and launch the current command package. `.` denotes the current directory; test files are excluded. |
| `go test -v ./...` | Build/run tests, printing their names. `./...` selects package directories under this module; it does not cross into every other module. |
| `go vet ./...` | Check for suspicious code patterns; it does not replace running the behavior tests. |

For a persistent native binary, run from the same directory:

```bash
mkdir -p bin
GOTOOLCHAIN=local go build -o bin/first-program .
./bin/first-program
```

`mkdir -p bin` creates the output directory if needed. `build` compiles without running. `-o` chooses the output filename. `./bin/first-program` runs that built executable. It prints the same two lines. Go's runtime is included in the executable; no JVM is needed. This build targets macOS/ARM64 here, not Linux containers. The ignored `bin/` directory holds generated files only.

## 8. Expected terminal output

`go run .` and the built binary print exactly:

```text
Learning Go for backend engineering.
Next: practice variables, functions, structs, and errors.
```

`go test -v ./...` should contain:

```text
=== RUN   TestWelcomeMessageUsesArguments
--- PASS: TestWelcomeMessageUsesArguments (0.00s)
=== RUN   Example
--- PASS: Example (0.00s)
PASS
ok      example.com/go-lab/first-program    0.123s
```

The elapsed durations, alignment, and possible `(cached)` marker vary. `go vet` and a successful `go build` normally print nothing. If compilation fails, these expected program/test lines do not appear; read the first diagnostic before editing.

## 9. Hands-on exercise

Complete [EXERCISES.md](EXERCISES.md): trace values, change the goal, keep expected output in sync with intentional behavior, and prove the built executable only changes after rebuilding.

## 10. Debugging challenge

Use [DEBUGGING.md](DEBUGGING.md) to swap the formatter's argument order, reproduce a failing test, explain why the compiler accepts it, and repair it. Also learn to recognize a missing module, an undefined variable, and an invalid repeated declaration. Leave the module passing when finished.

## 11. Five knowledge checks

1. How do a module, a package, and `main` differ? Does the filename alone make an executable?
2. Why is `message := welcomeMessage(language, goal)` valid, but assigning a never-declared `message` with `=` invalid?
3. Which function creates text, which operation prints it, and where do the two inputs originate?
4. What does `t *testing.T` mean, and how does a failed comparison make `go test` fail?
5. Why can swapped string arguments compile yet fail tests, and how does `go run .` differ from `go build`?

Answer before reading [INTERVIEW.md](INTERVIEW.md).

## 12. Recap and microservice connection

Execution enters `main`, which creates local values, calls functions, receives a result, and writes output. The later Investor service will use the same explicit flow to load configuration, construct dependencies, register handlers, and start listening. Returning a value separately from producing output is our first small step toward testable boundaries.

## 13. Completion criteria

- [ ] Run the program and explain both output lines.
- [ ] Explain all three `:=` declarations and the inputs/result of `welcomeMessage`.
- [ ] Run tests and explain what each check protects.
- [ ] Build/run the native executable and explain stale output before rebuilding.
- [ ] Diagnose the swapped-argument failure and return the code/tests to passing.
- [ ] Answer the five questions independently and record confidence/gaps in [PROGRESS.md](../../PROGRESS.md).

Assistant verification establishes that the supplied code works; these checkboxes remain for your own practice. This lesson does not complete Phase 1.

## 14. Next topic

After your confirmation, lesson 2 will explore modules, packages, imports, `go.mod`/`go.sum`, and a small Go workspace. First complete or discuss today's exercises; if a command fails, keep its exact diagnostic so we can debug this implementation.

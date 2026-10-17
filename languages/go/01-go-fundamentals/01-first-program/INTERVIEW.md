# Lesson 1 — Five knowledge checks with explanations

Answer in your own words before opening the explanations. Use the actual program to support your answers rather than memorizing a definition.

## 1. How do a module, a package, and `main` differ?

A module is the project/dependency boundary identified by `go.mod`. Packages group source files for compilation. An executable command uses `package main` and a `func main()` entry point. Naming a file `main.go` alone is insufficient: a file in another package does not become an executable by its filename. The application entry point here is distinct from the test runner Go builds for `go test`.

**Java bridge:** a Go package is not a class, and a module is not the same concept as Java's `module-info.java`. Compare project/dependency organization and source grouping separately.

## 2. Why does the initial declaration use `:=` rather than `=`?

`message` does not yet exist in `main`, so `:=` declares it and initializes its value, with type `string` inferred from the function result. `=` would require a previously declared assignable variable. Later, `goal = "API development"` assigns to the existing `goal`. A single-name `goal := ...` repeated in the same block fails because there is no new variable. Later lessons cover multiple-name declarations and nested scopes.

**Common mistake:** calling inferred typing dynamic typing. The variable still has a static type.

## 3. Which function creates text, which operation prints it, and where do inputs originate?

`main` supplies the string literals `"Go"` and `"backend engineering"`. `welcomeMessage` receives those argument values as its parameters and returns the text built with `fmt.Sprintf`. It does not print. `fmt.Println(message)` prints that returned value, and `printNextSteps` prints the second literal line. These inputs are currently in source; they are not silently loaded from environment variables or a network request.

**Microservice connection:** this distinction helps separate business computations from effects such as writing an HTTP response. Our short stdout example ignores printing errors; future I/O code needs explicit error decisions.

## 4. What does `t *testing.T` mean, and how can a failed comparison fail the test command?

`t` is the test function's parameter. `*testing.T` is its pointer type; Go's runner supplies the test state. The `*` appears in a type, not as multiplication or a dereference expression. No `&` address operation is needed in this test. `if got != want` checks behavior; `t.Errorf` reports the difference and marks the test failed. It does not throw a Java-style exception or immediately end the test. After execution the runner produces a nonzero exit status if the test failed.

**Common mistake:** assuming an error-looking string alone fails a test. It is the testing method's action on the runner's state that records failure. See [the testing tutorial](https://go.dev/doc/tutorial/add-a-test).

## 5. Why can swapped string arguments compile yet fail tests, and how do run/build differ?

Both arguments still satisfy the function/formatter's expected types, but their meanings are reversed. The compiler does not know the intended sentence; the test compares against it. The nondefault arguments expose hardcoded implementations, while the executable example checks the user's complete output.

`go run .` compiles and launches the current package. `go build -o bin/first-program .` writes a binary without launching it. Running that binary later executes the already-built code: changing the source does not change the file until rebuilding.

**Pass criterion:** explain these five answers, demonstrate the debugging drill, and leave the supplied behavior passing. Record your confidence and gaps; one successful lesson is not completion of Phase 1.

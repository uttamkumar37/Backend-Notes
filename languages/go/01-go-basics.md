# 1. Go basics

[Start](README.md) → **Go** → [HTTP](02-http-api.md)

Goal: explain a function's inputs, outputs, values, and failures.

Read the overview first. Then use [types](#types-and-conversions), [collections](#arrays-slices-and-copying), [functions](#functions-and-multiple-results), [interfaces](#interfaces-and-method-sets), [errors](#error-design), and [standard library](#standard-library-map) as references. Each section is a separate study session if needed.

## Run first

From the repository root:

```bash
cd languages/go/practice
go run ./basics
go run ./packages
go test -v ./...
```

Read [basics/main.go](practice/basics/main.go) alongside its test. `package main` declares an executable package; `func main()` is its entry point. `import "fmt"` brings in formatted output. `language := "Go"` declares a string; `language = "Rust"` would change it. `welcomeMessage(language, goal)` passes two values and returns a string. `fmt.Println` prints it; `fmt.Sprintf` builds it using `%s` placeholders.

## Modules and packages

A **module** is a dependency/version boundary described by `go.mod`. A **package** groups Go files in one directory. Module path + directory forms its import path.

Trace [packages/main.go](practice/packages/main.go) → [greeting.go](practice/packages/greeting/greeting.go). `Welcome` is capitalized, so another package can call it. Lowercase names are package-private; there is no Java-style `public` keyword.

`go.sum` records downloaded module checksums; it is not a lockfile. These examples use only the standard library, so neither module needs one. When adding a library, choose a version, run `go mod tidy`, and review both module files. `go.work` coordinates local modules; it is optional here. [Module reference](https://go.dev/ref/mod).

## Read the syntax

| Construct | Meaning | Backend use |
| --- | --- | --- |
| `var count int` | Integer, initially `0` | Request counter |
| `const max = 100` | Named constant | Page-size bound |
| `name := "Asha"` | Declare and infer type inside a function | Local input |
| `if err != nil { ... }` | Branch on failure | Stop an invalid write |
| `switch status { ... }` | Choose a case; no implicit fallthrough | Error mapping |
| `for _, item := range items` | Visit values; `_` discards the index | Validate a list |
| `[]string{"IN", "GB"}` | Slice of strings | Supported labels |
| `make(map[string]int)` | Allocate a writable map | Count by label |

Zero values: numbers `0`, booleans `false`, strings `""`, pointers/maps/slices/interfaces `nil`. Appending to a nil slice works; writing to a nil map panics. `len(string)` counts bytes; `range` decodes Unicode code points, which can differ from user-perceived characters.

An inner `:=` can **shadow** an outer name. If an outer value unexpectedly stays unchanged, check whether you declared a new variable in a smaller scope. [Language specification](https://go.dev/ref/spec).

## Structs, pointers, interfaces

Read [project/store.go](project/store.go) when ready:

- `type Investor struct { ... }` groups fields. Tags such as `` `json:"name"` `` name JSON fields.
- Assigning an `Investor` copies its fields. The store returns snapshots, so callers cannot replace its stored name directly.
- `*Store` is a pointer type; `&Store{...}` creates a value and takes its address. In an expression, `*pointer` reads the pointed-to value. Dereferencing nil panics.
- `func (s *Store) Create(...)` has a pointer receiver. Calls share the store and lock. Never copy a used mutex.
- Copying a slice can share its backing array; copying a map shares underlying data. A pointer alone provides no synchronization.

**Java bridge:** structs are values; pointers are also passed by value. Interfaces are satisfied by methods automatically, without `implements`. Define a small interface when a caller needs interchangeable behavior. A typed nil pointer inside an interface can make that interface unequal to `nil`.

Generics help reuse an algorithm across types. Start with an ordinary function until type parameters solve visible repetition.

## Errors and resources

Functions commonly return `(value, error)`. Handle failure before using the value. `fmt.Errorf("operation: %w", err)` wraps a cause; `errors.Is` checks the chain, and `errors.As` extracts a matching error type. Avoid matching error strings. [Error chains](https://go.dev/blog/go1.13-errors).

The project wraps `ErrInvalid`; HTTP code maps that identity to 400. `panic` is for exceptional failures, not ordinary invalid input.

`defer` schedules cleanup at function return, in reverse registration order. Check acquisition succeeded before deferring cleanup. `os.Exit` skips defers; return errors through the owning function first. Do not rely on garbage collection to close files or HTTP response bodies.

Use `encoding/json` for JSON, `strings` for text, and `time.Duration` for durations. Money needs deliberate currency/rounding rules: use integer minor units or an appropriate decimal representation rather than assuming binary floating point is exact.

## Types and conversions

Go checks types at compile time. `int` is convenient for local indexing; use an explicit size such as `int64` when a storage or wire contract requires it. An integer's range and arithmetic behavior matter when accepting external values.

| Type | Useful for | Watch for |
| --- | --- | --- |
| `bool` | Flags and conditions | No implicit truthiness |
| `int`, `int64` | Counts, IDs, minor currency units | Overflow; architecture-dependent `int` size |
| `uint64` | Deliberately nonnegative binary values | Subtracting below zero wraps |
| `float64` | Measurements and approximate calculations | Decimal rounding; NaN/Inf in external data |
| `string` | Immutable byte sequences | Byte length is not character count |
| `byte` | Alias for `uint8` | One UTF-8 character can occupy several bytes |
| `rune` | Alias for `int32`; commonly a code point | A code point is not always a displayed character |

Conversions are explicit: `int64(count)` converts an integer value. `strconv.Atoi("42")` parses text and can fail. `string(65)` produces `"A"`, not `"65"`; `strconv.Itoa(65)` formats the digits. [strconv reference](https://pkg.go.dev/strconv).

Validate range **before** narrowing a value. Choosing `int64` does not prevent an overflowing addition. At business boundaries, define maximum balances, counts, sizes, and timestamps deliberately.

## Declarations, assignment, and scope

A short declaration needs at least one new non-blank variable in its scope. Package-level variables use `var`, not `:=`. Prefer locals and explicit dependencies over mutable global state.

Complete program — save as a temporary `main.go` and run `go run main.go`:

```go
package main

import "fmt"

func main() {
	var attempts int
	const maxAttempts = 3
	name := "Asha"
	name = "Mira"
	if attempts == 0 {
		name := "Local"
		fmt.Println(name, maxAttempts)
	}
	fmt.Println(name, attempts)
}
```

Output: `Local 3`, then `Mira 0`. The inner declaration creates a different `name`. Change it to `name = "Local"` and predict the second line.

## Constants, custom types, and enums

`const` holds compile-time values. `iota` can generate related numeric constants, but a generated number is not automatically a stable public API value.

`type InvestorID int64` creates a distinct defined type; it helps prevent mixing unrelated IDs. `type InvestorID = int64` is an alias, so it does not create that separation. Use custom types where they protect meaning, not for every primitive.

For external statuses, explicit strings such as `"active"` are often easier to inspect. Validate unknown values when parsing JSON or database rows.

## Conditions, loops, and switch

Conditions must be booleans. A `switch` can match values or boolean expressions. A case does not fall through unless explicitly requested; ordinary backend branching rarely needs `fallthrough`.

Use `for` for counted loops, conditions, and indefinite loops. `break` leaves the nearest loop/switch/select; `continue` starts the next loop iteration. A `return` ends the function.

In `for i, value := range items`, `value` is a copy. Assign `items[i]` to replace a slice element. The loop-variable rules depend on the module's language version: Go 1.22+ changed variables declared by the loop. Assignment to pre-existing loop variables still reuses those variables. Do not copy an old closure workaround without understanding the current code. [Loop-variable explanation](https://go.dev/blog/loopvar-preview).

## Arrays, slices, and copying

An array's length is part of its type. A slice describes a window into an array: pointer, length, and capacity. `append` returns a new slice value; always keep that result.

Complete program:

```go
package main

import (
	"fmt"
	"slices"
)

func main() {
	original := []int{10, 20, 30}
	shared := original[:2]
	shared[0] = 99
	independent := slices.Clone(original)
	independent[1] = 77
	grown := append(shared[:len(shared):len(shared)], 40)
	grown[0] = 55
	fmt.Println(original)
	fmt.Println(independent)
	fmt.Println(grown)
}
```

Output: `[99 20 30]`, `[99 77 30]`, `[55 20 40]`. The full slice expression restricts capacity; appending must allocate for the extra element. It does not make the original shared slice immutable.

`slices.Clone` copies elements, not objects reached through nested pointers/maps/slices. `copy(dst, src)` copies up to the shorter slice length; allocating only capacity gives a zero-length destination. A small subslice can keep a large backing allocation alive. [Slices explained](https://go.dev/blog/slices-intro), [slices package](https://pkg.go.dev/slices).

## Maps and sets

Use `value, ok := lookup[key]` to distinguish a missing entry from an entry containing the value type's zero value. `delete` removes a key; deleting an absent key is safe. Map traversal order is unspecified.

For a set, `map[string]struct{}` records membership without storing additional meaning. Map keys must be comparable; slices cannot be keys. Concurrent read/write access requires coordination.

`maps.Clone` makes a shallow map copy. If values contain pointers, slices, or maps, the nested data can remain shared. A value read from a map is a copy: read-modify-write the value if you need to change a struct field. [Maps explained](https://go.dev/blog/maps), [maps package](https://pkg.go.dev/maps).

## Strings, Unicode, and bytes

Complete program:

```go
package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	name := "Go界"
	fmt.Println(len(name), utf8.RuneCountInString(name))
	for offset, character := range name {
		fmt.Printf("%d:%c\n", offset, character)
	}
}
```

Output begins `5 3`; offsets are `0`, `1`, and `2`. A range index is the starting **byte offset**. Avoid slicing arbitrary user text at an unchecked byte offset.

Use `strings.TrimSpace`, `Cut`, `Join`, and `Builder` for ordinary text work. Case folding, normalization, and user-visible length need a clear product rule. Converting a string to `[]byte` produces mutable bytes; mutating those bytes does not rewrite the original string. [Strings and runes](https://go.dev/blog/strings).

## Functions and multiple results

Parameters and return types form the contract. Return a value plus an error when the caller needs both an answer and a failure signal. Keep functions small enough to describe in one sentence.

Variadic `func total(values ...int)` receives a slice-like sequence of arguments; `total(existing...)` expands a slice into arguments. Closures capture variables, so their lifetime and shared mutation matter.

Named return values can help a short API but can also hide which values a bare `return` uses. Prefer explicit returns when cleanup or branching makes the flow hard to follow.

## Pointers, receivers, and mutation

Complete program:

```go
package main

import "fmt"

type Counter struct {
	Value int
}

func (c Counter) CopyIncrement() {
	c.Value++
}

func (c *Counter) Increment() {
	c.Value++
}

func main() {
	counter := Counter{Value: 10}
	counter.CopyIncrement()
	counter.Increment()
	fmt.Println(counter.Value)
}
```

Output: `11`. The value receiver changes a copy. The pointer receiver reaches the original. Go can take the address of an addressable local when calling a pointer method; that convenience does not change interface method-set rules.

Use pointer receivers for mutation, synchronization-containing values, or when copying would be inappropriate. Pointer receivers do not automatically make a type faster. Decide ownership and measure performance. [Methods introduction](https://go.dev/tour/methods/1).

## Struct composition and constructors

Go has no class constructor keyword. A function such as `NewStore` or `newStore` creates a valid value and returns it. Return an error when initialization can fail.

Embedding promotes methods/fields; it does not establish Java-style inheritance or virtual overrides. Prefer a named field when it makes dependency ownership clearer. Expose only the behavior callers need.

The project returns copies of `Investor` so outside code cannot alter stored string fields directly. If you add `[]string` tags, revisit that promise: a shallow struct copy can still share the tag slice.

## Interfaces and method sets

An interface describes behavior. A type satisfies it when it has the required methods. Keep interfaces close to the caller that needs them and use concrete constructors when practical.

If `T` has only a pointer-receiver method, `*T` can satisfy the corresponding interface while `T` cannot. Calling a method through an addressable variable and assigning a value to an interface are different operations.

Use `value, ok := item.(SomeType)` for a checked type assertion. A one-result assertion panics when the dynamic type does not match. A type switch handles several dynamic types. Prefer typed domain contracts over passing `any` everywhere. [Interfaces introduction](https://go.dev/tour/methods/9).

## The typed-nil trap

Complete program:

```go
package main

import "fmt"

type Failure struct{}

func (*Failure) Error() string { return "failed" }

func main() {
	var pointer *Failure
	var err error = pointer
	fmt.Println(pointer == nil)
	fmt.Println(err == nil)
}
```

Output: `true`, then `false`. The interface carries a dynamic pointer type even though its pointer value is nil. Return literal `nil` when there is no error; do not accidentally return a typed nil as `error`.

## Generics and constraints

Complete program:

```go
package main

import "fmt"

type Amount int64

func Sum[T ~int | ~int64](values []T) T {
	var result T
	for _, value := range values {
		result += value
	}
	return result
}

func main() {
	fmt.Println(Sum([]int{1, 2, 3}))
	fmt.Println(Sum([]Amount{10, 20}))
}
```

Output: `6`, then `30`. `T` is a type parameter; the constraint permits types with the listed underlying types. `~int64` includes the defined `Amount` type. This sum has no overflow checks; business money arithmetic needs explicit bounds.

Use generics for reusable algorithms/containers. Use interfaces to express behavior selected through methods. Neither replaces domain modeling. [Generics tutorial](https://go.dev/doc/tutorial/generics).

## Error design

Choose an error contract intentionally:

| Pattern | Useful when | Example |
| --- | --- | --- |
| Sentinel | Caller needs a stable category | Not found |
| Typed error | Caller needs structured detail | Field validation details |
| Wrapped error | Caller needs operation + original cause | Repository read failed |
| Joined errors | Several cleanup actions failed | Write failure plus close failure |

Complete program:

```go
package main

import (
	"errors"
	"fmt"
)

var ErrMissing = errors.New("missing investor")

func lookup(id int) error {
	if id == 0 {
		return fmt.Errorf("lookup id %d: %w", id, ErrMissing)
	}
	return nil
}

func main() {
	err := lookup(0)
	fmt.Println(errors.Is(err, ErrMissing))
	fmt.Println(err)
}
```

Output: `true`, then `lookup id 0: missing investor`. Wrap when preserving the cause is part of the API; avoid exposing storage implementation details accidentally. Handle or return an error at each boundary; logging it at every layer creates duplicate noise.

## Defer, panic, and recover

Deferred call arguments are evaluated when `defer` executes, even though the call runs later. Deferring inside a long loop can retain resources until the surrounding function returns; use a per-item helper when each iteration should close a file promptly.

`recover` handles a panic only from an appropriate deferred function in the panicking goroutine. HTTP recovery middleware cannot catch a panic in an unrelated background goroutine. Recovering also does not prove partially updated state is valid.

Return ordinary validation/network failures as errors. Keep cleanup ownership obvious; consider errors from writes, flushes, closes, and transaction commits where they affect correctness.

## JSON and absent values

Only exported struct fields participate in normal JSON encoding. Tags define the external name; keep request/response DTOs separate when domain data has private or server-owned fields.

For a PATCH field, a pointer can distinguish absent from a supplied non-null value, but plain pointer decoding usually does not distinguish absent from explicit JSON `null`. If all three states matter, use a presence-aware representation or inspect raw fields.

The project uses the compatibility `encoding/json` API. Review its duplicate-key, UTF-8, number, and field-matching behavior. Newer JSON APIs have different defaults; migration must preserve or deliberately change the contract. [JSON reference](https://pkg.go.dev/encoding/json).

## Files, readers, and writers

`io.Reader` produces bytes; `io.Writer` consumes them. These small interfaces allow the same parser to read a file, network body, or test buffer. Reads may return bytes **and** an error; process the bytes returned before deciding what the error means.

`os.ReadFile` is convenient for small bounded files; streaming is better for large inputs. `io.ReadAll` needs a size policy for untrusted sources. `bufio.Scanner` has token-size limits, so configure or replace it when long records are expected.

Use `filepath` for OS file paths and `net/url` for URLs. Concatenating an untrusted path can escape an intended directory; choose a constrained filesystem API/design. [I/O contracts](https://pkg.go.dev/io), [OS APIs](https://pkg.go.dev/os).

## Environment, time, and cancellation

Environment variables enter the process at launch; Go does not load `.env` automatically. Parse and validate them once into a configuration struct, then pass that configuration explicitly. The production guide traces shell, IDE, container, and cluster sources.

Use `time.Duration` for an interval and `time.Time` for an instant. `time.ParseDuration("250ms")` can fail; `time.Duration(250)` means 250 **nanoseconds**, not milliseconds. Use UTC for interchange where the contract calls for it, and named locations for calendar scheduling.

An instant, a local calendar date, and “run every day at 9 AM” are different data models. Daylight-saving changes make a calendar day different from a fixed 24-hour duration. [time reference](https://pkg.go.dev/time).

## Memory, escape analysis, and garbage collection

Go automatically reclaims unreachable memory. A reachable cache, backing array, closure, or goroutine can retain memory for as long as the owner retains it.

Whether a value lives on stack or heap is a compiler decision. Returning a pointer to a local value can be valid; Go preserves the value's lifetime. A pointer does not imply that an allocation must be on the heap.

Start performance work by measuring allocations and retained memory. `GOGC` and `GOMEMLIMIT` influence collector behavior, but neither repairs an unbounded cache. The memory limit is a runtime target, not a complete operating-system RSS cap. [GC guide](https://go.dev/doc/gc-guide).

## Standard library map

| Need | Start with |
| --- | --- |
| Format text / errors | `fmt`, `errors` |
| Parse numbers | `strconv` |
| Text and bytes | `strings`, `bytes`, `unicode/utf8` |
| Sort/search collections | `slices`, `sort`, `maps` |
| JSON and encoding | `encoding/json`, `encoding/base64`, `encoding/hex` |
| HTTP and URLs | `net/http`, `net/url` |
| Files and streams | `os`, `io`, `bufio`, `path/filepath` |
| Time and cancellation | `time`, `context` |
| Coordination | `sync`, `sync/atomic` |
| Structured logs | `log/slog` |
| Random security tokens | `crypto/rand` |
| Tests and local HTTP doubles | `testing`, `net/http/httptest` |

Learn the contract of the package you are using. Similar names in Java and Go do not guarantee identical ownership, blocking, or error behavior.

## Module commands and tooling

Run these inside `practice/` or `project/`:

| Command | Meaning |
| --- | --- |
| `go run .` | Compile/run the current executable package |
| `go build ./...` | Compile package trees in this module |
| `go test ./...` | Run package tests; excludes nested modules |
| `go test -run TestName -v .` | Run matching tests in this package |
| `go fmt ./...` | Format Go source; changes files |
| `go vet ./...` | Inspect suspicious constructs |
| `go doc strings.Cut` | Read local package API documentation |
| `go list ./...` | List packages covered by a module pattern |
| `go mod tidy` | Reconcile requirements/checksums; review its diff |
| `go mod verify` | Verify downloaded dependency-cache integrity |

The `go` directive sets language/minimum-toolchain requirements; it is not a dependency lockfile. `GOTOOLCHAIN=local` prevents automatic toolchain switching for that command. `GOWORK=off` tests a module without a surrounding workspace.

## Package design and initialization

Package names describe a focused responsibility. Avoid generic dumping grounds such as `utils` when names like `validation` or `httpapi` explain the purpose better.

Package initialization happens before `main`; hidden network calls in `init` make startup and tests harder to control. Keep runtime construction in explicit functions. Go rejects import cycles; a cycle often signals unclear ownership between packages.

An `internal` directory restricts who can import its packages based on its enclosing tree. It is useful after the flat project grows enough to need package boundaries. A folder called `pkg` is a convention, not a language-level visibility mechanism. [Module layout](https://go.dev/doc/modules/layout).

## Versions, replace, and workspaces

Module versions follow a compatibility convention: major changes can break callers, while minor/patch changes have narrower intended meanings. Go modules at major version 2+ normally include the major version in their module/import paths. Read release notes and test; version numbers are not proof of correctness. [Module versions](https://go.dev/doc/modules/version-numbers).

A `replace` directive can point development at a local checkout or alternate module version. Review it before sharing/releasing; another developer or CI may not have your filesystem path. A workspace lets local modules cooperate without permanently rewriting each module's dependencies.

Test each module independently with `GOWORK=off` so undeclared dependencies are not accidentally supplied by your workspace. Commit intended module metadata and explain any nonstandard module source.

## Build constraints, embedding, and generation

Build constraints and filename suffixes select source for platforms/features. A file ending `_linux.go` can participate in a Linux build without participating in a macOS build. Test the supported build combinations; “works on my Mac” covers only one.

`go:embed` includes selected files at compile time. It is useful for templates, migration assets, or static resources, but changing an embedded file requires rebuilding the binary. Never embed production secrets. [embed reference](https://pkg.go.dev/embed).

`go generate` runs explicit generation commands; `go build` does not automatically run them. Pin generation tools, review generated diffs, and make the regeneration command reproducible.

## Reflection, unsafe, and CGO

Reflection inspects runtime types/values. Encoders and framework tooling use it; ordinary domain rules are usually clearer with static types and small interfaces. Reflective mutation depends on addressability/settable rules and may panic if used incorrectly. [reflect reference](https://pkg.go.dev/reflect).

`unsafe` bypasses ordinary type/memory protections; CGO crosses into C/native libraries. Both add portability, memory, debugging, and build considerations. Understand these when reading an existing codebase, but introduce them only for a concrete need after considering safer APIs.

## Java to Go translation

| Familiar Java habit | Go approach |
| --- | --- |
| Classes and constructors | Structs, methods, constructor functions |
| Inheritance | Composition and small interfaces |
| Exceptions for expected failures | Explicit returned errors |
| Beans and reflection-based wiring | Explicit constructor calls |
| Overloaded methods | Distinct names or a deliberate options design |
| Annotations for JSON | Struct tags and encoder behavior |
| `null` references | Several kinds of nil; typed-interface nil needs care |
| Thread-local request metadata | Pass context explicitly |
| Package visibility assumptions | Exported identifiers start with uppercase |
| Finalizers for cleanup | Explicit ownership and deferred cleanup |

## Practice

1. Give `Welcome` a second argument. Update the test and explain each changed line.
2. Copy a slice, change one element, and inspect the original. Repeat using an explicit element copy.
3. Rename `Welcome` to `welcome` temporarily. Run tests, explain the export error, then restore it.

Check yourself: why can `:=` hide a value? Why can a copied slice change the original? Why use `errors.Is`?

Continue when you can trace `main → imported function → result`, return an error, and explain who owns mutable data.

Then build a small validator from memory: accept a fictional name, normalize whitespace, reject empty input, and return a wrapped validation error. Add tests for valid text, whitespace-only input, multibyte text, and an oversized name. Explain what copying its input/output does.

## Quick answers

| Question | Answer |
| --- | --- |
| Is Go pass-by-reference? | Arguments are passed by value; a copied pointer can reach shared data. |
| Does `append` always allocate? | No; spare capacity can permit reuse of the backing array. |
| Does a shallow copy isolate nested maps? | No; define the copy depth needed by your contract. |
| Can I compare two slices with `==`? | No, except comparison to nil; choose an element-comparison function. |
| Can map order define API order? | No; sort explicitly. |
| Does an interface containing a nil pointer equal nil? | Not necessarily; its dynamic type is still present. |
| Does a pointer receiver provide thread safety? | No; synchronization and ownership do. |
| Should every helper become generic? | Only when the abstraction makes repeated behavior clearer. |
| Does GC replace `Close`? | No; close resources according to ownership. |
| Does a passing test mean mastery? | You also need independent explanation and debugging. |

# 4. Concurrency and cancellation

[Data](03-databases-cache.md) → **Concurrency** → [Production](05-production.md)

Goal: explain who owns work, who stops it, and who waits for it to finish.

After the overview, study [channels](#channel-behavior-table), [mutexes](#a-shared-counter-with-a-mutex), [cancellation](#a-cancellable-wait), [worker pools](#a-bounded-worker-pool), and [failure diagnosis](#races-deadlocks-and-leaks).

## Choose the simplest coordination

| Tool | Use it for | Common mistake |
| --- | --- | --- |
| Goroutine: `go fn()` | Concurrent work | Starting work with no owner or bound |
| `sync.Mutex` | Protect a shared invariant | Protecting a write while leaving reads unprotected |
| `sync.RWMutex` | Shared reads/exclusive writes | Assuming it always outperforms a mutex |
| Channel | Transfer values / coordinate | Blocking on a send after the receiver exits |
| `sync.WaitGroup` | Wait for owned tasks | Adding work too late or forgetting `Done` |
| `select` | Wait on communication/cancellation | Assuming it gives priority to cancellation |
| `sync.Once` | One initialization attempt | Expecting it to retry after failure or rotate a token |
| Atomics | Small independent atomic operations | Treating several atomic fields as one transaction |

In [store.go](project/store.go), the same lock protects both the ID counter and map. Keeping them together prevents duplicate IDs and unsafe map access. Returned investor structs contain independent value fields. Future slices/maps inside them would need an ownership review.

Synchronization makes writes visible to other goroutines; merely starting a goroutine or waiting an arbitrary time does not. [Go memory model](https://go.dev/ref/mem).

## Context is a cooperative stop signal

At an HTTP boundary, start from `r.Context()`. For a database or downstream call, derive a timeout, defer its cancellation, and pass that context into the operation.

`context.WithCancel` returns a child context and a cancel function. Cancellation closes its `Done()` channel; `ctx.Err()` describes cancellation/deadline expiry. Work must observe that signal or use APIs that do. Cancellation cannot kill arbitrary code.

An HTTP handler normally owns its request work. A job that must outlive the request needs a separate application lifecycle or durable queue, with explicit delivery/retry behavior. Avoid detached goroutines that silently lose work during shutdown. [Context guide](https://go.dev/blog/context).

**Java bridge:** goroutines are lightweight concurrent tasks, not Java threads you should start without bounds. Context makes cancellation explicit through the call chain; it is not thread-local state.

## Bound work and deadlines

A worker pool has a fixed worker count, a bounded queue, a producer that can stop, and a coordinator that waits for all workers. The sender/owner closes the channel once sending is finished. Receivers do not close a shared input channel.

Every blocking send/receive needs a plan for cancellation. Define overload behavior when the queue is full. More goroutines cannot create more database connections or remote capacity.

Give each operation part of the remaining request budget. If the caller has 500 ms left, a fresh 5-second downstream timeout cannot extend that caller's lifetime. Include pool waits and retries in the budget.

## Shutdown is part of correctness

Read `serve` in [main.go](project/main.go):

1. `signal.NotifyContext` observes interrupt/SIGTERM.
2. The server stops accepting connections and gets time to finish active requests.
3. Cleanup uses a fresh timeout from `context.Background()`, because the signal context is already cancelled.
4. The process waits for shutdown; if its deadline expires, it closes remaining connections and reports failure.

As dependencies are added, drain work before closing pools/clients. Explicitly stop and join background workers. Long-lived upgraded connections need their own shutdown handling. [HTTP shutdown](https://pkg.go.dev/net/http#Server.Shutdown).

## Concurrency, parallelism, and scheduling

Concurrency structures several activities so they can make progress independently. Parallelism executes work simultaneously. A single process can have concurrent goroutines even when only one is executing Go code at a time.

Goroutines are scheduled onto operating-system threads by the runtime. `GOMAXPROCS` affects simultaneous Go execution; it is not a goroutine-count limit or a database capacity setting. Check runtime/container defaults before overriding them.

I/O-bound and CPU-bound workloads need different capacity decisions. Starting more goroutines can hide waiting, but it can also increase memory, contention, and pressure on downstream services.

## Who owns a goroutine?

Before writing `go`, answer four questions: who starts it, how does it stop, who waits for completion, and where does failure go?

Request-owned work should normally end before the handler returns. Application-owned workers live under application shutdown. Durable business work belongs in a persistence/delivery design rather than relying on process memory.

A goroutine closure can retain request bodies, buffers, pointers, or authentication data. Avoid retaining a whole framework context when only a small immutable value is needed.

## Channel behavior table

| Operation | Behavior |
| --- | --- |
| Send to unbuffered channel | Waits for a receiver |
| Send to buffered channel | Waits if the buffer is full |
| Receive from empty open channel | Waits for a sender |
| Receive from closed channel | Drains buffered values, then returns zero value with `ok == false` |
| Send to closed channel | Panics |
| Close an already closed channel | Panics |
| Send/receive on nil channel | Blocks indefinitely unless another selected case proceeds |
| Close a nil channel | Panics |

A buffer absorbs a finite mismatch in timing; it does not provide unlimited capacity or prove the consumer finished processing. Use a separate acknowledgment when completion matters.

## A channel with one closing owner

Complete program — save as a temporary `main.go` and run `go run main.go`:

```go
package main

import "fmt"

func main() {
	values := make(chan int)
	go func() {
		defer close(values)
		for value := 1; value <= 3; value++ {
			values <- value
		}
	}()
	total := 0
	for value := range values {
		total += value
	}
	fmt.Println(total)
}
```

Output: `6`. The producer owns closure. The range finishes after closure and draining. A receiver returning early would leave this producer blocked; real interruptible pipelines need cancellation around blocking operations.

## select, default, and fairness

`select` waits until a communication case can proceed. When several are ready, do not depend on source order to prioritize one. A default case makes the select nonblocking.

A `for` loop containing a default case can spin and consume CPU when no useful work is ready. A nil channel can disable a case deliberately, but an accidental nil channel can make a pipeline hang.

Checking cancellation once before a send is not sufficient: the send may block after the check. Include cancellation in the same select as the blocking communication.

## A shared counter with a mutex

Complete program:

```go
package main

import (
	"fmt"
	"sync"
)

func main() {
	var mutex sync.Mutex
	var workers sync.WaitGroup
	count := 0
	for worker := 0; worker < 4; worker++ {
		workers.Add(1)
		go func() {
			defer workers.Done()
			for item := 0; item < 1000; item++ {
				mutex.Lock()
				count++
				mutex.Unlock()
			}
		}()
	}
	workers.Wait()
	fmt.Println(count)
}
```

Output: `4000`. Register each task before starting it. WaitGroup waits for completion; the mutex protects concurrent updates. Reading after all these workers finish avoids a simultaneous read/write.

Run this temporary program with `go run -race main.go`. Removing the mutex creates a race; `count++` is not one indivisible synchronized operation.

## Lock the invariant, not just a field

The store's next ID and map insertion form one operation. Lock both together. Protecting each field independently can still break the relationship between them.

Keep critical sections short. Avoid network I/O, slow callbacks, or user-controlled code while holding a lock. Use a consistent ordering when acquiring multiple locks. Do not copy a used lock-containing value.

An RWMutex cannot be casually upgraded from a read lock to a write lock. Choose a safe retry/recheck design instead. Measure before replacing a plain mutex with a more complex synchronization strategy. [sync reference](https://pkg.go.dev/sync).

## Context trees and causes

A child context inherits parent cancellation/deadline and may impose a shorter deadline. Cancelling a child does not cancel its parent. Always release a derived context's resources by calling its cancel function when its owner is done.

Pass context as the first parameter to I/O-oriented functions. Do not pass nil. Context values are for request-scoped metadata across API boundaries, not a hidden optional-argument/configuration bag.

`context.WithCancelCause` can preserve a more specific reason while `ctx.Err()` retains its cancellation category. Detaching cancellation with an API such as `WithoutCancel` creates a new ownership problem; it does not define how work will eventually stop. [context reference](https://pkg.go.dev/context).

## A cancellable wait

Complete program:

```go
package main

import (
	"context"
	"fmt"
	"time"
)

func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	fmt.Println(wait(ctx, time.Hour))
}
```

Output: `context canceled`; it does not wait an hour. Unlike `time.Sleep`, this operation listens for cancellation while waiting. Timer APIs have changed across Go releases; check the documentation for your module/toolchain rather than copying older drain/reset recipes. [Timer reference](https://pkg.go.dev/time#Timer).

## Timeouts, tickers, and retries

Use one overall deadline and budget each attempt within it. Retries consume both elapsed time and downstream capacity. A retry delay should also be cancellable.

A ticker supplies periodic opportunities; it is not a durable scheduler and may drop ticks when a receiver is slow. Stop periodic work during shutdown. Several replicas starting together may need jitter to avoid synchronized load.

Avoid creating a new timer repeatedly in a hot loop unless that allocation/lifetime is understood. Profile first; correctness comes before timer micro-optimizations.

## A bounded worker pool

Complete program. Four inputs, two workers, cancellable producer/sends, one results closer:

```go
package main

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

func squareAll(ctx context.Context, input []int, workerCount int) ([]int, error) {
	if workerCount < 1 {
		return nil, fmt.Errorf("worker count must be positive")
	}
	jobs := make(chan int)
	results := make(chan int)
	var tasks sync.WaitGroup
	tasks.Add(workerCount + 1)
	go func() {
		defer tasks.Done()
		defer close(jobs)
		for _, value := range input {
			select {
			case <-ctx.Done():
				return
			case jobs <- value:
			}
		}
	}()
	for worker := 0; worker < workerCount; worker++ {
		go func() {
			defer tasks.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case value, ok := <-jobs:
					if !ok {
						return
					}
					select {
					case <-ctx.Done():
						return
					case results <- value * value:
					}
				}
			}
		}()
	}
	go func() {
		tasks.Wait()
		close(results)
	}()
	output := make([]int, 0, len(input))
	for value := range results {
		output = append(output, value)
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	sort.Ints(output)
	return output, nil
}

func main() {
	output, err := squareAll(context.Background(), []int{1, 2, 3, 4}, 2)
	fmt.Println(output, err)
}
```

Output: `[1 4 9 16] <nil>`. Results are sorted because completion order is not promised. The WaitGroup includes the producer, so results close only after every producer/worker exits.

This bounds active work, not total result memory: the collected output still grows with input. The arithmetic is a small-data demonstration without overflow handling. For large streams, bound input/output storage and define ordering/error semantics too.

## errgroup and parallel calls

`errgroup.WithContext` coordinates goroutines that return errors. The derived context is cancelled after the first non-nil task error or when `Wait` completes. Tasks still need to observe cancellation and return; the group does not kill them.

Limit concurrency deliberately. A blocked `Group.Go` call waiting for a concurrency slot is itself a point to reason about; do not assume every admission path is automatically cancellable. Do not reuse the group's cancelled context for new work after `Wait`.

Use this library when error-returning parallel tasks justify it; pin the external module version when adding it. The current repository does not depend on it. [errgroup reference](https://pkg.go.dev/golang.org/x/sync/errgroup).

## Once, atomics, and immutable snapshots

`sync.Once` ensures one execution, including an initialization attempt that fails. Store the error explicitly if callers need it; use a different design if failures should permit retries. Credential renewal is repeated lifecycle work, so it does not belong inside a one-time initializer.

An atomic operation protects a specific supported operation/value. Several atomic updates do not automatically form one consistent business transaction. An immutable configuration snapshot can be published atomically only if readers cannot mutate data reachable through it.

`sync.Map` has specialized use cases. A typed map with an ordinary mutex is often easier to reason about when several entries/fields participate in one invariant.

## Backpressure and load shedding

For every queue choose a size and full-queue behavior: wait within a budget, reject, drop according to policy, or persist for later. State which caller learns that work was rejected or lost.

Bound fan-out. Ten incoming requests each starting 100 downstream calls can create 1000 operations before retries. Bulkheads divide capacity so one slow dependency cannot occupy every worker.

## Races, deadlocks, and leaks

| Failure | Meaning | Evidence to collect |
| --- | --- | --- |
| Data race | Unsynchronized conflicting memory access | Race report with access stacks |
| Logical race | Individually synchronized steps violate a larger rule | Interleaving and invariant test |
| Deadlock | Work cannot progress because of circular/blocking dependencies | Goroutine stacks and channel/lock ownership |
| Goroutine leak | Work remains alive beyond its intended lifetime | Repeated workload plus persistent waiting stacks |
| Starvation/contention | Work progresses too slowly under competition | Latency, lock/block profiles, queue waits |

Arbitrary sleeps make concurrency tests timing-sensitive. Use channels/barriers to signal “started,” “blocked,” and “allowed to finish,” plus bounded test deadlines. Runtime goroutine counts include background activity; one exact count is not a complete leak test.

## Shutdown with background workers

Choose whether accepted work drains or is cancelled. Stop admission first, then signal/drain owned work, wait within a usable cleanup deadline, and close dependencies after their users stop.

The signal context and a cleanup context have different purposes. Passing an already-cancelled signal context into `Shutdown` can skip the grace period. An application also needs a policy for workers that ignore cancellation.

In Kubernetes, make the total termination budget cover traffic removal, request/job draining, and cleanup. An interrupt reaching a shell wrapper is not proof it reached the actual application process.

## Practice

From the repository root:

```bash
bash languages/go/check.sh race
```

1. In a temporary branch, remove the store's read locking. Explain the race report and restore the protection.
2. Implement a bounded worker pool; cancel it while its producer is blocked. Verify every owned goroutine exits.
3. Run the built API, send an interrupt, and explain why cleanup needs a context separate from the signal.

The race detector checks executed paths; a passing run does not prove every possible path is safe. Deadlocks and goroutine leaks need separate reasoning/tests. [Race detector](https://go.dev/doc/articles/race_detector).

Check yourself: who closes a channel? Why doesn't cancellation kill a goroutine? Why can `sync.Once` be wrong for refreshing credentials?

Additional exercises: cancel the worker-pool example before it starts; add a worker error and define whether partial results are returned; protect a two-field invariant; reproduce a deadlock in a disposable program and identify both waiting stacks. Run appropriate cases under the race detector, and describe what it cannot prove.

## Quick answers

| Question | Answer |
| --- | --- |
| Does starting a goroutine guarantee it runs before main exits? | No; the owner must wait where completion is required. |
| Does a buffered send prove processing finished? | No; it can mean only that the value entered the buffer. |
| Who should close a results channel with several workers? | A coordinator after all senders have finished. |
| Does a context cancel a database write retroactively? | No; cancellation and the actual commit outcome must be reasoned about separately. |
| Can `select` always prioritize cancellation? | Do not depend on source-order priority among ready cases. |
| Does `sync.Once` retry a failed initializer? | No; choose another design if retry is required. |
| Does a clean race run prove freedom from deadlock? | No; they are different failure modes. |
| Is an atomic counter enough for a transfer invariant? | No; coordinated business state needs a suitable transaction/synchronization design. |

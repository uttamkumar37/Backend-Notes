# Lesson and engineering standards

The learner already knows Java, Spring Boot, REST, and SQL, and is a beginner in Go. Use that backend experience as a bridge while explaining all new Go syntax. Work without office credentials, repositories, VPN, or paid cloud services.

## One lesson at a time

Each lesson fits approximately 1–2 hours and contains, in order:

1. Today's objective and a concrete success condition.
2. Prerequisite concepts, with links to earlier local lessons.
3. Theory: **what → why → when → how**.
4. Java/Spring comparison, including where the analogy breaks.
5. A small runnable example with minimal dependencies.
6. Line-by-line explanation of all new syntax and behavior.
7. macOS commands with exact working directory and prerequisites.
8. Expected output; label variable values such as timestamps/ports.
9. Practical exercise with observable verification.
10. Debugging challenge: symptom, reproduction, diagnostics, cause, fix, recovery.
11. Exactly five knowledge-check questions and explained answers.
12. Short recap connecting syntax to real microservice behavior.
13. Completion criteria and honest progress evidence.
14. Next topic, followed by a learner confirmation checkpoint.

Do not advance automatically after publishing a lesson. The learner specifically requested confirmation before progression, and environment confirmation before Phase 1. Documentation of future plans is not teaching/completing those phases.

## Required files when a Go topic is implemented

Provide `README.md`, `EXERCISES.md`, `DEBUGGING.md`, and `INTERVIEW.md` for each manageable topic. Add appropriate runnable `.go` sources and meaningful `_test.go` files, with module boundaries explained. Infrastructure/configuration topics can have configuration files and automated smoke checks instead of an artificial `main.go`.

Session 0 is a shell inventory, so it has an inspection script rather than a Go module. Future topics are added when reached; avoid dozens of empty folders or repeated boilerplate pretending to be completed lessons.

## Code explanation checklist

For every Go example, identify:

- Each import, why it is needed, and the package's role.
- Every new keyword, operator, declaration, scope, value origin, and variable type.
- Function parameters, return types, multiple results, call sites, and error propagation.
- Why `:=` declares and `=` assigns at that point; highlight shadowing.
- What `&` takes the address of, what `*` denotes in a type or expression, and nil behavior.
- Struct values, copies, method sets, pointer/value receivers, interfaces, and ownership where present.
- What happens internally, which operation blocks, and when resources are released.
- How to run/debug it, expected terminal output, and common beginner mistakes.
- Java comparison and realistic microservice use, while labeling educational shortcuts.

Explain syntax present in the example; do not force pointers or concurrency into a lesson that does not use them. Once a concept is established, link back and still explain the example's value flow and failure paths.

## Configuration provenance example requirement

When teaching `os.Getenv("LOG_LEVEL")`, trace the whole journey: a nonsecret value is exported in Terminal or set in an IDE launch configuration; the child process inherits it; Podman can inject it into a container; Kubernetes can supply it through a ConfigMap reference; the application reads its own environment at runtime. Explain defaults, missing versus empty values with `LookupEnv`, validation, safe inspection of this named field, and when a process restart is needed. A `.env` file is not automatically read by Go or inherited by a container without explicit loading/injection.

When teaching cancellation, explain who creates/cancels the context, how goroutines observe `Done()`, why cancellation does not kill arbitrary work, and how pod termination reaches `http.Server.Shutdown` and dependency cleanup. A shutdown cleanup deadline must remain usable after the signal context has been cancelled.

## Engineering and verification

Use idiomatic names, focused packages, explicit errors, constructor injection where useful, context-aware I/O, safe concurrency, bounded timeouts, resource cleanup, and graceful shutdown. Keep production secrets out of source, logs, fixtures, and generated documents; use fictional investors and development-only dummy credentials.

At each phase end, the current platform must compile, run, and pass relevant tests. Test behavior rather than restating implementation. Integration tests use real isolated dependencies where needed; they are explicitly opt-in until dependencies are available. Every Go module has automated tests; a command ending in `[no test files]` is not test coverage. Do not introduce broken examples into the default passing module suite.

For each test record the command, module directory, dependencies, expected assertions, and actual result. Clearly separate assistant verification from learner understanding. Use formatting/vet/fast tests for ordinary changes; add race, integration, end-to-end, benchmark, or profile checks when relevant. Never claim a cluster deploy, Sentinel failover, Vault renewal, or passing integration suite that was not actually executed.

Keep all learning content inside `languages/go/`. Record stable versions when installing and pin versions/checksums in implemented builds. Recheck official documentation before version-sensitive setup. Future cloud publishing, authenticated account actions, and repository-wide workflow activation are outside the initial local-note setup.

## Phase assessment

Ask the learner to implement a representative feature independently, explain its startup/value/error flow, run appropriate tests, diagnose one introduced failure, and review the design as a PR. Record gaps and repeat practice as needed. Only mark the phase complete after independent explanation and implementation are demonstrated.

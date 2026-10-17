# Learning progress

**Current lesson:** [Phase 1, Lesson 2 — modules, packages, imports, and workspaces](01-go-fundamentals/02-modules-packages/README.md).

**Current state:** you reported Lesson 1 done and requested Lesson 2. Lesson 1 is recorded as complete based on your report; confidence and detailed exercise evidence were not supplied. Lesson 2 is available and its supplied independent/workspace tests pass; your Lesson 2 practice remains pending. No phase assessment is complete. Assistant verification and learner-reported completion are recorded separately.

Use `—` for no learner evidence, `In progress` for partial work, and `Done` only with evidence. Confidence is self-rated: 1 = need a worked example, 2 = can follow, 3 = can implement with references, 4 = can debug independently, 5 = can teach and review. Assess phase completion separately from passing tests.

## Phase checklist

| Topic | Theory completed | Code completed | Tests passed | Debugging completed | Confidence (1–5) | Remaining gaps / evidence |
| --- | --- | --- | --- | --- | --- | --- |
| 0. Environment | Learner authorized progression; knowledge not assessed | Inventory script available | Script syntax + live inventory verified by assistant | — | — | Prior exercise evidence unreported; editor SDK/extension and patch update remain unverified |
| 1. Fundamentals | Lesson 1 done (learner-reported); Lesson 2 notes available | Lessons 1–2 modules supplied; later lessons pending | Both modules pass independently and in the supplied workspace (assistant) | Lesson 1 done (learner-reported); Lesson 2 drills supplied | — | Lesson 2 practice, Lessons 3–18, and independent assessment pending |
| 2. Startup/config/shutdown | — | — | — | — | — | Pending |
| 3. Gin/API | — | — | — | — | — | Pending |
| 4. Zerolog/observability | — | — | — | — | — | Pending |
| 5. PostgreSQL/pools | — | — | — | — | — | Pending |
| 6. Cache/Sentinel | — | — | — | — | — | Pending |
| 7. Concurrency/context | — | — | — | — | — | Pending |
| 8. Vault | — | — | — | — | — | Pending |
| 9. Microservices | — | — | — | — | — | Pending |
| 10. Podman/hybrid development | — | — | — | — | — | Pending |
| 11. Kubernetes/OKD | — | — | — | — | — | Pending |
| 12. Testing/debugging | — | — | — | — | — | Pending |
| 13. CI/CD/production engineering | — | — | — | — | — | Pending |

## Session 0 detail

| Topic | Theory completed | Code/commands completed | Verification passed | Debugging completed | Confidence | Gaps |
| --- | --- | --- | --- | --- | --- | --- |
| Host OS and architecture | — | Assistant captured macOS/ARM64 | Observed machine output recorded | — | — | Explain Darwin versus Linux target |
| `PATH` and executable resolution | — | Inventory script available | Assistant found core commands | — | — | Learner temporary-child-shell drill pending |
| SDK and toolchain selection | — | Assistant read installed version | Native Go 1.27.1 observed | — | — | Review available patch update; explain per-command override |
| Editor and debugger | — | VS Code launcher observed | SDK/extension/breakpoint not checked | — | — | Learner chooses/configures editor |
| Podman client and VM | — | Assistant inspected versions/list | Existing running VM observed; engine workload untested | — | — | Explain execution locations; dependencies can wait |
| Deferred tools | — | Assistant inventoried availability | Missing entries documented | — | — | Identify phase for each tool |

## Phase 1 lesson 1 detail

Module directory: `languages/go/01-go-fundamentals/01-first-program/`. Installed verification toolchain: Go 1.27.1 on `darwin/arm64`; no SDK upgrade or dependency installation was performed.

| Topic | Theory completed | Code completed | Tests passed | Debugging completed | Confidence | Gaps |
| --- | --- | --- | --- | --- | --- | --- |
| Module, package, entry point | Done (learner-reported) | Supplied code; lesson complete by report | `go vet ./...` clean; run/build output verified by assistant | Lesson completion reported; no separate diagnostic supplied | — | No gaps reported; confidence unspecified |
| Variables, arguments, result, printing | Done (learner-reported) | Supplied code; lesson complete by report | `TestWelcomeMessageUsesArguments` and `Example` pass | Swapped-argument drill verified by assistant; lesson done by report | — | No gaps reported; confidence unspecified |
| First automated test | Done (learner-reported) | Supplied test; lesson complete by report | `GOTOOLCHAIN=local go test -v ./...` passes both checks | Recovery verified by assistant; lesson done by report | — | Detailed exercise/test evidence not separately supplied |

## Phase 1 Lesson 2 detail

Module directory: `languages/go/01-go-fundamentals/02-modules-packages/`. Workspace file: `workspace/go.work`, selecting this module and Lesson 1. Go 1.27.1 remains the installed SDK. The only new external dependency is test helper `github.com/google/go-cmp v0.7.0`; Go generated its actual checksums.

| Topic | Theory completed | Code completed | Tests passed | Debugging completed | Confidence | Gaps |
| --- | --- | --- | --- | --- | --- | --- |
| Packages, import paths, exports | Notes available; learner pending | Command and `greeting` package supplied | Command output example and external-package input test pass (assistant) | Export-visibility failure and repair verified in temporary copy | — | Trace import and complete learner drill |
| Module requirements and checksums | Notes available; learner pending | `go.mod` and generated `go.sum` supplied | Independent tests, cache verification, and idempotent tidy checked (assistant) | Read-only test detects missing checksum; tidy repairs it in temporary copy | — | Explain content/metadata hashes and dependency commands |
| Workspace and independent modules | Notes available; learner pending | Relative-path `workspace/go.work` supplied | Workspace tests cover both modules; independent Lesson 2 tests pass (assistant) | Missing member failure and recovery verified in temporary copy | — | Run member patterns, repair a fault, and explain `GOWORK=off` |

No phase is marked complete. Complete Lesson 2 practice and confirm before moving to Lesson 3.

## Learning log

For each completed lesson add a dated row, including the module directory, exact verification command/result, debugging cause/fix, and what you can explain without copying.

| Date | Lesson/topic | Theory | Code | Tests / observed result | Debugging evidence | Confidence | Remaining gap / next session |
| --- | --- | --- | --- | --- | --- | --- | --- |
| 2026-10-17 | Initial machine assessment (assistant) | Learner pending | Inventory script created | `bash -n` and live inventory succeeded | Learner drill pending | Not assessed | Confirm environment before Phase 1 |
| 2026-10-17 | Learner requested next phase; Phase 1 lesson 1 supplied | Explanations available; learner pending | Separate standard-library module and tests | Unit + executable-example checks pass; vet/run/build verified | Temporary swapped-argument fault fails tests; restoration passes | Not assessed | Complete lesson 1 exercises and five checks before lesson 2 |
| 2026-10-17 | Learner reported Lesson 1 done; requested Lesson 2 | Lesson 1 complete by report; Lesson 2 notes supplied | Lesson 2 command/library, pinned test helper, workspace | Independent/workspace tests, native build/output, vet/cache verification pass | All three failure/recovery drills verified in temporary copies | Not specified | Complete Lesson 2 exercises and five checks before Lesson 3 |

## Phase assessment record

At the end of a phase record: task attempted without a solution, explanation of design/error paths, test evidence, injected-fault diagnosis, unresolved gaps, and learner confirmation. Repeat the relevant exercise if the main concepts cannot yet be explained and implemented independently. The estimate in the roadmap is a planning aid, not a deadline for checking boxes.

No assessments have been attempted yet.

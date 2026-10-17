# Phase 1 — Go fundamentals

**Current lesson:** [02 — Modules, packages, imports, and workspaces](02-modules-packages/README.md).

You reported Lesson 1 complete and requested Lesson 2. Today, trace a local import into its package, inspect real dependency metadata, and use a small workspace to test both lesson modules. The original [practice program](../practice/01-basics/hello_backend.go) and Lesson 1 source are preserved.

| Lesson | Material | Learner status |
| --- | --- | --- |
| 01. First module and program | [Notes](01-first-program/README.md), [exercises](01-first-program/EXERCISES.md), [debugging](01-first-program/DEBUGGING.md), [five questions](01-first-program/INTERVIEW.md), runnable code and tests | Complete (learner-reported); confidence unspecified |
| 02. Modules and packages | [Notes](02-modules-packages/README.md), [exercises](02-modules-packages/EXERCISES.md), [debugging](02-modules-packages/DEBUGGING.md), [five questions](02-modules-packages/INTERVIEW.md), code, tests, and workspace | Started; learner practice pending |
| 03–18. Remaining fundamentals | [Complete lesson sequence](../notes/roadmap.md#phase-1--go-fundamentals) | Planned; not yet delivered |

Run the current lesson's independent commands from `02-modules-packages/`, which contains its `go.mod`. A single pinned test helper demonstrates `go.sum`; the program itself uses only standard/local packages. Its `workspace/` example groups the two existing modules. No PostgreSQL, Valkey, Vault, containers, or cluster is needed.

After the lesson, record your own evidence in [PROGRESS.md](../PROGRESS.md). Phase 1 remains in progress until its later lessons and independent practical assessment are complete. Confirm Lesson 2 before continuing to Lesson 3.

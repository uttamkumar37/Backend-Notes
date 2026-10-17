# Go Backend & Cloud-Native Engineering Lab

A practical path from Java/Spring Boot experience to Go backend engineering, using a personal Apple Silicon Mac and fictional investor data.

**Learning root:** `Backend-Notes/languages/go/`. Keep all Go notes, examples, projects, scripts, and infrastructure here. The master prompt's name, `go-cloud-native-engineering-lab`, is the curriculum name; this existing folder is its home. A nested Git repository is unnecessary.

## Start here

1. Continue with [Phase 1, lesson 2: modules, packages, imports, and workspaces](01-go-fundamentals/02-modules-packages/README.md).
2. Refer back to [Session 0](notes/00-environment/README.md) and the [observed environment](notes/environment-assessment.md) when needed.
3. Follow the [13-phase roadmap](notes/roadmap.md), one lesson at a time.
4. Record your own understanding in [PROGRESS.md](PROGRESS.md).

Run the safe inventory from any Terminal directory:

```bash
bash /Users/uttamkumar/Uttam-workspace/GitHubProject/Backend-Notes/languages/go/scripts/check-environment.sh
```

It reports versions and tool availability. It does not install software, change configuration, start a VM, or read application credentials. Its output and each underlying command are explained in Session 0.

**Current checkpoint:** Lesson 1 is complete by your report. Lesson 2 is available with runnable code, dependency files, a workspace, automated tests, exercises, and debugging practice; its learner completion remains pending. Future lessons and Investor Platform implementations will be added as their lessons are reached. A roadmap entry describes planned work, not working software.

## Reference notes

| File | Purpose |
| --- | --- |
| [SETUP.md](SETUP.md) | Install and verify tools when their phase needs them; distinguish macOS, VM, and container processes |
| [COMMANDS.md](COMMANDS.md) | Commands with working-directory and prerequisite information |
| [TROUBLESHOOTING.md](TROUBLESHOOTING.md) | Current setup fixes and a catalog of future office-style failure drills |
| [PROGRESS.md](PROGRESS.md) | Theory, code, test evidence, debugging, confidence, and gaps |
| [Lesson format](docs/lesson-format.md) | Teaching rules, code explanation standards, and phase assessments |
| [Original master prompt](notes/MASTER-PROMPT.md) | Preserved requirements; this folder location overrides the suggested standalone repository |

## Repository layout

```text
languages/go/
├── README.md
├── SETUP.md
├── COMMANDS.md
├── TROUBLESHOOTING.md
├── PROGRESS.md
├── 01-go-fundamentals/
│   ├── README.md
│   ├── 01-first-program/       # Lesson 1: first module, program, and tests
│   └── 02-modules-packages/    # Lesson 2: packages, dependency files, and workspace
│       ├── go.mod
│       ├── go.sum
│       ├── main.go
│       ├── main_test.go
│       ├── greeting/           # Importable package and its test
│       ├── workspace/go.work
│       ├── README.md
│       ├── EXERCISES.md
│       ├── DEBUGGING.md
│       └── INTERVIEW.md
├── notes/
│   ├── MASTER-PROMPT.md
│   ├── roadmap.md
│   ├── environment-assessment.md
│   └── 00-environment/
│       ├── README.md
│       ├── EXERCISES.md
│       ├── DEBUGGING.md
│       └── INTERVIEW.md
├── practice/
│   └── 01-basics/hello_backend.go
├── scripts/check-environment.sh
└── docs/lesson-format.md
```

The existing `hello_backend.go` is preserved. Lesson 1 explains an adaptation in its own tested module. Its presence and passing tests do not mean the learner has completed the fundamentals.

The eventual platform will live at `projects/investor-platform/`, with infrastructure under `infrastructure/`. See the roadmap for planned phase folders and module boundaries. No office VPN, repositories, credentials, or paid cloud deployment is required.


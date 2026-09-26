# Backend Engineering Learning

A clean workspace for learning backend engineering one language at a time.

## Structure

```text
.
├── backend-concepts/
│   └── README.md
└── languages/
    ├── go/
    │   ├── README.md
    │   ├── notes/
    │   │   └── roadmap.md
    │   └── practice/
    │       └── 01-basics/
    │           └── hello_backend.go
    └── java/
        ├── README.md
        ├── notes/
        │   └── roadmap.md
        └── practice/
            └── 01-basics/
                └── HelloBackend.java
```

## Idea

- `backend-concepts/` keeps language-independent backend concepts.
- `languages/java/` keeps Java-specific notes and practice.
- `languages/go/` keeps Go-specific notes and practice.
- No Gradle, Maven, framework, or deep folder structure at the start.
- Use simple files first. Add tools only when they solve a real problem.

## Start

```bash
cd languages/java/practice/01-basics
javac HelloBackend.java
java HelloBackend
```

```bash
cd languages/go/practice/01-basics
go run hello_backend.go
```

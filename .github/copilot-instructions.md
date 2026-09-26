# Repository Copilot Instructions

## Repository Overview

This repository is a simple backend engineering learning workspace. It is intentionally small and beginner-friendly at the structure level.

The repo separates:

- `backend-concepts/` for backend concepts that apply to any language
- `languages/java/` for Java-specific learning and practice
- `languages/go/` for Go-specific learning and practice

Do not turn this into a full framework project unless explicitly requested.

## Current Structure

```text
backend-concepts/
  README.md
languages/
  java/
    README.md
    notes/
      roadmap.md
    practice/
      01-basics/
        HelloBackend.java
  go/
    README.md
    notes/
      roadmap.md
    practice/
      01-basics/
        hello_backend.go
```

## Development Style

- Keep examples small and easy to run.
- Prefer plain language and practical examples.
- Avoid unnecessary folders.
- Avoid build tools until they are actually needed.
- Keep language-specific content inside that language folder.
- Keep shared backend concepts inside `backend-concepts/`.
- Do not add large generated content.
- Do not add placeholder-heavy files.

## Java Rules

- Use plain Java files for now.
- Run Java examples with `javac` and `java`.
- Do not add Maven, Gradle, Spring Boot, or package-heavy folder structure unless explicitly requested.
- Keep Java practice files easy to read and modify.

Example:

```bash
cd languages/java/practice/01-basics
javac HelloBackend.java
java HelloBackend
```

## Go Rules

- Use plain Go files for now.
- Run Go examples with `go run`.
- Do not add `go.mod`, frameworks, or web server structure unless explicitly requested.
- Keep Go practice files easy to read and modify.

Example:

```bash
cd languages/go/practice/01-basics
go run hello_backend.go
```

## Learning Content Rules

- Add one topic at a time.
- Prefer clear notes over long theory dumps.
- Every practice example should have a simple command to run it.
- When adding a topic, update the related language roadmap only if it helps the learner.
- Keep beginner, practice-first flow.

## Quality Rules

- Do not claim something was tested unless it was actually run.
- Do not leave commented-out code.
- Do not leave TODO placeholders unless explicitly requested.
- Avoid unrelated refactoring.
- Keep names simple and consistent.
- Never commit secrets, tokens, API keys, or credentials.

## Git Rules

- Do not commit automatically unless explicitly requested.
- Do not push automatically unless explicitly requested.
- Never force-push unless explicitly requested.
- Do not rewrite Git history unless explicitly requested.
- Keep commit messages concise and professional.
- Do not add AI co-author or AI attribution trailers.


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
    01-go-basics.md ... 05-production.md
    SOURCES.md
    practice/
    project/
    check.sh
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

- Keep Go explanations clear and sections short; cover topics in depth within the existing guides.
- Use the existing `practice/` and flat `project/` modules.
- Keep the five numbered notes as the main reading order.
- Put exercises and debugging beside the relevant explanation; do not recreate separate per-topic README, EXERCISES, INTERVIEW, or DEBUGGING files.
- Add frameworks and infrastructure when their learning exercise needs them.
- Keep learner completion separate from supplied material and passing tests.

Example:

```bash
cd languages/go/practice
go run ./basics
go run ./packages
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


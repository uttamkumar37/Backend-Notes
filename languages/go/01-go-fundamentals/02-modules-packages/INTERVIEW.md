# Lesson 2 — Five knowledge checks with explanations

Try answering independently before reading. Point to the actual files and command output as evidence.

## 1. Which packages belong to this module, and where does the library import resolve?

The two production packages are the root executable and `greeting/`. The module path is `example.com/go-lab/module-basics`, so adding `/greeting` yields the library's full import path. Go finds that package in the active local module's `greeting` directory. It does not search for a class named `greeting` or download our own package from the example domain. The external `greeting_test` package supports testing the library's public API; it does not require another module.

**Java bridge:** distinguish dependency-project identity from package/source organization rather than treating every package as a separate Maven project.

## 2. Why can `main` call `greeting.Welcome`, and what if it becomes lowercase?

`main.go` imports the package, and uppercase `Welcome` is exported. Calling a lowercase `welcome` from a different package fails even if its types and body are unchanged. A lowercase declaration is visible to other source files of its own package, so “unexported” does not mean private to a file. The package name and function name play different roles in the qualified call.

**Common mistake:** using a `public` keyword borrowed from Java, or assuming import grants access to every name.

## 3. What do `go.mod` and `go.sum` record, and why is a test helper a direct requirement?

`go.mod` identifies the module, Go baseline, and dependency requirements. `go.sum` records integrity hashes for module content/metadata; it does not itself choose all dependency versions. Test imports count when maintaining a module's requirements, so our direct use of `cmp` in a test keeps `github.com/google/go-cmp` in `go.mod`. The two checksum entries cover content and `/go.mod` metadata. Lesson 1 needs no `go.sum` because it uses only standard packages.

**Common mistake:** assuming a `go.sum` entry guarantees that code is safe, up to date, or even currently imported. Dependency review and vulnerability checks answer different questions.

## 4. How do tidy, download, and verify differ?

`go mod tidy` reconciles requirements/checksums with source and tests and may edit dependency files. `go mod download` populates the dependency cache and may access the network. `go mod verify` checks cached module content for changes after download; it does not validate our greeting's behavior. Run each from the intended module, review changes, and use normal tests for correctness.

**Common mistake:** clearing every cache or upgrading all dependencies before reading the first error. A wrong local import path is usually fixed in source, not by downloading an unrelated module.

## 5. Which modules does the workspace select, and why also test with `GOWORK=off`?

`workspace/go.work` selects Lesson 2 at `..` and Lesson 1 at `../../01-first-program`, relative to the workspace file. It groups two local modules without merging their metadata. The explicit test patterns cover both member directory trees. `GOWORK=off` tests a module using its own declared configuration; it helps detect reliance on a developer's workspace that another build environment may not have.

**Microservice connection:** several services can share one module, and one repo can have several modules. Choose boundaries for clear ownership/versioning needs; neither service count nor folder count mandates a new `go.mod`.

**Pass criterion:** trace imports to disk, explain generated metadata, run independent and workspace tests, and repair a documented fault. Record confidence/gaps and confirm Lesson 2 before moving to Lesson 3; Phase 1 remains in progress.

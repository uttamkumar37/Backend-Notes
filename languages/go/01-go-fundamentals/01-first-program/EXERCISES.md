# Lesson 1 exercises

Allow 20–30 minutes. Start in this module's directory, with [README.md](README.md) open. Keep all edits within this lesson, and restore its supplied behavior before finishing so the checked-in example/test output remains consistent.

## 1. Predict before running

Without executing the program, write the values of `language`, `goal`, and `message` immediately after `message` is declared. Describe the call order and whether `welcomeMessage` writes anything to the terminal.

Then run:

```bash
GOTOOLCHAIN=local go run .
```

**Verify:** output matches the two lines in the lesson. Explain the difference between a returned string and printed text.

## 2. Change a goal using assignment

Immediately after `goal := "backend engineering"`, insert `goal = "cloud-native services"`. This changes an existing variable, so use `=`. Keep the assignment before `message` is created.

**Predict:** the first output line becomes `Learning Go for cloud-native services.`; the second line stays the same.

Run the program, then `GOTOOLCHAIN=local go test -v ./...`.

**Verify:** the parameter test still passes because its inputs are independent of `main`. The output example fails because its expected line still describes the original goal. This is a deliberate contract mismatch, not a broken Go installation. Update only that expected line to reflect the intentional change and confirm both checks pass.

Now move the assignment **after** `message := welcomeMessage(language, goal)`, before printing. Predict the first line again. `message` already contains the earlier result; assigning a different goal does not recompute it. Run to confirm, then restore the supplied source and example expectation manually.

## 3. Add an independent input check

In `main_test.go`, add `TestWelcomeMessageUsesAnotherGoal`, using the same test structure already explained. Call the function with language `"Go"` and goal `"API development"`; expect `Learning Go for API development.`. Do not change implementation to hardcode your test values.

**Verify:** `GOTOOLCHAIN=local go test -v ./...` lists your new test and it passes. Explain the function's parameters and why `t` is supplied by the test runner. This is practice writing the existing pattern; later lessons introduce broader test-case organization.

## 4. Observe build versus run

Use the build commands in README to produce `bin/first-program`. Run it once, then temporarily change the goal in source and run `go run .` with the local toolchain setting.

Run the old `./bin/first-program` again **before** rebuilding. Its output is unchanged because editing source does not rewrite an existing executable. Rebuild and verify the new goal appears. Restore the supplied goal, rebuild, and leave tests passing.

**Verify:** describe which command compiles and launches, which only compiles, and which launches an already-built file. Build artifacts are ignored under `bin/`.

## 5. Diagnose and explain

Complete the argument-order drill in [DEBUGGING.md](DEBUGGING.md). Answer the five checks in [INTERVIEW.md](INTERVIEW.md) before reading their explanations.

Record what you changed, actual output, test command/result, one failure's cause/fix, confidence, and remaining gaps in [PROGRESS.md](../../PROGRESS.md). Restore intentional drill faults; keep your additional passing test if desired. Do not mark the entire phase complete.

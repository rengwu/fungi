---
title: Sessions, results and recovery
updated: 2026-10-22
---
# Sessions, results and recovery

**ID: FEAT-sessions**

Fungi owns one interaction and one immutable terminal result per session. This implements the [session contract](../contract.md#result-and-lifecycle) for the accepted first milestone.

## Entry points

Use `open`, `collect` and `reopen` in the Go CLI. An executable helper uses `open --exec HELPER --delivery wait -- ARGS`. See [setup](../verification/setup.md) for builds, isolation and fixtures. No AI account is needed for core checks.

## Expected outcomes

- Submission persists data and artifact bytes before acknowledgement.
- Close, cancellation and failure remain distinct from a submitted response.
- Repeated identical submission returns the saved result; conflicting terminal events cannot replace it.
- Collect is repeatable and non-destructive. A held OS lock identifies a live owner; a saved launch without an owner/result is interrupted.
- Losing the waiting CLI does not cancel the owner. The worker uses a detached process, ordinary pipes and an OS file lock; there is no local HTTP server.
- Reopen creates a new linked session and restores the latest acknowledged JSON/image draft.
- Native helper output is drained before interpreting exit. The helper need not depend on Electron.
- Malformed output or a closed stdout ends the helper protocol and terminates the helper, even if its process would otherwise keep running. Persisted JSON is bounded including its final newline, so acknowledged boundary-sized records remain readable.

## Verification recipe

Run `go test -race -v . ./examples/helper` from the repository. The tests build real CLI/helper executables in a temporary directory. They exercise concurrent sessions, lost waiting clients, helper exit, repeated collection, duplicate submissions, file confinement, saved image drafts, interrupted recovery and failed publication.

For a visible result:

```sh
go build -o verification-output/helper ./examples/helper
FUNGI_STATE_DIR="$PWD/verification-output/sessions" ./fungi open --exec "$PWD/verification-output/helper" --delivery wait
```

Read the returned artifact path. Run `collect` with the exact returned session and the same state directory; it must return the same envelope. Keep the result/draft directories until evidence is saved. Never use repeated collection to wait for human input.

## Coverage and limitations

The macOS and Linux arm64 race-enabled suites passed. Hosted core execution also passed on Windows and Linux x64; the follow-up report records target-specific failures and repairs. Failure injection covers publication before acknowledgement; the UI separately exercised an actual filesystem permission error. One user-controlled macOS sleep/wake cycle preserved the note and image, accepted a new edit and returned the result. These tests do not establish power-loss or storage-hardware durability. Windows skips directory flushing after flushing and renaming each file; no directory-entry durability across power loss is promised.

## Verification history

- [2026-10-22 verified baselines](../verification/runs/2026-10-22-baselines.md) — saved renderer startup failure and identical collection, plus 18 core tests on four hosted targets.

- [2026-10-22 hosted checks](../verification/runs/2026-10-22-hosted.md) — Windows/Linux/macOS core execution and target-specific renderer findings.

- [2026-10-22 platform verification](../verification/runs/2026-10-22-platform.md) — malformed/closed helper output and exact JSON-limit regressions, plus the macOS core race suite.

- [2026-10-20 skill proof](../verification/runs/2026-10-20-skill-proof.md) — real skill recovery, Linux arm64 core execution, one macOS sleep/wake cycle and a completed two-hour helper/owner/callback run.
- [2026-10-20 no-setup compatibility](../verification/runs/2026-10-20-no-setup.md) — self-contained collection receipts and real native artifact sessions.

- [2026-10-20 first proof](../verification/runs/2026-10-20-first-proof.md) — core checks and coverage gaps.

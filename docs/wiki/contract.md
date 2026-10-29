---
title: Session engine contract
updated: 2026-10-28
---
# Session engine contract

The public protocol 2 reference is embedded in `fungi help helper`, `delivery`,
`cli`, `installation`, `examples` and `migration`. The [accepted engine-only
scope](planning.md#accepted-engine-only-scope-on-2026-10-28) assigns all UI
implementations to other projects.

## Process ownership

Each interaction has one detached Go session owner and an external executable
helper communicating through stdin/stdout. The helper supplies its own runtime;
Fungi launches it without a shell and removes the caller's CODEX_THREAD_ID from
its environment. Helpers are trusted native programs, not sandboxed plugins.

The engine serializes helper messages and holds an exclusive session file lock.
Startup readiness is bounded; no human timeout follows readiness. No shared
service or renderer remains after a completed interaction.

## Result and lifecycle

Input is snapshotted at launch. Drafts are optional complete snapshots of JSON
and attachments. Results distinguish submitted, cancelled, closed and failed.
Only submission has answer data. Closing an interface must be reported by its
helper and never implies approval. Each message is acknowledged after durable
persistence. Identical duplicate terminal results return the committed record;
a different terminal result is rejected.

Artifacts are copied from confined paths under outputDir, hashed and saved
before acknowledgment. Read advertised size limits. Final responses persist
independently of drafts. Session storage is separate from managed installations;
collection is non-destructive and version pruning does not remove session data.

## Failure and recovery

A failed helper start or exit before readiness produces a durable
helper_startup_failed result. Exiting after readiness without a terminal message
produces helper_exited. Malformed output or blocked protocol writes cannot keep
the owner indefinitely waiting for a broken helper. Committed results survive
notification failure and helper exit.

Explicit reopen creates a new session from original input, helper arguments and
the last acknowledged draft/artifacts. Active and submitted sessions cannot be
reopened. Original records remain unchanged. There is no automatic restart or
remote desktop transport. External helpers/assets must remain available.

## Delivery boundary

Bind the exact caller at launch, outside helper data. Auto uses the supported
Codex callback or reports manual fallback. A host-managed pending/background
call may use wait. Manual delivery requires the user's done cue before one
collection. Never poll while a person responds. Saved result state is separate
from notification state; accepted notification is not proof of host consumption,
and ambiguous callbacks are not automatically resent.

## Verification boundary

The [feature map](features/index.md), [setup](verification/setup.md) and
[headless run](verification/runs/2026-10-28-headless.md) record exact checks.
Legacy UI reports remain historical. Protocol 2 explicitly withdraws the old
page/CLI contract; protocol-1 session data requires its original executable.
Public message names/types/meanings are stable within protocol 2; helper schemas
and internal storage layout are outside that promise.

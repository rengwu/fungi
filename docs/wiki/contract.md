---
title: Session engine contract
updated: 2026-10-27
---
# Session engine contract

This page records the engine's responsibilities and invariants. The maintained
public API, commands and schemas live in the embedded manual:

| Reference | Scope |
| --- | --- |
| `fungi help cli` | Session, discovery and installation commands |
| `fungi help html` | Four-method page API, limits and protocol 1 compatibility |
| `fungi help delivery` | Launch receipts, outcomes, artifacts and recovery |
| `fungi help helper` | Native executable protocol and acknowledgments |
| `fungi help form` | Experimental flat JSON form schema |
| `fungi help installation` | Packages, runtimes, updates and rollback |

Use the [authoring guide](authoring.md) and [feature map](features/index.md)
for entry points and verification recipes. [Current decisions](planning.md#current-decisions)
establish the scope; dated reports retain their original build boundaries.

## Process ownership

Go owns sessions, persistence and delivery. Each interaction has a detached
owner process and a private state directory; an OS file lock establishes
liveness without relying on a reusable PID. Pipes connect the launcher,
owner and renderer/helper. There is no local HTTP listener or permanent daemon.
Losing the launcher's waiting connection does not itself cancel the session.

Wails supplies the native web window using the OS webview. Skills supply local
HTML/assets or a trusted executable helper. The small built-in JSON renderer
uses the same result path. Wails pages and native helpers are trusted local
application code; browser storage isolation is best effort. Fungi keeps its
own sessions, drafts, results and artifacts separate.

## Result and lifecycle

A session has at most one authoritative terminal result: `submitted`,
`cancelled`, `closed` or `failed`. The owner serializes competing outcomes.
A committed submission survives window exit and notification failure.
Identical repeated submissions return the saved receipt; conflicting repeats
fail without changing the result or sending another completion notice.

Persist referenced artifact bytes before acknowledging a draft or result.
Failed persistence must not report success or replace previously acknowledged
state. Each draft/submission supplies its complete data and attachment list.
Artifact paths are confined to the helper output directory before the owner
copies them into content-addressed storage. The advertised limits apply before
success. JSON response data remains opaque to the engine.

Reopening creates a new linked session with the current caller, saved launch
input and the latest acknowledged draft; the old session remains unchanged.
An active or submitted session cannot be reopened. Source files must still be
available. Unsaved UI state and historical copies of source assets are not
reconstructed. Collection is repeatable and non-destructive; saved sessions
remain until explicitly removed.

## Failure and recovery

A helper exit without a terminal result is a failure. If an owner exits
without a committed result, collection reports an interrupted session and
leaves recovery explicit. Startup and protocol operations have bounded
timeouts; ordinary human waiting has no Fungi deadline after readiness.

Window close records `closed` and preserves acknowledged drafts. The bridge
orders save/submit/cancel calls before asynchronous attachment encoding;
a rejected operation does not prevent a later retry. There is no automatic
DOM-state capture, heartbeat or restart. A Wails web-content crash may leave
a blank window alive: close it before explicit collection and reopening.
See [interface recovery](features/interfaces.md#verification-recipe).

## Delivery boundary

Bind the caller at launch, outside page-controlled data. Use one return route:
an exact-caller callback, a host-managed pending/background call, or the manual
“done” fallback. Never infer a destination from the latest conversation or poll
collection while a person is responding. No harness setup is required.

Keep saved-result state separate from notification state. An armed route is
not guaranteed delivery; accepted notification is not proof of consumption.
Do not blindly resend after an ambiguous acknowledgment. Preserve the result
for exact-session collection regardless of notification outcome. Reopening
captures the caller anew; it never reuses the old session's callback target.

## Verification boundary

[Sessions](features/sessions.md), [interfaces](features/interfaces.md),
[delivery](features/delivery.md) and [installation](features/installation.md)
record their expected behavior and executable checks. The
[platform record](platforms.md) distinguishes current Wails observations from
historical Electron proofs. The [October 27 contract and AFK run](verification/runs/2026-10-27-afk-contract.md)
records local protocol, callback, two-hour idle and sleep/wake checks.

These results qualify their tested builds, platforms and harness modes.
Clean consumer-machine setup, physical Windows/Linux Wails input and broader
idle/failure behavior remain separate verification work. Protocol 1 preserves
documented names, types and meanings; the experimental form schema, renderer
internals and on-disk layout are outside that public compatibility promise.

---
title: Authoring contract and AFK verification
updated: 2026-10-27
result: Passed
---
# Authoring contract and AFK verification

**Passed for the selected local contract, callback, two-hour idle and sleep/wake
checks.** The user requested these tasks
while clean-machine and physical Windows/Linux verification are unavailable.
Scope: [discovery](../../features/discovery.md),
[interfaces](../../features/interfaces.md) and [delivery](../../features/delivery.md).

## Build and changes

Local package `v0.1.0-local.20261027` uses Wails 2.17.0 on macOS arm64, with
Codex CLI 0.164.0. The [BUILD record](../../evidence/2026-10-27-afk-contract/BUILD.json)
records base `a94e3ffd5dd1625cca724f74f8f40b3db5e894e7`, dirty source and binary/helper
hashes. [Local checks](../../evidence/2026-10-27-afk-contract/local-checks.json)
record source hashes, environment, completed session identities and cleanup.

The manual now states the protocol 1 compatibility boundary: documented page
methods, session commands and messages preserve names, types and meanings;
optional metadata may be added. Renderer internals and the experimental form
schema remain excluded. Discovery examples use auto delivery so they do not
assume a host can sustain a pending wait call. No new API layer was introduced.

Bridge regression tests cover repeated readiness, empty/binary attachments,
additive metadata, falsy/nested JSON values and acknowledged submission results.
Existing ordering, failure/retry, saved-result and callback-uncertainty checks
remain in place. The renderer and session engine implementations are unchanged.

## Results

| Scenario | Expected | Observed | Result |
| --- | --- | --- | --- |
| Core and bridge regression | Preserve public behavior and avoid implicit host wait assumptions | `go vet ./...`, `go test -race ./...` and all eight Node tests passed; offline discovery checks require protocol 1 and auto launch examples | Passed |
| Native page contract | Real page calls persist, cancel, reopen and return the saved PNG | All 21 [renderer checks](../../evidence/2026-10-27-afk-contract/renderer.json) passed with no exclusions, including draft and Blob restoration | Passed |
| Older skill compatibility and fallback | Unchanged authored page works; missing exact caller selects manual | Native priority page moved Verify to first, submitted and closed; one collection returned exact order and matching SVG, with no draft | Passed |
| Short Wails idle | One local timer returns unchanged state and PNG after five seconds | [Result](../../evidence/2026-10-27-afk-contract/idle-short.json) reports 5.186 seconds, preserved input and expected marker; PNG bytes/hash verified | Passed |
| Short exact-caller callback | Notification wakes this conversation, followed by one collection | Notice resumed this conversation after the prior turn ended. One collection returned the expected marker after 10.031 seconds, preserved input and a verified 3,278-byte PNG | Passed |
| Two-hour Wails idle | No model polling; one completion cue and one collection after the timer | The page timer completed after 7,200.196 seconds of performance-clock time and 7,667.373 seconds of wall time, spanning two OS-recorded sleep/wake cycles. The callback resumed this conversation and one collection verified unchanged note/input and PNG | Passed |
| User-controlled sleep/wake | Saved state survives actual sleep; post-wake edit and PNG return | macOS recorded clamshell sleep at 01:23:15 and wake at 01:29:10 (+08:00), 355 seconds later. The user appended “edited after sleep”; the exact-conversation callback resumed the agent and one collection verified the note, preserved input and unchanged PNG hash | Passed |

Older page source is copied unchanged from the
[earlier authored example](../../evidence/2026-10-26-discovery/variation.html.txt).
Its [manual receipt](../../evidence/2026-10-27-afk-contract/manual-receipt.json),
[new result](../../evidence/2026-10-27-afk-contract/prior-skill-result.json) and
[SVG bytes](../../evidence/2026-10-27-afk-contract/priority.svg.txt) preserve the check.
This is a maintainer-operated compatibility test, not an independent fresh-agent
authoring trial.

The [sleep/wake result](../../evidence/2026-10-27-afk-contract/sleep-wake-collect.json)
and [verified PNG](../../evidence/2026-10-27-afk-contract/sleep-wake.png) match the
prepared baseline, with the user's new note. The
[observation record](../../evidence/2026-10-27-afk-contract/sleep-wake-observation.json)
records one collection and actual callback consumption. The
[power-log excerpt](../../evidence/2026-10-27-afk-contract/sleep-wake-power.log)
contains only Sleep/Wake events within the test interval. Total session wall
time was 398.133 seconds; the page's performance clock reported 45.742 seconds.
The OS events, rather than that timing difference alone, establish actual sleep.

The [long-idle result](../../evidence/2026-10-27-afk-contract/idle-long-collect.json)
returned the original note and [PNG](../../evidence/2026-10-27-afk-contract/idle-long.png).
Its [observation record](../../evidence/2026-10-27-afk-contract/idle-long-observation.json)
confirms one collection after the completion notice; the
[delivery record](../../evidence/2026-10-27-afk-contract/idle-long-delivery.json)
records acceptance. The [power-log excerpt](../../evidence/2026-10-27-afk-contract/idle-long-power.log)
contains the 355-second clamshell sleep and a later 117-second idle sleep within
the interval. No model status checks were made while awaiting this result.

## Completed session records

All used the local package under
`dist/release-v0.1.0-local.20261027/fungi-darwin-arm64/` and separate state roots
under `verification-output/afk-contract-2026-10-27/`. Each receipt's exact
collect array and environment was used once, on its completion notice.
All results are collected; no test session remains pending. The receipt paths
below identify the original run and are no longer runnable after cleanup.

| Test | Receipt and current evidence |
| --- | --- |
| Short callback `f484f5719d260079a748467a5fba8418` — collected once | [Receipt](../../evidence/2026-10-27-afk-contract/callback-short-receipt.json), [accepted delivery](../../evidence/2026-10-27-afk-contract/callback-short-delivery.json), [collected result](../../evidence/2026-10-27-afk-contract/callback-short-collect.json), [verified PNG](../../evidence/2026-10-27-afk-contract/callback-short.png) |
| Two-hour idle `7126d5d6c8ad8f797bd8fcbc54de704b` — collected once | [Receipt](../../evidence/2026-10-27-afk-contract/idle-long-receipt.json), [baseline](../../evidence/2026-10-27-afk-contract/idle-long-baseline.json), [collected result](../../evidence/2026-10-27-afk-contract/idle-long-collect.json) |
| Sleep/wake `f930bce64f4158e7b1d5eb5d40301b54` — collected once | [Receipt](../../evidence/2026-10-27-afk-contract/sleep-wake-receipt.json), [baseline](../../evidence/2026-10-27-afk-contract/sleep-wake-baseline.json), [collected result](../../evidence/2026-10-27-afk-contract/sleep-wake-collect.json) |

The [short fixture source](../../evidence/2026-10-27-afk-contract/idle-short.html.txt)
and [extended fixture source](../../evidence/2026-10-27-afk-contract/idle-extended.html.txt)
differ only in the displayed heading for manual sleep/wake. The extended fixture
uses one timer for automatic submission, or an explicit button for sleep/wake;
it has no polling loop or heartbeat. No system sleep settings were changed.
After every result was collected and the package processes had exited, the
disposable package and test directories were removed. Collected results,
artifact bytes, baseline records, fixture sources and build hashes remain in
the evidence directory. The user's installation was unchanged.

## Limits

No clean-machine, Windows/Linux desktop or independent-agent qualification is
claimed. These bounded tests cannot establish arbitrary AFK duration. OS sleep and
background timer throttling can affect wall-clock duration; the timer alone is
not evidence of sleep. The previously observed intermittent Mac media failure
is not resolved by this passing renderer run. No release or user installation
was changed.

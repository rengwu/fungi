---
title: Callback across an intervening chat message
updated: 2026-10-27
result: Passed
---
# Callback across an intervening chat message

**Passed for the selected local Codex callback scenarios.** A human response
returned after an intervening chat message, and a separate timed submission
completed during a running tool call. Both callbacks resumed the original
conversation after its active turn ended, without “done,” completion polling,
a replacement waiter or relaunching either session. Scope:
[FEAT-delivery](../../features/delivery.md) and
[FEAT-sessions](../../features/sessions.md).

## Build and method

The [build](../../evidence/2026-10-27-interleaved-message/BUILD.json) was made
from clean commit `c9be5eb8f41bd05043da2b2c3d550f80f57ae123` on macOS 27 arm64,
using Wails 2.17.0 and Codex CLI 0.164.0. The
[environment](../../evidence/2026-10-27-interleaved-message/environment.json),
[version](../../evidence/2026-10-27-interleaved-message/version.json) and
[runtime check](../../evidence/2026-10-27-interleaved-message/doctor.json)
identify the tested package. No runtime code or bridge API changed.

The first session used the unchanged
[reply example](../../evidence/2026-10-27-interleaved-message/reply.html.txt) and
[test instructions](../../evidence/2026-10-27-interleaved-message/input.json).
Its [receipt](../../evidence/2026-10-27-interleaved-message/receipt.json) armed an
exact-conversation callback. The user sent “continue the test” while the form
was open, then submitted in response to the agent's cue. One collection on the
automatic notice returned `hihi this is the short note` from the original session.

That result was saved about eight seconds before the core-check subprocess
started. It verifies the intervening-message path and deferred callback, but
not submission during that subprocess. A second, independent session used the
existing [idle page](../../evidence/2026-10-27-interleaved-message/active-tool-fixture.html.txt)
with a [five-second timer](../../evidence/2026-10-27-interleaved-message/active-tool-input.json).
The driver launched it once and immediately began a bounded Go test in the
same tool call. No additional human action was needed. Its
[receipt](../../evidence/2026-10-27-interleaved-message/active-tool-receipt.json)
also armed the original conversation's callback.

## Results

| Scenario | Expected | Observed | Result |
| --- | --- | --- | --- |
| Session across another user message | Original form returns a response after chat continues | Session `f1837b1eb7b0d0777ab64231c5e5682f` returned the user's note after “continue the test”; [result](../../evidence/2026-10-27-interleaved-message/result.json) | Passed |
| Automatic return after the active turn | Callback resumes the original conversation without “done” | First callback arrived in a subsequent turn; collected exactly once. [Delivery](../../evidence/2026-10-27-interleaved-message/delivery.json) was accepted and the recorded caller matched | Passed |
| Completion during a running tool | Submission persists while the agent's bounded command remains active | Timed page finished at 07:40:22.844 UTC, inside the 07:40:17.699–07:40:29.696 command interval; [command timing](../../evidence/2026-10-27-interleaved-message/second-active-work.json) and [result](../../evidence/2026-10-27-interleaved-message/active-tool-result.json) | Passed |
| Queued callback and artifact return | Busy-turn completion later resumes the same conversation with intact data | Second callback arrived after turn end; one collection verified the session, marker, input and 3,278-byte PNG SHA-256; [delivery](../../evidence/2026-10-27-interleaved-message/active-tool-delivery.json), [PNG](../../evidence/2026-10-27-interleaved-message/active-tool-image.png) | Passed |
| Bounded work commands | Verification work finishes without routing its fixtures into this conversation | [Race suite](../../evidence/2026-10-27-interleaved-message/active-work.log) and [wait-path test](../../evidence/2026-10-27-interleaved-message/second-active-work.log) passed with caller routing removed from their test environments | Passed |

## Timing and interpretation

[Observations](../../evidence/2026-10-27-interleaved-message/observations.json)
retain both receipts, event order, returned data and timing records. The first
result file was saved at 07:35:16.604 UTC and its accepted-delivery file at
07:35:16.970. The first work command ran from 07:35:24.659 to 07:36:35.634.
The user's message and the assistant's submit cue preceded the first command;
the message's recorded timestamp is when the tool wrote the observation,
not an exact timestamp of the chat message itself.

For the second session, the result file was saved at 07:40:22.859 UTC and
accepted delivery recorded at 07:40:23.073, both inside the active command.
The page's own `finishedAt` independently places submission in that interval.
File times are observations of the original session files, not copied evidence.
Callback observation times record when the handling turn wrote its evidence;
they are not precise notification-latency measurements.

In both cases the completion notice was consumed after the assistant ended its
active turn. The observed behavior is queued completion and automatic resumption,
not immediate interruption of the running tool. Queue acceptance alone was not
counted as a pass: each notice actually reached this conversation and caused
one exact-session collection.

## Limits and retained state

This qualifies the tested Codex callback route on this build. Claude and other
hosts' background-wait routes, explicit task cancellation, host shutdown and
callback-receiver failure were not exercised. The human-message scenario and
mid-tool submission were separate sessions; the first did not overlap the
Go subprocess itself.

Both completed sessions, their inputs/artifacts and the disposable package
remain under ignored local build/test directories for inspection. No matching
worker, renderer or waiter remained after collection. Durable evidence is in
this report's evidence directory. No runtime changes, commit, release or user
installation changes were made by this test.

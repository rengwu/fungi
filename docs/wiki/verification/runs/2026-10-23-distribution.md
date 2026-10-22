---
title: Command installation and update proof
updated: 2026-10-23
result: Untested
---
# Command installation and update proof

**Untested overall: the newest Mac job and its single retry never obtained
a hosted runner.** Automatic retries have stopped; no callback is pending.
Windows and both Linux jobs now pass on `6d587ac`; the initial four-platform
run, Mac runtime-sharing CI and selected local Mac scenarios also passed.
This run covers [installation](../../features/installation.md)
and preservation of [sessions](../../features/sessions.md) through updates.
It does not establish a public release or clean consumer-machine support.

## Build and environment

Local base `c768218`, with implementation edits; [source hashes](../../evidence/2026-10-23-distribution/source-hashes.json)
record the dirty source checkpoint. macOS 27.0 arm64, Go 1.26.4, Node 22.19.0
and Electron 44.9.0. [Package A](../../evidence/2026-10-23-distribution/package-a.json)
and [package B](../../evidence/2026-10-23-distribution/package-b.json) identify
the actual executable hashes and sizes. Their versions are `v0.1.0-install.1`
and `v0.1.0-install.2`. Later packaging changes increase archive compression
and substitute a compact portable README; those two archives predate these
packaging-only changes. [Initial CI run 39843208109](https://github.com/rengwu/fungi/actions/runs/39843208109)
passed on `8d4f23b6900f736cc7ee9bc4b05d222a87aa1704`; its
[metadata](../../evidence/2026-10-23-distribution/ci-39843208109.json) and
[log](../../evidence/2026-10-23-distribution/ci-39843208109.log) preserve exact job scopes.

The file-sharing follow-up uses local packages
[`v0.1.0-install.3`](../../evidence/2026-10-23-distribution/package-c.json) and
[`v0.1.0-install.4`](../../evidence/2026-10-23-distribution/package-d.json).
Their [source checkpoint](../../evidence/2026-10-23-distribution/runtime-reuse-source.json)
identifies the edits atop `8d4f23b`. These include the final compact README
and maximum Go archive compression. Package C is 329,066,779 bytes unpacked
and 137,878,147 bytes compressed (313.8 and 131.5 MiB).

Local installation interactions use synthetic fixtures and isolated session/install
directories. These fixtures submit no real questionnaire and queue no agent
callback. A separate CI-wait helper exercises an actual callback, recorded
below. No release is published and no shell profile/PATH is changed. Installed
CLI children use an empty PATH; the shell bootstrap uses standard OS tools.

## Scenarios

| Scenario | Result | Observation |
| --- | --- | --- |
| Core regression and archive tests | Passed | [Race-enabled log](../../evidence/2026-10-23-distribution/go-tests.log): 20 root tests and two archive test groups; Go vet and seven Node tests also passed |
| Managed A → B → C, recovery and cleanup | Passed | Actual CLI/helper processes; stable receipt path, failed checksum preserves A, concurrent mutation rejected, running/recoverable A retained, draft restored through C, old result immutable, artifact retained, rollback and later pruning succeed |
| Invalid downloads/archive paths | Passed | Oversize, partial and missing HTTP responses rejected; ZIP traversal/absolute/drive paths, duplicate files, escaping links and TAR hardlinks rejected; normal symlinks/modes survive Unix round trips |
| Shell bootstrap on development Mac | Passed | [Download evidence](../../evidence/2026-10-23-distribution/mac-install.json): exactly one matching archive plus checksum fetched from a local fixture server, checksum checked, space-containing destination installed |
| Open Electron session across update | Passed | A acknowledges draft and PNG, B activates, A submits, original stable receipt collects the matching saved answer; [PNG](../../evidence/2026-10-23-distribution/submitted.png) hash verified, rollback selects A |
| Installed renderer regression | Passed | [All 21 checks](../../evidence/2026-10-23-distribution/mac-renderer.json), including WebGL, cancel/reopen, restored Blob, submission and hash through the stable installed command |
| Hosted macOS, Windows and Linux install/update paths | Passed | Initial CI completed all four jobs, including actual platform bootstraps and local archive updates; Ubuntu 24.04 intentionally omits UI |
| Runtime sharing, changed files and cleanup | Passed | [Race-enabled follow-up](../../evidence/2026-10-23-distribution/runtime-reuse-tests.log): unchanged file shares an inode, changed file stays separate, removing the old directory preserves the new runtime; managed lifecycle regression passes |
| Actual Mac runtime sharing during update | Passed | [Install evidence](../../evidence/2026-10-23-distribution/mac-runtime-reuse.json): runtime executable shares an inode; open session submits with matching PNG hash; rollback succeeds |
| Renderer regression with shared runtime | Passed | [All 21 checks](../../evidence/2026-10-23-distribution/mac-runtime-renderer.json) pass through the stable installed command after the sharing update and rollback, including WebGL and draft/Blob recovery |
| Hosted runtime-sharing follow-up | Passed | Mac passes on `af2dbee`; [Windows](../../evidence/2026-10-23-distribution/windows-runtime-reuse.json), [Ubuntu 22.04](../../evidence/2026-10-23-distribution/linux22-runtime-reuse.json) and [Ubuntu 24.04 core](../../evidence/2026-10-23-distribution/linux24-runtime-reuse.json) pass on `6d587ac` |
| Windows disconnect regression follow-up | Passed | The handshake correction passes on Windows and both Linux runners at `8bd24a0`; the following wiki check stopped their installer/renderer stages |
| Fresh-checkout generated wiki | Passed | Corrected generation passes a local export of staged files and the Windows/Linux hosted checks; the CRLF failure remains recorded below |
| Latest complete CI matrix | Untested | [Run 39848109685](https://github.com/rengwu/fungi/actions/runs/39848109685) passes Windows and both Linux jobs on `6d587ac`; attempts 1 and 2 Mac ran zero steps due to hosted capacity; no further retry is armed |
| CI-wait callback through managed installation | Passed | [Callback evidence](../../evidence/2026-10-23-distribution/ci-callback.json): both CI attempts saved their failure result, delivered a Codex notice to the originating conversation and were each collected once |

The first unsafe-archive test attempt reused a ZIP header object and failed in
the test writer. A fresh duplicate header corrected the fixture. The first
managed-path assertion also exposed macOS `/var` versus `/private/var` aliases;
package-reference matching now resolves existing symlinks before retaining
assets. The rerun passes. No failed product assertion was simply removed.

[Follow-up CI 39844997609](https://github.com/rengwu/fungi/actions/runs/39844997609)
on `af2dbee1ea02bb0d2fe36f892e7fab985a6663c5` failed overall because the
Windows driver terminated PowerShell at exactly 180 seconds (`killed: true`,
`SIGTERM`); it had requested both release files but produced no stage output.
The update/sharing assertions therefore did not execute there. Preserve the
[failed log](../../evidence/2026-10-23-distribution/ci-39844997609.log) and
[job metadata](../../evidence/2026-10-23-distribution/ci-39844997609.json).
The prior Windows run took about 165 seconds for this entire scenario.

The correction adds bootstrap stage messages, records elapsed time and failure
signal metadata, and gives Windows extraction a 300-second fixture bound.
It also replaces per-chunk PowerShell progress with the concise stage messages;
[Microsoft's module issue](https://github.com/microsoft/PowerShellForGitHub/issues/227)
documents this source of download overhead. The original timeout's exact
phase is unknown, so progress rendering is a plausible contributor, not an
established root cause. All functional assertions remain unchanged. No human
wait timeout or host security setting changed.

The [next Windows attempt](https://github.com/rengwu/fungi/actions/runs/39846391476)
on `58edcc6` stopped before the installer in
`TestDisconnectedWaiterDoesNotCancel`: the test killed the CLI after one
second and expected its ten-second helper's result eleven seconds later.
The result was absent; [job log](../../evidence/2026-10-23-distribution/windows-120060671977.log)
preserves the failure. Those sleeps did not establish that the owner had
started or the helper had completed. Without worker diagnostics from that
attempt, delayed startup is an inference rather than a proven cause.

The correction uses the existing bounded helper gate: wait for an acknowledged
draft/readiness marker, kill the actual waiting CLI, assert no result exists,
release submission and require the saved submitted result. It adds worker
diagnostics on failure. The fixture's observable-condition deadline is thirty
seconds, matching the product startup bound; there is no fixed startup sleep
or human polling. The queued Mac job in that superseded run was cancelled;
both Linux jobs had already passed; [job metadata](../../evidence/2026-10-23-distribution/ci-39846391476.json)
records that cancellation. The new committed-source rerun will cover
all four targets.

The [local race-enabled handshake and managed-update retest](../../evidence/2026-10-23-distribution/fixture-handshake-tests.log)
passed before this follow-up was committed; it does not establish the Windows
result.

[CI 39847366853](https://github.com/rengwu/fungi/actions/runs/39847366853)
on `8bd24a0` passed core and Node tests on Windows and both Linux targets,
including the corrected disconnect assertion. All three then failed the wiki
freshness check, before renderer/installer stages. The [Windows log](../../evidence/2026-10-23-distribution/windows-120064065885.log)
and [metadata](../../evidence/2026-10-23-distribution/ci-39847366853.json)
record it. This was reproduced in a fresh local export of that commit:
the saved Windows job log had CRLF in the working tree when the wiki hashed
attachments, but Git stored LF per repository policy. Normalize the evidence
before rebuilding and verify a fresh staged checkout as well as the working
tree. No wiki assertion is removed. The superseded queued Mac job was cancelled.

The corrected Windows run completed the bootstrap in **11,773 ms**, then passed
the real Electron session, shared-file, collected PNG and rollback assertions.
Both Linux jobs also passed their declared scopes. The final
[attempt 1 metadata](../../evidence/2026-10-23-distribution/ci-39848109685.json)
records an overall CI failure because Mac was cancelled after fifteen minutes
without an assigned runner or any executed steps. GitHub's
[annotations](../../evidence/2026-10-23-distribution/mac-capacity-attempt1.json)
state that the job was not acquired by a hosted runner and identify macOS arm64
capacity constraints. This leaves the latest Mac scope untested; it is not a
failed product assertion. The single Mac-only retry on the same source also
ended after fifteen minutes without a runner or any steps. Its
[final metadata](../../evidence/2026-10-23-distribution/ci-39848109685-attempt2.json),
[Mac job record](../../evidence/2026-10-23-distribution/mac-job-attempt2.json) and
[annotations](../../evidence/2026-10-23-distribution/mac-capacity-attempt2.json)
preserve the same capacity diagnosis. The earlier Windows/Linux successes
remain valid; they were not rerun in attempt 2. There is no third automatic
attempt or pending watcher. The latest hosted Mac scope remains externally
blocked until runner capacity is available.

Fungi's separate CI-wait helper saved that failed CI conclusion and sent a
Codex notice to the originating conversation. The agent received it and
collected the exact session once. The helper's `watchError: exit status 1`
reflects the CI conclusion, not a failed Fungi delivery. The
[callback evidence](../../evidence/2026-10-23-distribution/ci-callback.json)
records both observed returns. The retry's notice also reached this
conversation, and its exact session was collected once. Each wait used one
local helper with no model polling. Both helpers have finished; no callback
is pending. This proves the installed callback path for these fixtures, not
callback delivery during an update. The final matrix result is not inferred
from the earlier Mac success.

## Size observations

[Filesystem measurement](../../evidence/2026-10-23-distribution/runtime-storage.json)
of the two retained Mac versions records 658,133,714 logical bytes and
335,661,895 bytes in unique regular files: **322,471,819 bytes (307.5 MiB)
shared**. These figures exclude the small stable launcher and filesystem
overhead; they are not download sizes. Electron's locales, GPU fallback and
license notices remain intact. File sharing does not reduce full-archive
download traffic or the space needed to stage an update. Unsupported hard
links fall back to independent verified copies.

The Windows hosted update shares **386,397,410 bytes (368.5 MiB)** of runtime
files; both Linux updates share **297,192,653 bytes (283.4 MiB)**. These are
retained-storage savings for an unchanged Electron version, not reductions
to first-download sizes.

## Limits and follow-up

The local HTTP fixture proves the bootstrap download path; updates use an
explicit local archive, not a live GitHub release. HTTP failures are tested
separately. A saved receipt was collected through the stable command; actual
agent queue delivery during an update is not claimed here. UI actions in this
run are automated page fixtures, not native human keyboard/mouse input.

Physical Windows and Linux installation, fresh-shell PATH integration,
downloaded/quarantined consumer packages, signing/notarization, release
publication, interruption during selector replacement and power-loss recovery
remain unverified. Existing platform and harness limits remain in force.
Verification was performed by the implementing agent; no independent review
is claimed.

---
title: Command installation and update proof
updated: 2026-10-23
result: Failed
---
# Command installation and update proof

**Failed pending retest of the Windows regression described below.**
The initial four-platform CI run and selected local Mac scenarios passed.
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

All local interactions use synthetic fixtures and isolated session/install
directories. No real questionnaire is submitted, no agent callback is queued,
no release is published and no shell profile/PATH is changed. Installed CLI
children use an empty PATH; the shell bootstrap uses standard OS tools.

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
| Hosted runtime-sharing follow-up | Untested | Mac and both Linux jobs passed on `af2dbee`; Windows bootstrap hit the fixture's 180-second timeout before installation completed. Corrected driver/bootstrap rerun pending |
| Windows disconnect regression follow-up | Failed | The subsequent Windows core run found no result after fixed startup/completion sleeps; the handshake-based correction awaits hosted retest |

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

The next Windows attempt on `58edcc6` stopped before the installer in
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
both Linux jobs had already passed. The new committed-source rerun will cover
all four targets.

The [local race-enabled handshake and managed-update retest](../../evidence/2026-10-23-distribution/fixture-handshake-tests.log)
passed before this follow-up was committed; it does not establish the Windows
result.

## Size observations

[Filesystem measurement](../../evidence/2026-10-23-distribution/runtime-storage.json)
of the two retained Mac versions records 658,133,714 logical bytes and
335,661,895 bytes in unique regular files: **322,471,819 bytes (307.5 MiB)
shared**. These figures exclude the small stable launcher and filesystem
overhead; they are not download sizes. Electron's locales, GPU fallback and
license notices remain intact. File sharing does not reduce full-archive
download traffic or the space needed to stage an update. Unsupported hard
links fall back to independent verified copies.

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

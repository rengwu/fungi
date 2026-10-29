---
title: Verified environments and portable distribution
updated: 2026-10-28
---
# Verified environments and portable distribution

Current source builds a headless Go engine with CGO disabled, protocol 2 and
package format 3. All UI implementations belong to other projects. Engine
packages have no GUI, webview or native renderer dependency. Release targets
remain Mac arm64, Windows x64 and Linux x64; cross-build support alone is not
runtime qualification.

The [headless run](verification/runs/2026-10-28-headless.md) records current local
execution, package checks and limitations. Hosted CI, publication, signing and
consumer-machine installation require new evidence. No 0.2 public release is
implied. Previous 0.1 alpha archives remain UI packages.

Everything below records the former Wails/Electron distributions and their
historical limits. Those renderer requirements and tests no longer govern the
engine. Use [verification setup](verification/setup.md) for current checks.

## Historical Wails evidence

The [prototype](verification/runs/2026-10-25-wails.md) and
[follow-up](verification/runs/2026-10-25-wails-followup.md) preserve measured
behavior and failures. The [default-package integration](verification/runs/2026-10-25-wails-integration.md)
extends those checks through ordinary `open`, installer, update and rollback,
with exact revisions, artifacts and runner environments. The Ubuntu 22-built
archive also passes unchanged on Ubuntu 24, including explicit runtime setup.

The archives measured about 5.8 MiB (Mac arm64), 6.3 MiB (Linux amd64) and
6.8 MiB (Windows amd64), excluding system runtimes. macOS arm64 and hosted
Windows x64 pass all 21 rendering checks. Ubuntu 22.04
and 24.04 use Xvfb for hosted checks; unavailable WebGL and direct local media
remain explicit exclusions. These are tested environments, not a claim that
every machine with the same OS name works. Physical Windows/Linux Wails input,
clean consumer-machine setup, other architectures, long sleep/wake and public
signing remain separate qualification work.

The [public alpha checks](verification/runs/2026-10-26-release.md) passed all
installer/update paths but observed one intermittent Mac media metadata/seek
assertion failure. Hosted and repeated local follow-ups passed on the same
archive; the cause remains unconfirmed. Media behavior remains provisional.
The first Mac archive's actual deployment target is macOS 26.0.

The [current Wails AFK check](verification/runs/2026-10-27-afk-contract.md) now
passes one OS-confirmed 355-second Mac sleep/wake cycle: the user's post-wake
edit and original PNG returned through the exact-conversation Codex callback.
The separate two-hour page timer also passed, spanning two recorded sleep/wake
cycles and 2 h 7 min 47 s of wall time. Its unchanged note and PNG returned through
one callback and one collection, without polling. This is bounded Mac evidence,
not an indefinite or cross-platform sleep/wake guarantee.

## Accepted limits

- Skill pages are trusted local application code and can access Wails native APIs.
- Sessions have independent processes and saved data. Browser storage isolation
  is best effort; macOS/Linux can share localStorage. Use Fungi drafts for recovery.
- Linux direct local audio/video can fail through the `wails:` asset scheme.
  The same WAV's metadata loaded through a Blob in the probe; audible playback
  and large streaming files are unverified. No transparent workaround is promised.
- After a crash, collect a committed result or reopen an acknowledged draft.
  Close any surviving blank window first. No automatic restart or heartbeat.

These are [accepted product limits](planning.md#accepted-wails-limits-and-next-milestone-on-2026-10-25).
A green Linux job excludes the known direct-media failure and can record a
missing WebGL context; it never relabels either as a successful render.
Other failed checks still fail the job.

## Skill authoring and failed startup

Skills bring their HTML, CSS, JS, libraries and assets; they do not need a
platform-specific plugin. Use the same four-method bridge. Feature-detect
WebGL, media and engine-specific APIs, with a suitable alternative or a clear
requirement when unavailable. See [authoring](authoring.md).

A missing helper reports `renderer_unavailable`; a missing engine reports
`runtime_missing`. A helper that starts but exits before readiness produces
one saved `renderer_startup_failed` result and points to its `worker.log`.
Collect returns the same saved failure. Inspect the actual error; do not
bypass host security policy to manufacture a passing result.

The [desktop recipe](verification/desktop.md) covers visible input, draft
recovery and returned PNGs. A successful build or headless fixture alone does
not prove physical desktop interaction.

## Historical Electron evidence

Earlier development packages bundled Electron. Their observations remain in
[hosted tests](verification/runs/2026-10-22-hosted.md),
[baseline tests](verification/runs/2026-10-22-baselines.md),
[physical Windows tests](verification/runs/windows-2026-10-22-203907.md) and
[distribution tests](verification/runs/2026-10-23-distribution.md).
Ubuntu 24.04's former Chromium sandbox-startup failure belongs to that
renderer; it is not the current Wails support policy. The older no-install
Linux decision was superseded by explicit opt-in runtime setup.

The current updater accepts format-2 Wails packages only. Electron format-1
installation/rollback support was removed during the cleanup. Old portable
receipts are historical; any still-needed portable session requires its original
files. No old session paths are silently migrated.

---
title: Verified environments and portable distribution
updated: 2026-10-26
---
# Verified environments and portable distribution

Current packages contain the Go CLI and a Wails v2.17.0 window helper, with
one archive per OS/CPU. They contain no Electron, Node or browser engine.
The system web engine is a runtime prerequisite:

| Platform | Engine | Missing-runtime setup |
| --- | --- | --- |
| macOS 26+ (first alpha, arm64) | WKWebView supplied by macOS | OS repair/update; no separate runtime installer |
| Windows | Microsoft WebView2 | Explicit bootstrap `-SetupRuntime`, official signed Microsoft installer |
| Ubuntu 22.04 / 24.04 | WebKitGTK 4.1 and its runtime dependencies | Explicit bootstrap `--setup-runtime`, official Ubuntu apt packages |

`fungi doctor` checks the package/runtime without opening a window. Ordinary
launches never install software. Other Linux distributions need compatible
libraries supplied by their operator; the bootstrap does not guess package
managers. A graphical desktop session is still required. Runtime availability
alone does not prove display, GPU, codec or input behavior.

The [installation contract](features/installation.md) covers bootstrap,
updates, rollback and migration from old development previews. The current
integration passed its accepted scope. The first public distribution is an
unsigned alpha; signed/notarized consumer downloads remain separate work.

## Current Wails evidence

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

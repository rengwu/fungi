---
title: Wails simplicity, Linux media and manual recovery
updated: 2026-10-25
result: Failed
---
# Wails simplicity, Linux media and manual recovery

**Failed for direct Linux media; manual recovery passed locally.** The user's
[clarification](../../planning.md#wails-simplicity-clarification-on-2026-10-25)
accepts native access and best-effort browser storage isolation. These are no
longer replacement blockers. This follow-up does not revise the observations
or historical verdict in the [initial assessment](2026-10-25-wails.md).

## Scope and build

[FEAT-interfaces](../../features/interfaces.md) and
[FEAT-sessions](../../features/sessions.md): test the suspected WebKit scheme
allowlist cause on real platforms, retain the existing rendering fixture,
and exercise explicit draft recovery after stopping the entire Wails host.
No attempt was made to read a user's clipboard, change OS policy, install
runtimes on their machine, or reconfigure a harness.

The experiment at `639db19aff51fc2c9579355a1bfe1d9047a3e6b9` adds only a
Linux process-local `WEBKIT_GST_ALLOWED_URI_PROTOCOLS` setting. Wails remains
v2.17.0, with Go 1.26.4 and Node 22.19.0 for development drivers. The
[build manifest](../../evidence/2026-10-25-wails-followup/build.json) records
local source/executable hashes. The
[hosted run](https://github.com/rengwu/fungi/actions/runs/40191649597) and
[metadata](../../evidence/2026-10-25-wails-followup/media-ci.json) identify
all four targets. The implementing agent also verified this work.

## Results

| Scenario | Result | Observation and evidence |
| --- | --- | --- |
| Local macOS renderer fixture | Passed | [All 21 assertions](../../evidence/2026-10-25-wails-followup/mac-renderer.json), including media seeking, cancel/reopen, restored PNG and independent artifact hash |
| Hosted macOS | Passed | [All 21 assertions](../../evidence/2026-10-25-wails-followup/wails-macos-latest-renderer.json) |
| Hosted Windows | Passed | [All 21 assertions](../../evidence/2026-10-25-wails-followup/wails-windows-latest-renderer.json) |
| Ubuntu 22.04 | Failed | [20/21 assertions](../../evidence/2026-10-25-wails-followup/wails-ubuntu-22.04-renderer.json); direct media still fails despite allowing the scheme |
| Ubuntu 24.04 | Failed | [20/21 assertions](../../evidence/2026-10-25-wails-followup/wails-ubuntu-24.04-renderer.json); same failure |
| Blob comparison on Linux | Passed | [Ubuntu 22](../../evidence/2026-10-25-wails-followup/wails-ubuntu-22.04-host.json) and [Ubuntu 24](../../evidence/2026-10-25-wails-followup/wails-ubuntu-24.04-host.json) load metadata for the same WAV via Blob; direct URL returns media error 4. This checks metadata, not audible playback or arbitrary codecs |
| Linux loader diagnosis | Passed | At `3b62acb188e01a81b3ef518d2724f83d2bc1c33b`, the unchanged renderer with added GStreamer logging reports `No URI handler implemented for "wails"` on [Ubuntu 22](../../evidence/2026-10-25-wails-followup/ubuntu-22.04-media.log) and [Ubuntu 24](../../evidence/2026-10-25-wails-followup/ubuntu-24.04-media.log). This diagnoses the failure; it does not repair it |
| Full Wails host exit and manual recovery on Mac | Passed | [Disposable session evidence](../../evidence/2026-10-25-wails-followup/mac-recovery.json): save a distinctive draft, SIGKILL only the exact helper process, collect once to receive `failed/helper_exited`, explicitly reopen and submit the restored draft; original failure stays unchanged |
| Existing default platform checks | Passed | [Foundation CI metadata](../../evidence/2026-10-25-wails-followup/foundation-ci.json): existing four-platform core/Electron/install/update scopes still pass |

The recovery check initially looked for `reopenedFrom` on the result envelope;
it belongs to `launch.json`. Correcting that test assumption produced the
successful run above without a product change. The test script and synthetic
sessions remain under `verification-output/`.

## Interpretation and limits

The allowlist experiment did not fix media and was removed. The renderer
source is restored byte-for-byte to the earlier prototype. Upstream WebKitGTK
2.52.6's web media source registers HTTP, HTTPS and Blob protocols; permitting
an additional scheme does not supply its missing loader. See
[primary sources and analysis](../../research/system-renderers.md#linux-media-follow-up).
No framework fork, HTML rewriting, transparent Blob conversion or local media
server was added.

The [diagnostic run](https://github.com/rengwu/fungi/actions/runs/40192079604)
and [metadata](../../evidence/2026-10-25-wails-followup/diagnostic-ci.json)
retain the same 21/21 Mac/Windows and 20/21 Linux outcomes. Detailed
[Ubuntu 22](../../evidence/2026-10-25-wails-followup/ubuntu-22.04-diagnostic.json)
and [Ubuntu 24](../../evidence/2026-10-25-wails-followup/ubuntu-24.04-diagnostic.json)
results still load Blob metadata successfully. The host probe now enables
targeted GStreamer logs only in its disposable Linux test process.

Best-effort browser storage and native APIs are accepted requirements, not
failures to hide. Fungi still owns separate session directories and durable
records. Full host exit is distinct from a web-content-only crash: the latter
can leave a blank native window and live owner. Close that window before
explicit reopen. Web-content-only crash detection across platforms, audible
Linux playback, long idle/sleep under Wails, clean-machine runtime setup,
release packaging and updater integration remain unqualified by this run.

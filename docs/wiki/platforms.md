---
title: Verified environments and portable distribution
updated: 2026-10-22
---
# Verified environments and portable distribution

Fungi ships as a portable directory containing the Go executable and pinned
Electron runtime. Copy or extract the whole directory; no installer, Node,
Go, npm or separate browser is required on a verified environment. It does
not install missing OS libraries or change sandbox policy.

The [Linux distribution decision](planning.md#linux-distribution-decision)
keeps this no-install model and limits support to explicitly verified
environments. These are development-preview baselines, not a claim that every
machine with the same OS name works. OS libraries, security policy, desktop
session and graphics capabilities still matter. Signing/notarization and
clean-machine distribution have not passed acceptance.

## Current evidence

| Environment actually exercised | Verified scope | Limits |
| --- | --- | --- |
| Development Mac, macOS 27.0 arm64 (26A5368g) | Native input, paste, saved drafts, recovery, image artifacts, exact-conversation return; portable renderer with empty PATH | One short sleep/wake cycle; two-hour idle proof uses a native helper, not the Electron window |
| Hosted macOS 26.6.2 arm64 | Core, actual Electron workflow and all 21 renderer checks, including WebGL | Hosted runner; not a clean consumer machine |
| Hosted Windows Server 2025 x64, 10.0.26100 | Core, actual Electron workflow and all 21 renderer checks, including WebGL | Physical input and an extracted archive on consumer Windows remain unverified |
| Hosted Ubuntu 22.04.5 x64 with Xvfb | Core, actual Electron cancel/reopen/submit workflow, artifacts and 20 renderer checks | WebGL 2 context unavailable; no physical desktop/input or clean-machine claim |
| Hosted Ubuntu 24.04.5 x64 with Xvfb | Core only; portable renderer fails before readiness under tested sandbox policy | **Unsupported for UI in this environment** |
| Linux arm64 under OrbStack | Core tests in an isolated unprivileged container | Desktop renderer untested |

The [original hosted report](verification/runs/2026-10-22-hosted.md) preserves
the failed broad rendering matrix and exact runner images. Selecting a
narrower baseline does not change those observations. Hosted runners update;
each later report identifies the image actually tested.
The [follow-up baseline run](verification/runs/2026-10-22-baselines.md) passes
the selected scopes below, including the unsupported-startup diagnostic test.

## Skill authoring

Ordinary HTML/CSS, forms, modules, Canvas 2D, SVG, workers, WebAssembly and
local media use the same four-method bridge. Skills bring their libraries and
assets. Feature-detect WebGL and codecs; provide an alternative or explain
when an interface requires an unavailable feature. Fungi does not override
Chromium's GPU blocklist or switch off its sandbox to enable an interface.

No platform plugin or installation detector is needed in a skill. See
[authoring](authoring.md) for the mechanical launch and response procedure.

## Failed startup

If the renderer starts but exits before readiness, the CLI returns nonzero
with one saved `failed` result and `error.code: renderer_startup_failed`.
Its message locates the session's `worker.log`; inspect that file for the
actual runtime or OS error. Collection returns the same saved failure.
The code identifies a startup failure, not its cause: a missing display,
library or sandbox permission must be diagnosed from the log.

On the tested Ubuntu 24.04 runner, Chromium requests a root-owned sandbox
helper after unprivileged namespace creation is denied. That environment is
outside the portable UI baseline. Record the failure; use a verified desktop
environment. An installer, privileged helper setup or changed host policy is
not part of the accepted distribution. Do not automatically retry with
`--no-sandbox`.

## What a green check means

CI names its scope explicitly. macOS and Windows require WebGL and all
workflow checks. Ubuntu 22.04 requires the same workflow but records a null
WebGL context as **unavailable**, never as a successful rendering check. If a
context exists but renders incorrectly, the check fails. Ubuntu 24.04 tests
the core and the saved unsupported-startup error with sandbox diagnostics;
its green job does not claim UI support. Unexpected successful startup there
requires review and a full rendering run before widening the baseline.

The [desktop recipe](verification/desktop.md) covers the remaining physical
archive handoff. A successful cross-build alone does not add a supported
environment.

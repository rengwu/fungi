---
title: Verified environments and portable distribution
updated: 2026-10-23
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

For the eventual release, command installation and updates are now requirements;
see the [proposed delivery flow](#planned-command-installation-and-updates).
These convenience commands are not implemented yet.

## Current evidence

| Environment actually exercised | Verified scope | Limits |
| --- | --- | --- |
| Development Mac, macOS 27.0 arm64 (26A5368g) | Native input, paste, saved drafts, recovery, image artifacts, exact-conversation return; portable renderer with empty PATH | One short sleep/wake cycle; two-hour idle proof uses a native helper, not the Electron window |
| Hosted macOS 26.6.2 arm64 | Core, actual Electron workflow and all 21 renderer checks, including WebGL | Hosted runner; not a clean consumer machine |
| Hosted Windows Server 2025 x64, 10.0.26100 | Core, actual Electron workflow and all 21 renderer checks, including WebGL | Hosted runner; consumer desktop coverage is recorded separately below; clean-machine/archive acceptance remains unverified |
| Physical Windows 11 Pro x64, 10.0.26200 | Source-built portable package with empty child PATH; agent-operated native input, acknowledged draft, close/reopen, comparison PNG/hash and cancellation | 17 core and six Node tests passed; symlink-dependent tests and all 21 renderer fixture checks unexecuted for lack of fixture privilege. No clean-machine/archive, clipboard or local WebGL claim |
| Hosted Ubuntu 22.04.5 x64 with Xvfb | Core, actual Electron cancel/reopen/submit workflow, artifacts and 20 renderer checks | WebGL 2 context unavailable; no physical desktop/input or clean-machine claim |
| Hosted Ubuntu 24.04.5 x64 with Xvfb | Core only; portable renderer fails before readiness under tested sandbox policy | **Unsupported for UI in this environment** |
| Linux arm64 under OrbStack | Core tests in an isolated unprivileged container | Desktop renderer untested |

The [original hosted report](verification/runs/2026-10-22-hosted.md) preserves
the failed broad rendering matrix and exact runner images. Selecting a
narrower baseline does not change those observations. Hosted runners update;
each later report identifies the image actually tested.
The [follow-up baseline run](verification/runs/2026-10-22-baselines.md) passes
the selected scopes below, including the unsupported-startup diagnostic test.
The [Windows desktop handoff](verification/runs/windows-2026-10-22-203907.md),
completed on 2026-10-23, passes its visible desktop journey. Its aggregate is
**Untested** because required symlink-dependent scenarios could not execute on
that account. This adds the observed Windows 11 workflow without extending
the hosted renderer capability results to this machine.

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
environment. An OS dependency installer, privileged helper setup or changed
host policy is not part of the accepted distribution. Do not automatically retry with
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

## Planned command installation and updates

The [2026-10-23 requirement](planning.md#command-installation-and-updates-on-2026-10-23)
is a simple command to install and update Fungi on each supported platform,
without requiring a language toolchain or package manager. The recommended
first release uses complete prebuilt archives and a small bootstrap:

| Platform | Proposed bootstrap | Subsequent update |
| --- | --- | --- |
| macOS | Shell script using curl or wget and OS archive tools | `fungi update` |
| Supported Linux environments | Same shell bootstrap, selecting the matching OS/CPU archive | `fungi update` |
| Windows | PowerShell script using built-in download, ZIP extraction and hashing | `fungi update` |

These are proposed commands, not current CLI capabilities. The bootstrap should
select one OS/CPU package, download and verify it, place it in a user-owned
location without administrator access, and make the stable `fungi` command
available. It should report any required shell restart or PATH adjustment.
Manual archive download remains available. No Go, Node, npm, Git, Homebrew or
separate browser should be needed to run the installed package.

Keep release selection and update logic in Go rather than maintaining separate
updaters per platform. Updates should be explicit and repeatable, with useful
exit codes and machine-readable results for agents. Start with full-package
downloads; delta updates and a background updater are unnecessary for this
target. Optional Homebrew/npm distribution can use the same artifacts later.
Package-manager-owned installations should update through their manager rather
than be overwritten by Fungi's updater.

Stage and verify the complete next version separately, then activate it through
a stable command path. A failed download, verification or activation must leave
the last working installation usable. Support selecting a specific release and
returning to the previous version. Keep session state outside replaceable
installation files. Running sessions must finish normally during an update,
including on Windows where running binaries may prevent replacement.

Current receipts/callbacks contain absolute executable paths, and reopening can
refer to assets in an older package. Define how those references survive updates
before enabling cleanup. Retain the current release, a rollback release and any
older versions still needed by sessions; prune only unreferenced versions.
Avoid treating terminal sessions as automatically safe to discard: their saved
receipts and drafts can still be collected or reopened later.

Acceptance must cover clean-machine installation, repeated install/update,
interrupted and invalid downloads, concurrent update attempts, activation
failure, rollback, and updates while a session is open or awaiting collection.
Verify old callback commands, draft reopening and artifact access after update
and cleanup on each supported platform. Signing/notarization and authentic
release publication remain part of release readiness. No such result is claimed
by the current portable-package tests.

Platform references: Windows PowerShell provides
[HTTP downloads](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/invoke-webrequest?view=powershell-5.1),
[ZIP extraction](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.archive/expand-archive?view=powershell-5.1)
and [SHA-256 hashing](https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.utility/get-filehash?view=powershell-5.1).
Use `-UseBasicParsing` for downloads on PowerShell 5.1. An optional npm channel
adds the [Node/npm prerequisite](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm/),
so it should not be the only installation path.

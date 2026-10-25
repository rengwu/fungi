---
title: Wails default packages, runtime setup and updates
updated: 2026-10-25
result: Passed
---
# Wails default packages, runtime setup and updates

**Passed for the accepted integration scope**, including all four native CI
environments and one shared Linux archive across Ubuntu 22/24. This report covers the
[accepted Wails scope](../../planning.md#accepted-wails-limits-and-next-milestone-on-2026-10-25),
not full Electron parity or public release readiness. Earlier direct Linux
media failures remain failures in their original reports and an explicit
exclusion here.

## Build and scope

Affected features: [interfaces](../../features/interfaces.md),
[sessions](../../features/sessions.md) and
[installation](../../features/installation.md). The implementation at
`737d8595676caa6b1d87458236a25ce8a9e37642` makes normal `fungi open` select the
Wails helper. It adds format-2 packages, runtime preflight, `doctor` and
explicit bootstrap runtime setup; removes Electron rendering and runtime
sharing; and retains format-1 validation for preview rollback.

The [first complete CI run](https://github.com/rengwu/fungi/actions/runs/40200305216)
and [saved metadata](../../evidence/2026-10-25-wails-integration/ci.json) identify
all jobs. Wails is v2.17.0; development drivers use Go 1.26.4 and Node 22.19.0.
The [final five-job run](https://github.com/rengwu/fungi/actions/runs/40201171453)
at `4d9321583e990cda8cc6f4555dbce8670cf152e3` repeats all native checks and adds
Ubuntu 22-built archives on Ubuntu 24; [metadata](../../evidence/2026-10-25-wails-integration/compatibility-ci.json)
records five successful jobs.

All actual UI/install checks use the packaged helper through the ordinary CLI,
not the prototype `--exec` override. Installed CLI children have an empty PATH.
The implementing agent also performed the verification; no independent review
is claimed.

| Environment | Recorded image/build | Scope |
| --- | --- | --- |
| Local macOS arm64 | 27.0, build 26A5368g | Packaged automated fixture, real install/update/rollback, native comparison/form input |
| Hosted macOS arm64 | 26.6.2, image 20261021.0445.8 | All 21 rendering checks and actual bootstrap/update/rollback |
| Hosted Windows x64 | Server 2025 10.0.26100, image 20261021.278.1 | All 21 rendering checks and PowerShell bootstrap/update/rollback with existing WebView2 |
| Hosted Ubuntu x64 | 22.04.5, image 20261012.322.1 | Rendering with accepted media exclusion, active-UI update/rollback, missing-runtime container test |
| Hosted Ubuntu x64 | 24.04.5, image 20261012.334.1 | Same scope, initially with its own native build |

[Local build hashes](../../evidence/2026-10-25-wails-integration/local-build.json),
[Mac manifest](../../evidence/2026-10-25-wails-integration/macos-latest-build.json),
[Windows manifest](../../evidence/2026-10-25-wails-integration/windows-latest-build.json),
[Ubuntu 22 manifest](../../evidence/2026-10-25-wails-integration/ubuntu-22.04-build.json)
and [Ubuntu 24 manifest](../../evidence/2026-10-25-wails-integration/ubuntu-24.04-build.json)
record the executable identities. Environment log excerpts are retained beside
them. The native Mac journey uses `v0.1.0-wails.4`: identical core/helper source
to the first CI revision, with a packaging-only Homebrew license-path repair.

## Results

| Selected scenario | Result | Evidence |
| --- | --- | --- |
| Core validation, session/delivery behavior, update failure/rollback and archive tests | Passed | Core vet and tests in all four CI jobs; race detector on Unix. Runtime-unavailable fixtures verify saved launch failure and unchanged activation; legacy format-1 packages remain readable |
| Helper assets and JS bridge | Passed | Nested Wails helper vet/tests and six Node tests in CI; ranges, HEAD, MIME, root confinement, ordered Blob saves and submit after a failed save |
| Mac ordinary-open rendering | Passed | [Local](../../evidence/2026-10-25-wails-integration/local-renderer.json) and [hosted](../../evidence/2026-10-25-wails-integration/macos-latest-renderer-smoke.json): all 21 checks, cancel/reopen, PNG/hash |
| Windows ordinary-open rendering | Passed | [All 21 checks](../../evidence/2026-10-25-wails-integration/windows-latest-renderer-smoke.json), cancel/reopen, PNG/hash |
| Ubuntu ordinary-open rendering within accepted scope | Passed | [22.04](../../evidence/2026-10-25-wails-integration/ubuntu-22.04-renderer-smoke.json) and [24.04](../../evidence/2026-10-25-wails-integration/ubuntu-24.04-renderer-smoke.json): 20 passing checks; direct-media failure remains in `excluded` |
| Real bootstrap, active session during update, stable collection, artifact and rollback | Passed | [Local Mac](../../evidence/2026-10-25-wails-integration/local-install.json), [hosted Mac](../../evidence/2026-10-25-wails-integration/macos-latest-install-smoke.json), [Windows](../../evidence/2026-10-25-wails-integration/windows-latest-install-smoke.json), [Ubuntu 22](../../evidence/2026-10-25-wails-integration/ubuntu-22.04-install-smoke.json), [Ubuntu 24](../../evidence/2026-10-25-wails-integration/ubuntu-24.04-install-smoke.json). Download/checksum requests, immutable versions, live PNG and rollback are recorded |
| Missing Linux runtime and explicit setup | Passed | Real official [Ubuntu 22 container](../../evidence/2026-10-25-wails-integration/ubuntu-22.04-runtime-setup.json) and [Ubuntu 24 container](../../evidence/2026-10-25-wails-integration/ubuntu-24.04-runtime-setup.json): doctor exits 3; ordinary open saves `runtime_missing`; bootstrap without opt-in refuses; injected apt failure leaves no selected install; real opt-in apt setup succeeds and repeats without reinstalling |
| Installed Mac native comparison/recovery | Passed | [Native observations](../../evidence/2026-10-25-wails-integration/native-observation.json): native mouse/keyboard input, acknowledged draft, close, one collect, explicit reopen with the same 65% divider/After/note, submit, one collect; [result](../../evidence/2026-10-25-wails-integration/native-submitted.json) and visually inspected [PNG](../../evidence/2026-10-25-wails-integration/native-comparison.png) match the independently calculated hash |
| Installed primitive form and native file picker | Passed | Enter synthetic text, choose the comparison PNG through the native file dialog, observe preview and Draft saved, submit and collect once; [result](../../evidence/2026-10-25-wails-integration/native-form-submitted.json) preserves exactly the selected bytes/hash |
| Same Ubuntu 22 archive on Ubuntu 24 | Passed | [Archive identity](../../evidence/2026-10-25-wails-integration/compatibility-archive.json) matches the installer checksum. [Rendering](../../evidence/2026-10-25-wails-integration/compatibility-renderer-smoke.json), [live update/rollback](../../evidence/2026-10-25-wails-integration/compatibility-install-smoke.json) and [missing-runtime setup](../../evidence/2026-10-25-wails-integration/compatibility-runtime-setup.json) all pass within the same accepted media scope |

## Size and packaging

[Measured version-1 archives](../../evidence/2026-10-25-wails-integration/package-sizes.json):
Mac arm64 6,039,106 bytes (5.76 MiB), Linux amd64 built on Ubuntu 22
6,626,882 bytes (6.32 MiB), Windows amd64 7,101,027 bytes (6.77 MiB).
Each contains the CLI, native helper, small examples and dependency notices.
It contains no Chromium/Electron, full wiki/evidence or development toolchain.
The system runtime's disk/network size is separate from these package sizes.
A machine obtains one matching archive, not the development `dist/` directory.

The local package build initially failed because Homebrew places Go's LICENSE
beside `libexec`, outside GOROOT. The packaging fallback repaired that layout;
a fresh `v0.1.0-wails.4` build, installation and native input succeeded. The
partial `.3` directory was not used. No old session-referenced builds were
deleted. A documentation-only intermediate CI run was cancelled after fixing
an obsolete anchor; the replacement run identifies its tested revision.

## Accepted exclusions and remaining release work

Direct Linux local audio/video is still unreliable through the `wails:` scheme.
The fixture's failure is retained as excluded, not counted as a passing media
check. WebGL is feature-detected; Linux may explicitly record it unavailable.
Native Wails access, best-effort browser storage separation and explicit manual
crash recovery are accepted product limits, with no new framework fork or
recovery daemon.

Missing WebView2 installation on a Windows machine that actually lacks it is
**Untested**; Windows CI reused an existing runtime. The implemented path uses
the official Microsoft bootstrapper, Authenticode verification and an explicit
installer flag. No user machine was modified to manufacture absence. Linux's
missing-runtime setup proof is container-level and does not establish audible
playback, physical desktop input or offline clean-machine installation.

Other untested release work includes signed/notarized downloads, actual public
release assets, fresh interactive PATH behavior, physical Windows/Linux Wails
input, other OS versions/architectures and long Wails sleep/wake. Actual
automatic harness callbacks retain the older build-specific evidence; this
integration tests the unchanged delivery logic and manual native collection,
not a new long-lived human callback. No public release was published, no user
OS policy changed, and no harness configuration installed.

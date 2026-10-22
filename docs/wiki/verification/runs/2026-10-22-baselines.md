---
title: Verified platform baselines and startup diagnostics
updated: 2026-10-22
result: Passed
---
# Verified platform baselines and startup diagnostics

**Passed for the selected scope:** the user's accepted no-install platform
boundary, a saved and diagnosable renderer startup failure, and regression
checks of the supported workflows. This covers [sessions](../../features/sessions.md)
and [interfaces](../../features/interfaces.md). The earlier
[broad hosted run](2026-10-22-hosted.md) remains **Failed**; Ubuntu 24.04 has
not gained UI support and Ubuntu 22.04's virtual display still lacks WebGL.

## Build and scope

Product revision: `199759a5f0341d8b97b8683385a1d4d06080b8ef`, Electron
44.9.0, Go 1.26.4. [The CI run](https://github.com/rengwu/fungi/actions/runs/39758430600)
and [saved job record](../../evidence/2026-10-22-baselines/ci.json) identify
the exact source. [Observed results](../../evidence/2026-10-22-baselines/observations.json)
include the actual runner OS versions; logs retain image versions and checks.

The [accepted answer](../../journal/2026-10-22-grill-0c10be35-a07a-4a8d-8d5d-39fda5b2e1bc.md)
keeps portable distribution on explicitly verified environments. The
[platform record](../../platforms.md) makes those limits discoverable. It
does not prescribe an installer, elevated helper, changed host policy or
disabled sandbox. No four-method bridge or native-helper protocol changed.

The code now identifies a bundled renderer that exits before readiness as
`renderer_startup_failed`, retaining the failure in the normal immutable
result and pointing to its actual `worker.log`. It does not guess the cause
from an exit signal or parse arbitrary log text into remediation commands.

## Exercised scenarios

| Scenario | Result | Observation |
| --- | --- | --- |
| Local core and Node regression | Passed | 18 root Go tests with race detection, Go vet and seven Node tests; the new fixture starts the real CLI with a deliberately failing renderer executable, reads its stderr log and collects the identical saved failure |
| Existing helper exit behavior | Passed | Native helper retains `helper_exited`; an already-written valid submission survives a later nonzero helper exit |
| Hosted core on all four targets | Passed | Each target passes Go vet, 18 root tests, seven Node tests and wiki validation; race detector on Unix, ordinary tests on Windows |
| Hosted macOS 26.6.2 arm64 | Passed | [All 21 renderer checks](../../evidence/2026-10-22-baselines/macos-renderer.json), cancel/reopen/submit, restored Blob and independently checked PNG SHA-256 |
| Hosted Windows Server 2025 x64, 10.0.26100 | Passed | [All 21 renderer checks](../../evidence/2026-10-22-baselines/windows-renderer.json) and the same complete workflow |
| Hosted Ubuntu 22.04.5 x64 with Xvfb | Passed | [20 required checks](../../evidence/2026-10-22-baselines/ubuntu22-renderer.json) and the complete workflow; WebGL context explicitly unavailable and excluded from this selected baseline |
| Unsupported Ubuntu 24.04.5 startup error | Passed | [Nonzero CLI result](../../evidence/2026-10-22-baselines/ubuntu24-renderer.json) is saved as `renderer_startup_failed`; its indicated log contains Chromium's actual sandbox abort; no UI readiness or submission is claimed |
| Fresh portable macOS build, macOS 27.0 arm64 | Passed | [All 21 checks](../../evidence/2026-10-22-baselines/mac-package-renderer.json), full workflow and PNG hash with an empty child PATH |

The hosted diagnostic logs are [Windows](../../evidence/2026-10-22-baselines/windows.log),
[macOS](../../evidence/2026-10-22-baselines/macos.log),
[Ubuntu 22.04](../../evidence/2026-10-22-baselines/ubuntu22.log) and
[Ubuntu 24.04](../../evidence/2026-10-22-baselines/ubuntu24.log).

The renderer driver retains each actual capability outcome. Its default
requires WebGL. Only the explicit Ubuntu 22.04 scope permits a null context;
an available context with incorrect pixels still fails. Ubuntu 24.04's
separate negative test requires the saved startup error and sandbox log. If
that renderer starts successfully in a future image, the check fails with an
instruction to review the platform matrix and run the complete fixture.

All actual renderer launches use an empty child executable search path and
no inherited Codex thread ID. Linux uses the runner's existing Xvfb. No
machine configuration, renderer sandbox or GPU blocklist was bypassed. These
are disposable fixtures, not agent answers in a human questionnaire.

## Portable handoff and remaining limits

[Build metadata and source hashes](../../evidence/2026-10-22-baselines/packages.json)
identify the fresh macOS arm64, Windows x64 and Linux x64 portable directories
created at `2026-10-22T09-50-20-*`. The Mac package ran locally. Hosted Windows
and Linux used matching source builds; their cross-packaged archives still
need physical target-machine acceptance using the bundled `VERIFY.md`.
The packager's `runtimeVerified: false` means packaging alone does not certify
a target; the runtime observations are recorded separately here.

Consumer Windows input/clipboard, clean-machine security prompts, public
signing/notarization, arbitrary Linux desktops, every GPU/codec/library and
power-loss durability remain unverified. Existing harness and AFK evidence
remains version- and route-specific; this run does not expand it. The
implementing agent performed verification; no independent review is claimed.
Full Fungi-native jstack authoring integration remains deferred.

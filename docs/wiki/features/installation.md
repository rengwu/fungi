---
title: Installation and updates
updated: 2026-10-28
---
# Installation and updates

**ID: FEAT-installation**

The accepted distribution contains one Go engine executable with docs and
license metadata. Other projects distribute all UI/helper implementations.

## Entry points

Build with Go 1.26; package with Node 22+ using `scripts/package.mjs`. The packager
supports explicit OS/CPU cross-builds with CGO disabled. Release targets remain
Mac arm64, Windows x64 and Linux x64 until further execution qualification.
Run an extracted executable or create a managed installation with `install`.
Shell/PowerShell bootstrap verifies a matching archive/checksum; there are no
webview checks, runtime setup flags or system package installations.

## Expected behavior

Engine manifests use format 3 and protocol 2, with OS, architecture, version and
engine SHA-256. Install/update validates these and the candidate executable
before activation. Format-1/2 UI packages and mismatched protocols are rejected.
Managed versions are immutable. The stable launcher selects the active release;
existing owners continue through updates. Rollback selects a retained engine
release. Prune protects current, previous, active and session-referenced
versions. External helpers/assets are maintained by their provider.

State stays outside managed installations. POSIX bootstrap prints a PATH entry;
PowerShell updates it unless NoPath is requested. No background updater exists.
Existing protocol-1 sessions require their original binary/state; use a separate
engine installation during migration. Published 0.1 assets remain unchanged.

## Build and verification recipe

Build two unique versions and run:

```sh
node scripts/package.mjs --release v0.2.0-test.1
node scripts/package.mjs --release v0.2.0-test.2
node scripts/install-smoke.mjs ARCHIVE_A ARCHIVE_B
```

Use the actual archive paths printed by packaging. The driver serves archives
from a local fixture server, bootstraps into a path with spaces, runs installed
commands with an empty PATH, keeps an external helper active through update,
collects its saved artifact, and verifies rollback and help. `go test` covers
failed checksums, unsupported formats, session pins and pruning. All data is
under a unique verification-output directory; no user installation is changed.

## Public release workflow

The engine CI runs on macOS, Windows and Ubuntu without GUI setup. The Ubuntu
22 archive is also exercised unchanged on Ubuntu 24. Publication requires a
reviewed tag and notes. The public-download driver checks the exact published
engine and preserves evidence. See the repository releasing manual.

## Coverage and limits

Cross-compilation is not runtime verification. Hosted workflow execution,
public publication, signing and clean consumer-machine setup require separate
evidence. Previous renderer/runtime-setup reports remain historical.

## Verification history

- [Engine-only migration](../verification/runs/2026-10-28-headless.md) — protocol 2; see run-specific results.


- [2026-10-26 terminal help](../verification/runs/2026-10-26-terminal-help.md) — one-line install examples and alpha.2 public installation/update checks.

- [First public alpha](../verification/runs/2026-10-26-release.md) — publication and public-download checks; all public install/update paths passed; an intermittent Mac media probe failure remains recorded.

- [2026-10-26 cleanup](../verification/runs/2026-10-26-cleanup.md) — removed Electron preview compatibility; current Wails installation/update/rollback checks.

- [2026-10-26 offline discovery](../verification/runs/2026-10-26-discovery.md) — installed help, examples, update/rollback and native throwaway reply.

- [2026-10-25 Wails integration](../verification/runs/2026-10-25-wails-integration.md) — default packages, installed native input, update/rollback and explicit Linux runtime setup.

- [2026-10-23 distribution proof](../verification/runs/2026-10-23-distribution.md)

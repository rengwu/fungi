---
title: Verification setup
updated: 2026-10-28
---
# Verification setup

Current verification targets the headless protocol-2 engine. Other projects
own UI rendering and its tests. Historical UI recipes/results remain attached
to their dated reports; they are not current engine prerequisites.

## Build and launch

Use Go 1.26 and Node 22+ for development. No GUI toolchain is required.

```sh
go vet ./...
go test -race ./...
node scripts/package.mjs --release v0.2.0-test.1
node scripts/package.mjs --release v0.2.0-test.2
node scripts/engine-smoke.mjs /absolute/package/fungi
node scripts/install-smoke.mjs ARCHIVE_A ARCHIVE_B
```

Use fresh version identifiers and the paths printed by packaging. On Windows
omit -race and use fungi.exe. The engine itself has no Node/Go runtime dependency.
The Node fixture is an external helper used only by verification.

## Engine and helper checks

`engine-smoke.mjs` copies the supplied executable into an isolated directory,
clears PATH and display variables, and exercises doctor, offline help, concurrent
helpers, JSON/artifact fidelity, manual close, explicit draft/artifact reopen
and rejected UI arguments. It invokes the external helper by an absolute Node
executable and script path. `helper-fixture.mjs` has no UI or model calls.

Tests use bounded file gates solely to synchronize automated helpers. Never
reuse these waits for a person's session. `go test` covers owner interruption,
malformed helpers, startup failures, blocked writes, immutable results, limits,
legacy-protocol rejection and uncertain notification behavior.

## Installation and updates

`install-smoke.mjs` serves real archives/checksums over a local fixture server,
runs the platform bootstrap in a fresh path with spaces, keeps an external
helper active through update, verifies its saved artifact and rolls back.
Installed commands run with an empty PATH. The Windows driver uses NoPath;
POSIX never edits shell configuration. Only disposable installations change.

The engine package must have format 3, protocol 2, and no renderer directory.
CI runs on macOS, Windows, Ubuntu 22 and Ubuntu 24 without Xvfb or webview setup.
Linux compatibility also executes the older-baseline archive on Ubuntu 24.
Cross-built binaries still need execution evidence before platform qualification.

## Native helper idle probe

```sh
go build -o verification-output/idle-helper ./scripts/idle-helper
```

Pass an input file with seconds and a unique marker through `open --exec`.
The helper saves a draft and submits timing data/artifact after one local timer.
Use wait only for a suitable host pending call. For an exact-conversation
callback, launch once with auto, preserve the receipt, end the turn and collect
once on notification. Never poll. A finite test does not establish indefinite
waiting or arbitrary host sleep/wake behavior.

## Observation and cleanup

Each driver saves evidence under its own verification-output directory. Inspect
result JSON, outcome, artifact bytes/hash and original/reopened IDs. Preserve
BUILD.json, engine hash, source revision and dirty-tree fingerprint alongside
the report. Remove only test-owned resources after saving evidence. The drivers
release outstanding synthetic helpers on failure; they never use real interviews.

## Offline discovery verification

`go test -run 'TestOfflineHelpAndDiscovery|TestTerminalHelp|TestHeadlessDoctor' .`
checks the standalone reference and storage probe. The installer driver repeats
help/discovery before update, after update and after rollback.

## Public alpha downloads

`node scripts/public-release-smoke.mjs vX.Y.Z-alpha.N` checks an explicitly
published protocol-2 engine release, anonymously installs it in a disposable
path, runs doctor/update and the helper smoke, and records evidence. It does
not publish, change persistent PATH or install helper runtimes. Use historical
checkouts for protocol-1 public downloads.

## Wails diagnostics

Historical Wails diagnostics belong to the protocol-1 reports, including the
[prototype record](runs/2026-10-25-wails.md). Those drivers were removed with the
bundled UI. All new interface-specific verification belongs to its UI project.

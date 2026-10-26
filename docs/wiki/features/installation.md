---
title: Installation and updates
updated: 2026-10-26
---
# Installation and updates

**ID: FEAT-installation**

Implement the user's [command installation and update requirement](../planning.md#command-installation-and-updates-on-2026-10-23)
with one prebuilt package per OS/CPU and shared Go update logic. This is a
public alpha implementation; signing and clean consumer-machine acceptance
remain separate work. Skills use the same four page methods.

## Entry points

`fungi --version` reports JSON build identity. `fungi describe installation`
describes these commands:

| Command | Behavior |
| --- | --- |
| `fungi install --archive FILE --sha256 HASH [--dir DIRECTORY]` | Install a verified local package; default directory is `~/.fungi` |
| `fungi update` | Fetch the matching archive/checksum from the latest stable GitHub release |
| `fungi update --version vX.Y.Z` | Select a specific published release |
| `fungi update --archive FILE --sha256 HASH` | Apply an explicitly supplied local release |
| `fungi rollback` | Switch to the previous retained release without downloading |
| `fungi prune` | Remove obsolete versions not required for rollback or sessions |

Updates, rollback and pruning operate on managed installations. Unpacked
portable folders remain usable directly; they are not silently adopted.
Successful mutations emit one JSON `installation` result with the stable
command, current/previous versions, removed versions, cleanup warnings and
`reusedRuntimeBytes` (retained as zero for old clients; Wails packages do not
share a bundled runtime).
Errors use the existing JSON error envelope and nonzero exit status.

The repository's `scripts/install.sh` and `scripts/install.ps1` download and
verify one archive, extract it privately and invoke the Go installer. POSIX
requires a shell, curl or wget, tar and shasum or sha256sum. Windows uses
PowerShell's download, ZIP and hashing commands. Neither requires Go, Node,
npm or Homebrew. Windows adds the user PATH entry unless `-NoPath` is given;
POSIX prints the directory to add, without editing shell startup files.

`FUNGI_INSTALL_DIR` selects a bootstrap destination. The bootstrap-only
`FUNGI_RELEASE_BASE_URL` permits an explicit mirror/test server; it does not
change the installed updater's GitHub source. Public alpha assets and installers
live at `https://github.com/rengwu/fungi/releases`. Alpha install/update commands
select an explicit version; ordinary `fungi update` means latest stable, which
is unavailable until a stable release exists. The first release targets macOS 26+
arm64, Windows amd64 and Ubuntu 22.04/24.04 amd64. Unsupported bootstrap targets
fail with an explanation before downloading.

## Expected behavior

The installation contains a small stable `bin/fungi` entry point and immutable
`versions/vX.Y.Z` directories. The entry point delegates through a committed
`installation.json` selector; it is not replaced while Windows is executing
it. The installation layout is versioned separately from session protocol 1.

Stage a complete package, verify its archive and executable hashes, validate
OS/CPU/runtime metadata, and execute its version probe before switching the
selector. Checksums come from the same HTTPS release source; they are not a
separate cryptographic publisher signature. An invalid download/package or
failed preactivation probe leaves the selected installation unchanged.
Interrupted staging is disposable and an interrupted first install can retry.
Version names cannot be reused with different archive bytes.

New format-2 packages contain the core and Wails helper, with hashes for both.
The updater probes the helper's system runtime before activation. No browser
engine is bundled or shared between versions. The earlier Electron file-sharing
optimization is removed.

`fungi doctor` checks runtime availability without a window: exit 0 means the
runtime is available, 3 means it is missing, and other nonzero results indicate
an invalid command or a helper-check failure. A desktop session and graphics
capabilities are separate requirements. Ordinary launches, checks and updates
never install software.

The bootstrap checks the downloaded package before installing it. Missing
runtimes stop setup unless the user explicitly passes `--setup-runtime` (shell)
or `-SetupRuntime` (PowerShell). Ubuntu 22.04/24.04 use configured distro
repositories for `libwebkit2gtk-4.1-0` and its runtime dependencies; administrator
permission may be needed. Windows downloads Microsoft's official Evergreen
bootstrapper, verifies its Microsoft Authenticode signature and runs it. Without
elevation this requests per-user installation. Existing runtimes are reused;
macOS WebKit is supplied by the OS. No package-manager guessing, private browser
cache or automatic installation during a skill launch is added.

A missing runtime, failed setup or failed preactivation check leaves an existing
Fungi installation selected. System package managers own any partial runtime
setup and updates; Fungi does not roll back OS packages.

Only format-2 Wails packages are accepted. Format-1 Electron preview support
was removed during the 2026-10-26 cleanup; unsupported packages are rejected
before any package executable is run. Rollback between retained Wails releases
continues to work. Historical migration/rollback evidence below describes older
builds, not current compatibility.

An update lock serializes mutations. Short activation locks coordinate version
selection and cleanup; downloads do not hold this lock. Each running process
holds an OS file lease on its version. There is no daemon or PID polling.

Managed session receipts/callbacks use the stable command. Package-owned source
and helper paths are retained for draft reopening, including sessions in
alternate state directories. Keep the current version, previous version,
versions in use and referenced packages. A closed/cancelled/interrupted session
can still need old assets. A submitted session cannot reopen and its bytes live
in session storage. Cleanup inspection failures retain the affected version
and produce warnings. Removing an obsolete package never removes external
session storage. Session storage inside the managed installation is rejected.

Legacy portable receipts still reference their original executable; this
installer does not migrate or delete existing `dist/` builds. It also does not
snapshot skill directories outside its installation: skill authors must keep
those assets available as before. Retaining a recoverable session may retain an
older package indefinitely. The stable launcher itself is not self-replaced;
a future incompatible launcher/layout change needs an explicit migration.

## Build and verification recipe

Build on the target OS/CPU with Go, Node and native build tools:

```sh
node scripts/package.mjs --release v0.1.0-test.1
node scripts/package.mjs --release v0.1.0-test.2
```

Each creates `dist/release-VERSION/fungi-OS-ARCH`, a `.tar.gz` (Unix) or `.zip`
(Windows) archive and a `.sha256` sidecar. `RELEASE.json` identifies the version,
OS/CPU, Wails version and both executable hashes. Archive names use Go OS/CPU names, such
as `darwin-arm64`, `windows-amd64` and `linux-amd64`. Packaging alone does not publish.
BUILD.json records the source commit and whether the checkout was dirty.
The build includes dependency license notices and small HTML examples. It omits
a browser engine, full wiki/evidence, skill bundles and compiled demo helpers.

Run `go test -race -v ./...` on Unix (`go test -v ./...` on Windows), nested
renderer checks from [setup](../verification/setup.md), and
`node scripts/install-smoke.mjs ARCHIVE_A ARCHIVE_B`. The latter uses a local
HTTP fixture for the actual bootstrap, updates from the second local archive,
keeps one real Wails session open through activation, checks its PNG hash,
verifies the new format has no bundled runtime, and rolls back. Installed CLI
children have an empty PATH. Windows uses `-NoPath` to avoid persistent test
changes. Both supported Ubuntu targets exercise UI under Xvfb, with direct
media recorded as an accepted limitation in the renderer fixture.

On Linux, `node scripts/runtime-setup-smoke.mjs ARCHIVE ubuntu-22.04` (or
`ubuntu-24.04`) uses a disposable Docker container without WebKitGTK. It checks
missing-runtime launch failure, declined setup, an injected package-manager
failure, real distro runtime setup and reuse. It does not prove a desktop UI
or audible playback in that container. No runtime is installed on the user's
host by this verification recipe.

Also run the existing renderer smoke driver against the stable installed
command. Preserve the generated evidence before removing test-owned builds.
No test may delete older user/session-referenced `dist/` packages.

## Public release workflow

An annotated `vX.Y.Z-alpha.N` tag triggers the alpha workflow. It reuses the
foundation checks with the actual release version, then publishes those same
archives after all jobs pass. Only the Ubuntu 22-built Linux archive is shipped;
it must also pass on Ubuntu 24. Each selected BUILD.json must name the tagged
commit and clean source. Both bootstrap scripts and all archive/script SHA-256
sidecars accompany the GitHub prerelease. No release is marked stable or Latest.

Post-publication jobs use anonymous HTTPS downloads on Mac, Windows and both
Ubuntu baselines. They run the public installer, doctor, the public update path,
example discovery and the existing renderer fixture. These are hosted-download
checks, not clean consumer-machine qualification. The complete maintainer recipe
is in [RELEASING](https://github.com/rengwu/fungi/blob/main/docs/RELEASING.md).
Published version bytes are never replaced; corrections use a new version.

## Coverage and limits

The [Wails integration run](../verification/runs/2026-10-25-wails-integration.md)
records current package/runtime results and build identities. The older
[distribution run](../verification/runs/2026-10-23-distribution.md) retains
Electron-specific observations. Clean-machine security prompts,
signing/notarization, installer PATH behavior in fresh
interactive shells, power loss and future protocol/launcher migrations are not
inferred from a local or hosted success. Full Fungi-native jstack integration
remains deferred.

## Verification history

- [First public alpha](../verification/runs/2026-10-26-release.md) — publication and public-download checks; all public install/update paths passed; an intermittent Mac media probe failure remains recorded.

- [2026-10-26 cleanup](../verification/runs/2026-10-26-cleanup.md) — removed Electron preview compatibility; current Wails installation/update/rollback checks.

- [2026-10-26 offline discovery](../verification/runs/2026-10-26-discovery.md) — installed help, examples, update/rollback and native throwaway reply.

- [2026-10-25 Wails integration](../verification/runs/2026-10-25-wails-integration.md) — default packages, installed native input, update/rollback and explicit Linux runtime setup.

- [2026-10-23 distribution proof](../verification/runs/2026-10-23-distribution.md)

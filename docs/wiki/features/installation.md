---
title: Installation and updates
updated: 2026-10-23
---
# Installation and updates

**ID: FEAT-installation**

Implement the user's [command installation and update requirement](../planning.md#command-installation-and-updates-on-2026-10-23)
with one prebuilt package per OS/CPU and shared Go update logic. This is a
development implementation; published releases, signing and clean consumer
machine acceptance remain separate work. Skills use the same four page methods.

## Entry points

`fungi --version` reports JSON build identity. `fungi describe installation`
describes these commands:

| Command | Behavior |
| --- | --- |
| `fungi install --archive FILE --sha256 HASH [--dir DIRECTORY]` | Install a verified local package; default directory is `~/.fungi` |
| `fungi update` | Fetch the matching archive/checksum from the latest GitHub release |
| `fungi update --version vX.Y.Z` | Select a specific published release |
| `fungi update --archive FILE --sha256 HASH` | Apply an explicitly supplied local release |
| `fungi rollback` | Switch to the previous retained release without downloading |
| `fungi prune` | Remove obsolete versions not required for rollback or sessions |

Updates, rollback and pruning operate on managed installations. Unpacked
portable folders remain usable directly; they are not silently adopted.
Successful mutations emit one JSON `installation` result with the stable
command, current/previous versions, removed versions, cleanup warnings and
`reusedRuntimeBytes` (the bytes shared with an existing renderer).
Errors use the existing JSON error envelope and nonzero exit status.

The repository's `scripts/install.sh` and `scripts/install.ps1` download and
verify one archive, extract it privately and invoke the Go installer. POSIX
requires a shell, curl or wget, tar and shasum or sha256sum. Windows uses
PowerShell's download, ZIP and hashing commands. Neither requires Go, Node,
npm or Homebrew. Windows adds the user PATH entry unless `-NoPath` is given;
POSIX prints the directory to add, without editing shell startup files.

`FUNGI_INSTALL_DIR` selects a bootstrap destination. The bootstrap-only
`FUNGI_RELEASE_BASE_URL` permits an explicit mirror/test server; it does not
change the installed updater's GitHub source. There is no public installer
URL or published release claim yet. Ordinary `fungi update` needs release
assets to exist at `rengwu/fungi`.

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

When the selected and incoming packages use the same Electron version, compare
the runtime's regular files by permissions, size and SHA-256. Identical files
can share filesystem hard links; unsupported links fall back to the verified
copies. Framework symlinks remain unchanged. Package files must stay immutable:
editing a shared file in place would affect each version using it. Removing an
old package removes its links without invalidating the retained version.
This reduces retained update storage, not the first download, update traffic
or staging space. No separate runtime cache or skill-author API is introduced.

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

From a development checkout with the pinned runtime already downloaded:

```sh
node scripts/package.mjs --release v0.1.0-test.1
node scripts/package.mjs --release v0.1.0-test.2
```

Each creates `dist/release-VERSION/fungi-OS-ARCH`, a `.tar.gz` (Unix) or `.zip`
(Windows) archive and a `.sha256` sidecar. `RELEASE.json` identifies the version,
OS/CPU, Electron version and CLI hash. Archive names use Go OS/CPU names, such
as `darwin-arm64`, `windows-amd64` and `linux-amd64`. No release is published.
The build keeps Electron intact, strips Go debug information, uses standard
maximum standard-library archive compression and omits full wiki/evidence, skill bundles and compiled
demo helpers. Small HTML examples and upstream license notices remain.

Run `go test -race -v ./...` on Unix (`go test -v ./...` on Windows) and
`node scripts/install-smoke.mjs ARCHIVE_A ARCHIVE_B`. The latter uses a local
HTTP fixture for the actual bootstrap, updates from the second local archive,
keeps one real Electron session open through activation, checks its PNG hash,
verifies identical runtimes share files and rolls back. Its installed CLI
children have an empty PATH. Windows uses
`-NoPath` to avoid persistent test changes. `--core-only` deliberately omits UI
on the unsupported Ubuntu 24.04 baseline. A finite automated fixture gate is
not a human-facing polling mechanism.

Also run the existing renderer smoke driver against the stable installed
command. Preserve the generated evidence before removing test-owned builds.
No test may delete older user/session-referenced `dist/` packages.

## Coverage and limits

The [distribution run](../verification/runs/2026-10-23-distribution.md) records
actual results and build identities. Clean-machine security prompts, public
release downloads, signing/notarization, installer PATH behavior in fresh
interactive shells, power loss and future protocol/launcher migrations are not
inferred from a local or hosted success. Full Fungi-native jstack integration
remains deferred.

## Verification history

- [2026-10-23 distribution proof](../verification/runs/2026-10-23-distribution.md)

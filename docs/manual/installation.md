# Installation and updates

Engine packages contain one Fungi executable, the offline manual, build metadata
and Go license. The executable also works alone. No webview, graphical desktop,
Go, Node or npm is needed to run the engine. External helpers bring their own
requirements. Use `fungi.exe` on Windows.

This checkout is the 0.2 development line (protocol 2, package format 3).
Previously published 0.1 alpha downloads contain the old bundled UI.

## Local build and installation

Developers need Go 1.26 and Node 22+ for packaging:

```sh
node scripts/package.mjs --release v0.2.0-test.1
```

The command prints the archive path and writes its SHA-256 sidecar. Extract the
archive and use the extracted executable, substituting those exact values:

```sh
./fungi install --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256 --dir "$HOME/.fungi-engine"
export PATH="$HOME/.fungi-engine/bin:$PATH"
fungi doctor
```

Add that PATH entry to your shell configuration to retain it. A separate install
keeps protocol-1 sessions and their UI installation available during migration.
The default managed install destination is `~/.fungi` when `--dir` is omitted.

## Release bootstrap

Once an engine release is published, use its shell or PowerShell installer with
an explicit version. Script options:

```text
sh install.sh [--version vX.Y.Z] [--dir DIRECTORY]
powershell -File install.ps1 [-Version vX.Y.Z] [-Directory DIRECTORY] [-NoPath]
```

Installers download the matching OS/CPU archive and verify its checksum. The
current release workflow targets macOS arm64, Windows x64 and Linux x64.
Cross-building additional Go targets is supported by the packaging script but
does not establish runtime qualification or published availability.

The shell installer needs curl or wget, tar, and shasum or sha256sum; it prints
PATH instructions. PowerShell uses Windows built-ins and updates user/current
PATH unless `-NoPath` is set. Neither installer installs system runtimes.
`FUNGI_INSTALL_DIR` sets a bootstrap destination; `FUNGI_RELEASE_BASE_URL`
selects a bootstrap test mirror. Omitted versions select latest stable.

## Managed versions

```sh
fungi update --version vX.Y.Z
fungi update --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
fungi rollback
fungi prune
```

Updates verify the archive checksum, platform, engine digest and executable's
protocol/version before activation. Versions are immutable. Existing sessions
keep running; rollback selects the previous retained engine release. Prune
protects current, previous, active and session-referenced versions. External
helpers and their assets are managed by their provider.

Keep session storage outside the installation. Portable copies do not
self-update. There are no background updates. Format-1/2 UI packages are rejected,
including on rollback. Checksums detect corruption; publisher signing and OS
security policy remain separate. See `fungi help migration` for old sessions.

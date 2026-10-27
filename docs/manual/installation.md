# Installation and runtime

Packages target macOS 26+ arm64, Windows x64 and Ubuntu 22.04/24.04 x64.
Download only your OS/CPU archive from https://github.com/rengwu/fungi/releases.
Keep the extracted folder intact: CLI, Wails helper, examples, manual and
licenses. No Go, Node, npm or bundled browser engine is required.
Use `fungi.exe` on Windows.

## Install

Shell:

```sh
curl -fsSL https://github.com/rengwu/fungi/releases/download/v0.1.0-alpha.2/install.sh | sh -s -- --version v0.1.0-alpha.2 && export PATH="$HOME/.fungi/bin:$PATH"
```

PowerShell:

```powershell
Invoke-WebRequest -UseBasicParsing https://github.com/rengwu/fungi/releases/download/v0.1.0-alpha.2/install.ps1 -OutFile "$env:TEMP/fungi-install.ps1"; & "$env:TEMP/fungi-install.ps1" -Version v0.1.0-alpha.2
```

Script options (scripts also ship in the source checkout):

```text
sh scripts/install.sh [--version vX.Y.Z] [--dir DIRECTORY] [--setup-runtime]
powershell -File scripts/install.ps1 [-Version vX.Y.Z] [-Directory DIRECTORY] [-SetupRuntime] [-NoPath]
```

Alphas require an explicit version. Omitting it selects latest stable; none is
published yet. `FUNGI_INSTALL_DIR` sets the default destination.
`FUNGI_RELEASE_BASE_URL` overrides bootstrap downloads for testing, not updates.

The shell installer needs curl or wget, tar, and shasum or sha256sum. It prints
PATH instructions; the one-liner above also sets PATH for the current terminal.
PowerShell uses Windows built-ins and updates user and current PATH unless
`-NoPath` is set. These tools are unnecessary for running an extracted package.

## Managed updates

```sh
fungi install --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
fungi update --version v0.1.0-alpha.2
fungi update # latest stable, once published
fungi update --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
fungi rollback
fungi prune
```

Install defaults to `~/.fungi`; override with `--dir DIRECTORY`. Add its `bin`
directory to PATH and use `bin/fungi`. Updates verify a GitHub release or a local
archive with a 64-character hexadecimal SHA-256 checksum. Checksums detect
corruption; they are not independent publisher signatures.

Versions are immutable. Updates validate the candidate and runtime before
activation; existing sessions keep their running version. Rollback selects the
previous retained release. Prune protects current, previous, active and
session-referenced versions. Keep session state outside the installation.
These commands require a managed install; portable folders do not self-update.
There are no background updates. Only format-2 Wails packages are accepted,
including for rollback; Electron previews are unsupported.

## System runtime

`fungi doctor` checks the webview without opening a window.

| Platform | Runtime |
| --- | --- |
| macOS | System WKWebView |
| Windows | Microsoft WebView2 |
| Linux | GTK3 and WebKitGTK 4.1 |

Exit 3 means a missing runtime; other failures can indicate a missing helper.
Normal `open`, `doctor` and `update` never install system software. Request setup
explicitly with `install.sh --setup-runtime` on supported Ubuntu or
`install.ps1 -SetupRuntime` on Windows. Setup may require network/admin access.
Other Linux distributions need their own WebKitGTK 4.1 setup.

On startup failure, read the error and worker log location, then run `fungi doctor`.
A window requires a desktop session; a successful runtime check does not verify
desktop, GPU or media support. Do not change security policy or install runtimes
automatically to hide a failure.

## Alpha limits

Downloads are unsigned and unnotarized; OS policy may block them. Installers do
not bypass security or execution policy. Ubuntu missing-runtime setup is tested;
Windows setup on a machine lacking WebView2 remains unverified. Physical
Windows/Linux input, other OS/CPU combinations and broader sleep/wake coverage still
need qualification. Hosted checks do not establish clean consumer-machine support.

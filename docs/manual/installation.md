# Installation and runtime

Obtain only the archive matching this machine's OS/CPU and keep its folder
together: Go CLI, Wails helper, examples, manual and licenses. No Go, Node,
npm or bundled browser engine is needed by the user. Use `fungi.exe` on Windows.
Public alpha downloads are at https://github.com/rengwu/fungi/releases.
They are not publisher-signed or notarized; OS security policy may block them.
The initial targets are macOS 26+ arm64, Windows x64 and Ubuntu 22.04/24.04 x64.

`fungi doctor` checks the system webview without a window:

| Platform | Runtime |
| --- | --- |
| macOS | System WKWebView |
| Windows | Microsoft WebView2 |
| Linux | GTK3 and WebKitGTK 4.1; current package baseline Ubuntu 22.04/24.04 |

Doctor exit 3 means a missing runtime; other failures include a missing helper.
A successful runtime check does not prove desktop/GPU/media availability.
Normal `open`, `doctor` and `update` never install system software. Repository
bootstrap scripts offer explicit `install.sh --setup-runtime` on supported
Ubuntu or `install.ps1 -SetupRuntime` on Windows. Other distributions need
their own WebKitGTK 4.1 setup. Runtime installation may need network/admin
access. Missing-runtime setup has been tested in Ubuntu containers; a Windows
machine actually lacking WebView2 remains a verification gap.

## Bootstrap scripts

Download the matching install.sh or install.ps1 asset from the chosen release.
For the first alpha:

```sh
curl -fsSL https://github.com/rengwu/fungi/releases/download/v0.1.0-alpha.1/install.sh | sh -s -- --version v0.1.0-alpha.1
export PATH="$HOME/.fungi/bin:$PATH"
```

```powershell
Invoke-WebRequest -UseBasicParsing https://github.com/rengwu/fungi/releases/download/v0.1.0-alpha.1/install.ps1 -OutFile "$env:TEMP/fungi-install.ps1"
& "$env:TEMP/fungi-install.ps1" -Version v0.1.0-alpha.1
```

Downloaded scripts (also in the source checkout) accept:

```text
sh scripts/install.sh [--version vX.Y.Z] [--dir DIRECTORY] [--setup-runtime]
powershell -File scripts/install.ps1 [-Version vX.Y.Z] [-Directory DIRECTORY] [-SetupRuntime] [-NoPath]
```

Omitted version selects the latest stable release. Alphas require an explicit
version; there is no stable release yet. `FUNGI_INSTALL_DIR` supplies the default target;
`FUNGI_RELEASE_BASE_URL` overrides the release download base for local testing.
The shell script needs curl or wget, tar and shasum or sha256sum. PowerShell uses
Windows built-ins and adds the user PATH unless `-NoPath` is given; the shell
script prints PATH instructions. They are not extra dependencies for running
an extracted package. Neither script bypasses OS security or execution policy.

## Managed updates

```sh
fungi install --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
fungi update --version v0.1.0-alpha.1
fungi update # latest stable; available after a stable release is published
fungi update --version vX.Y.Z
fungi update --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
fungi rollback
fungi prune
```

Install defaults to `~/.fungi`; override with `--dir DIRECTORY`. Add its `bin`
directory to PATH. Use the stable `bin/fungi` command thereafter. The updater
verifies a published GitHub release (latest stable or selected version), or an
explicit local archive/checksum. Checksum input is 64 hexadecimal characters.
Use the new release's exact version for alpha updates. Checksums detect corrupt
downloads; they are not a separate publisher signature.

Versions are immutable directories. Updates validate the candidate and its
runtime before activation; existing sessions retain their running version.
Rollback selects the previous retained release. Prune retains current,
previous, active and session-referenced versions. Portable folders do not
self-update; create a managed install for these commands. No background update
or OS policy changes occur. Keep session state outside the installation.

Only format-2 Wails packages are supported for installation and rollback.
Electron preview packages are no longer accepted. Physical Windows/Linux desktop input,
other OS/CPU combinations, signing/notarization and long Wails sleep/wake still
need release acceptance. Hosted checks are not clean consumer-machine proof.

For startup failure, read the returned error and worker log location, then run
`fungi doctor`. A desktop session is required to open a window. Do not change
host security policy or install runtimes automatically to hide a failure.

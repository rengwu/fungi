# Installation and runtime

Obtain only the archive matching this machine's OS/CPU and keep its folder
together: Go CLI, Wails helper, examples, manual and licenses. No Go, Node,
npm or bundled browser engine is needed by the user. Use `fungi.exe` on Windows.
This is still a development preview; public releases/signing are not established.

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

From a source checkout, the download scripts accept:

```text
sh scripts/install.sh [--version vX.Y.Z] [--dir DIRECTORY] [--setup-runtime]
powershell -File scripts/install.ps1 [-Version vX.Y.Z] [-Directory DIRECTORY] [-SetupRuntime] [-NoPath]
```

Omitted version selects latest. `FUNGI_INSTALL_DIR` supplies the default target;
`FUNGI_RELEASE_BASE_URL` overrides the release download base for local testing.
The shell script needs curl or wget, tar and shasum or sha256sum. PowerShell uses
Windows built-ins and adds the user PATH unless `-NoPath` is given; the shell
script prints PATH instructions. These scripts require published assets (or a
test release server). They are not extra dependencies for running an extracted
package.

## Managed updates

```sh
fungi install --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
fungi update
fungi update --version vX.Y.Z
fungi update --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
fungi rollback
fungi prune
```

Install defaults to `~/.fungi`; override with `--dir DIRECTORY`. Add its `bin`
directory to PATH. Use the stable `bin/fungi` command thereafter. The updater
verifies a published GitHub release (latest or selected version), or an explicit
local archive/checksum. Checksum input is 64 hexadecimal characters. Public
download commands require actual release assets; local archives work now.

Versions are immutable directories. Updates validate the candidate and its
runtime before activation; existing sessions retain their running version.
Rollback selects the previous retained release. Prune retains current,
previous, active and session-referenced versions. Portable folders do not
self-update; create a managed install for these commands. No background update
or OS policy changes occur. Keep session state outside the installation.

Old Electron preview updaters reject package format 2 and need the new bootstrap
once; legacy rollback support remains. Physical Windows/Linux desktop input,
other OS/CPU combinations, signing/notarization and long Wails sleep/wake still
need release acceptance. Hosted checks are not clean consumer-machine proof.

For startup failure, read the returned error and worker log location, then run
`fungi doctor`. A desktop session is required to open a window. Do not change
host security policy or install runtimes automatically to hide a failure.

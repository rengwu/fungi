---
title: Verify a portable desktop build
updated: 2026-10-25
---
# Verify a portable desktop build

Run on the target desktop after copying/extracting the **whole** portable
Wails package. Go, Node and npm are unnecessary. macOS uses WKWebView;
Windows needs WebView2 and Linux needs WebKitGTK 4.1. Run `fungi doctor`
first. The bootstrap can explicitly install missing Windows/Ubuntu runtimes;
ordinary launches never do. See [installation](../features/installation.md).

Record the OS/version/architecture, `BUILD.json` and runtime report. Hosted
fixtures prove rendering and persistence on their runner environments; this
recipe separately checks visible input and clean-machine behavior. Preserve
startup errors and the session's `worker.log`; do not bypass OS security
protections. Known Linux local-media and best-effort storage limits are in
[platforms](../platforms.md).

## Windows PowerShell

Open PowerShell in the extracted directory. Use a disposable state directory:

```powershell
$env:FUNGI_STATE_DIR = Join-Path $env:TEMP ('fungi-check-' + [guid]::NewGuid().ToString('N'))
Get-Content .\BUILD.json
.\fungi.exe doctor
.\fungi.exe describe
$opened = .\fungi.exe open .\examples\compare\index.html --delivery manual | ConvertFrom-Json
$opened
```

In the comparison window, move the divider, select a preference and enter a test
note. Wait for **Draft saved**, then close using the window's close button.
Collect once and reopen:

```powershell
.\fungi.exe collect $opened.session
$reopened = .\fungi.exe reopen $opened.session --delivery manual | ConvertFrom-Json
$reopened
```

The first outcome should be `closed`. The reopened window should restore the
divider, preference and note. Click **Send review**, then collect once:

```powershell
$result = .\fungi.exe collect $reopened.session | ConvertFrom-Json
$result | ConvertTo-Json -Depth 20
Get-FileHash -Algorithm SHA256 $result.artifacts[0].path
```

Expect `submitted`, matching data and a `comparison` PNG whose SHA-256 matches
the returned metadata. Open that PNG to check the captured view. A blocked
launch, lost draft, missing image or failed hash is a failed observation; keep
the session directory for diagnosis.

For a cancellation check, launch the comparison again with manual delivery,
click **Cancel**, and collect that exact session once. Expect `cancelled`.
No script execution-policy or machine configuration change is required.

## Linux or another Mac

Use the same sequence with POSIX commands and the paths from each receipt:

```sh
export FUNGI_STATE_DIR="$(mktemp -d)/sessions"
cat BUILD.json
./fungi doctor
./fungi describe
./fungi open ./examples/compare/index.html --delivery manual
# After closing the window with an acknowledged draft:
./fungi collect EXACT_SESSION_ID
./fungi reopen EXACT_SESSION_ID --delivery manual
# After submitting the reopened window:
./fungi collect NEW_SESSION_ID
```

Use `sha256sum ARTIFACT_PATH` on Linux or `shasum -a 256 ARTIFACT_PATH` on macOS.
Linux needs a graphical desktop and compatible WebKitGTK 4.1 libraries.
A container-only runtime installation check does not validate physical input,
clipboard, GPU behavior or clean desktop distribution. Record any failed
prerequisite and whether explicit runtime setup was used.

## Return the observations

Keep the build metadata, terminal results, any exact errors and the test-only
state directory. State whether draft restoration and the PNG looked correct.
Record automatic agent delivery separately using an already supported native
host route; this recipe deliberately uses manual delivery to isolate desktop
behavior. No model polling or repeated status calls are needed.

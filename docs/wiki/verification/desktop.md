---
title: Verify a portable desktop build
updated: 2026-10-22
---
# Verify a portable desktop build

Run on the target desktop after copying/extracting the **whole** portable
directory. Go, Node, npm and a separately installed browser are unnecessary.
This recipe is prepared for Windows and Linux; those desktop paths remain
unvalidated until executed there. Record the OS/version/architecture and
`BUILD.json` with the results. Do not bypass OS security protections or disable
Chromium's sandbox if startup fails; record the actual error.

## Windows PowerShell

Open PowerShell in the extracted directory. Use a disposable state directory:

```powershell
$env:FUNGI_STATE_DIR = Join-Path $env:TEMP ('fungi-check-' + [guid]::NewGuid().ToString('N'))
Get-Content .\BUILD.json
.\fungi.exe describe
.\fungi.exe open --exec .\examples\helper.exe --delivery wait
$opened = .\fungi.exe open .\examples\compare\index.html --delivery manual | ConvertFrom-Json
$opened
```

The helper should return `submitted` with a readable text artifact. In the
comparison window, move the divider, select a preference and enter a test
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
./fungi describe
./fungi open --exec ./examples/helper --delivery wait
./fungi open ./examples/compare/index.html --delivery manual
# After closing the window with an acknowledged draft:
./fungi collect EXACT_SESSION_ID
./fungi reopen EXACT_SESSION_ID --delivery manual
# After submitting the reopened window:
./fungi collect NEW_SESSION_ID
```

Use `sha256sum ARTIFACT_PATH` on Linux or `shasum -a 256 ARTIFACT_PATH` on macOS.
Linux needs a supported desktop environment and the shared OS libraries used
by Electron. A container-only Go pass does not validate windowing, graphics,
clipboard, sandbox startup or a clean distribution. Record missing library or
sandbox errors instead of prescribing installation as a passing no-setup test.

## Return the observations

Keep the build metadata, terminal results, any exact errors and the test-only
state directory. State whether draft restoration and the PNG looked correct.
Record automatic agent delivery separately using an already supported native
host route; this recipe deliberately uses manual delivery to isolate desktop
behavior. No model polling or repeated status calls are needed.

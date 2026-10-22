---
title: Windows agent handoff
updated: 2026-10-22
---
# Windows agent handoff

Continue Fungi's physical Windows verification from this checkout. The user
requested this handoff so an agent on their Windows machine can perform it
after pulling. No earlier conversation, Mac-local files, or sibling jstack
checkout is needed. [Their response](../journal/2026-10-22-grill-b8ce5a10-7261-48c1-9d80-ee54e103ec5e.md)
requests these instructions; it does **not** report a passing Windows test.

Hosted Windows Server already passes the core and renderer fixture. What
remains is the portable package's behavior on this actual desktop: launch,
visible input, acknowledged drafts, close/reopen, PNG submission and cancel.
Read [the verified baselines](runs/2026-10-22-baselines.md) and
[platform limits](../platforms.md) for context. Keep Go + Electron, the four
page methods, portable distribution and no-harness-setup behavior intact.

## 1. Get the correct checkout

Run in the Fungi repository using native Windows PowerShell, not WSL. The
work is on **`foundation-hardening`**, not `main`. Preserve local changes; do
not reset or clean them to switch branches.

```powershell
git status --short --branch
git fetch origin
git switch foundation-hardening
git pull --ff-only origin foundation-hardening
git rev-parse HEAD
```

Record the revision, dirty files, Windows version/build, CPU architecture and
agent/harness version. Use a new evidence directory for this run:

```powershell
$FungiRoot = (Get-Location).Path
$FungiRunId = 'windows-' + (Get-Date -Format 'yyyy-MM-dd-HHmmss')
$FungiEvidence = Join-Path $FungiRoot ('verification-output\' + $FungiRunId)
New-Item -ItemType Directory -Path $FungiEvidence | Out-Null
[System.Environment]::OSVersion.VersionString
Get-CimInstance Win32_OperatingSystem | Select-Object Caption, Version, BuildNumber, OSArchitecture
```

Keep the absolute paths in your task notes. Separate agent shell calls may
not retain PowerShell variables or environment overrides; restore their
recorded values when opening a new shell.

## 2. Build the portable folder

`dist/`, executables and `renderer/node_modules/` are ignored by Git. **A pull
does not download a ready-to-run package.** From a source checkout, building
requires Go 1.26.x, Node 22.12+ and npm. The hosted baseline used Go 1.26.4 and
Node 22.19.0. These are development prerequisites; the completed package must
run without invoking them. No Visual Studio or C compiler is required by the
Go packaging command (`CGO_ENABLED=0`).

Check existing tools first. If one is unavailable, record the missing build
prerequisite; that is not a renderer failure. A complete prebuilt Windows
package can instead be supplied and used for the desktop steps, but none is
stored in this repository. Do not assume the Mac's local archive is present.

```powershell
go version
node --version
npm.cmd --version
node -p "process.platform + ' ' + process.arch"
```

The commands below target Windows x64, matching the current hosted baseline.
Record an ARM/emulation environment explicitly rather than describing it as
native x64 coverage. Run each build command from the repository root and
check its exit code; stop the dependent build steps on a nonzero exit.

```powershell
npm.cmd ci --prefix renderer
node .\renderer\node_modules\electron\install.js
$FungiPackage = (node .\scripts\package.mjs --platform win32 --arch x64).Trim()
if ($LASTEXITCODE -ne 0) { throw 'Portable packaging failed' }
$FungiExe = Join-Path $FungiPackage 'fungi.exe'
Get-Content (Join-Path $FungiPackage 'BUILD.json')
Get-FileHash -Algorithm SHA256 $FungiExe
```

The Electron download is a build-time step, pinned by the lockfile. Preserve
the printed package path and its `BUILD.json`; verify the executable hash
matches `binarySHA256`. Never select the newest directory by guessing or use
a partially built folder. `runtimeVerified: false` in build metadata means
packaging itself does not certify runtime behavior.

## 3. Run the existing automated checks

From the repository root, run and retain each command's output and exit code:

```powershell
go vet ./...
go test -v ./...
node --test
node .\docs\wiki\_build.mjs --check
node .\scripts\renderer-smoke.mjs $FungiExe
```

The current suites contain 18 root Go tests and seven Node tests. The real
renderer fixture opens disposable windows, cancels, reopens, submits and
verifies the PNG hash. It requires all 21 renderer checks by default and
prints its `evidence.json` path. It clears the child PATH and caller ID,
showing that this packaged workflow uses its bundled runtime. Keep that
evidence directory. Do not relax WebGL checks or disable the sandbox to
manufacture a pass.

Some Go/Node tests and the renderer fixture create a file symlink to verify
path confinement. An ordinary Windows account may lack that privilege. If
fixture construction fails with a symlink permission error, retain the
failure, mark the affected checks unexecuted with that reason, and continue
the independent desktop journey below. Do not enable Developer Mode or run
elevated just to make the suite green. Fungi's normal HTML workflow does not
require the user to create a symlink. Other failures need diagnosis; do not
classify them as this known fixture prerequisite without evidence.

## 4. Exercise the visible desktop journey

Follow [the portable desktop checklist](desktop.md#windows-powershell), also
copied into `$FungiPackage\VERIFY.md`. Start in the package folder:

```powershell
Set-Location $FungiPackage
$env:FUNGI_STATE_DIR = Join-Path $FungiEvidence 'sessions'
```

Use this state directory throughout; omit the checklist's alternative
temporary-directory assignment. This keeps receipts, drafts, results,
artifacts and `worker.log` together. Preserve the receipt JSON from each
launch before interacting with its window.

Use the available native desktop driver to enter a distinctive **synthetic**
note such as `WINDOWS-ACCEPTANCE-<run ID>`, move the divider, choose a
preference and observe **Draft saved**. Close, collect once, reopen, inspect
restored values, submit and verify the PNG visually and by SHA-256. Launch a
separate comparison and cancel it. Match exact window titles/processes and
session IDs; do not interact with the user's unrelated Electron windows.

If the agent has no desktop-control tool, let the user perform those visible
steps and supply the observations. An automated self-submitting fixture is
not evidence that physical input was exercised. Do not replace the visible
steps with DOM setters, edited session files or a fabricated result.

This journey deliberately uses `manual` delivery to isolate desktop behavior.
Collect only after the corresponding observed close/submit/cancel, never in
a polling loop. If waiting for the human, end the turn and resume when they
report completion. Automatic harness return is a separate test and must use
an already supported native route, with no plugin/configuration setup.

On failure, preserve the exact command, error, revision and session. A
`renderer_startup_failed` result points to its diagnostic log. Record any OS
launch block without bypassing it. Within the existing foundation scope,
reproduce actual product defects, make the smallest repair, and rerun the
affected checks on a newly identified build. Keep the failed baseline too.

## 5. Save the result for the next agent

Create `docs/wiki/verification/runs/<FungiRunId>.md` with a matching H1/title,
the actual date and `result: Passed`, `Failed` or `Untested`. Record the tested
Git revision/dirty source, `BUILD.json`, binary hash, OS/architecture, driver,
commands and scenario-level observations. Any selected failed scenario makes
the aggregate Failed; use Untested when required scenarios remain unexecuted
and none failed. Separate fixture setup failures from observed product failures.

Copy relevant test-only UTF-8 logs, JSON receipts/results, screenshots and the
returned PNG under `docs/wiki/evidence/<FungiRunId>/`; link them relatively
from the report. Keep private machine/account data out of committed logs.
State explicitly which of launch, draft recovery, PNG/hash and cancellation
passed, and whether the agent or human operated the window. Do not claim a
clean-machine install-free launch merely because a source build worked on a
machine with development tools; record that limitation.

Link the report from `features/sessions.md`, `features/interfaces.md` and the
wiki index. Update `planning.md` and `platforms.md` only to the extent the
observations support. Preserve older reports. From the repository root:

```powershell
Set-Location $FungiRoot
node .\docs\wiki\_build.mjs
node .\docs\wiki\_build.mjs --check
git diff --check
git status --short
```

Commit only the intended report/evidence/docs and any justified repair to
`foundation-hardening`, then push that branch to `rengwu/fungi`. Leave
unrelated changes alone. Do not commit `dist/`, dependencies or raw session
directories. Report the result, commit and material gaps to the user. No
new grill is needed merely to start this already requested verification.

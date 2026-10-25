---
title: Verification setup
updated: 2026-10-23
---
# Verification setup

Run from `/Users/rengwu/Desktop/Projects/fungi`. Early runs predate Git initialization; current runs record the Git base, dirty source hashes and binary SHA-256.

For the user's Windows checkout, follow the [Windows agent handoff](windows.md).
It includes the correct branch, source build, actual desktop checks and result
recording without relying on this Mac's paths or a sibling jstack checkout.

## Build and launch

Development prerequisites on the tested machine: Go 1.26.4, Node 22.19.0, npm 10.9.3, macOS arm64. Electron is pinned to 44.9.0 in `renderer/package-lock.json`.

```sh
npm ci --prefix renderer
node renderer/node_modules/electron/install.js
go build -o fungi .
./fungi describe
go test -race -v . ./examples/helper
go vet . ./examples/helper
node scripts/package.mjs
```

The explicit Electron download is a **build** step; the portable output includes the runtime. Packaging produces a unique directory under `dist/`, preserving framework symlinks. Do not test a half-built directory. Read its `BUILD.json` and hash its `fungi` executable.

Use `FUNGI_STATE_DIR="$PWD/verification-output/sessions"` for local UI fixtures. Use a separate root for each harness batch. `open` returns only after readiness, or keeps one pending pipe in `wait` mode. Never poll `collect` while awaiting a person.

## UI driver and fixtures

Native UI verification used the connected desktop driver against the exact Electron app path. Multiple Electron development apps share a bundle ID, so select the full app path and verify the window title. Bringing up an app through the driver sometimes also launched Electron's default page; terminate only the extra test-created launcher process, never unrelated Electron apps.

The generated PNG fixture is 64×64 RGB, 180 bytes, with SHA-256 `1bd7e798fbfd3bd20de881435c7bfa7811db8466e73714a8ab7c1af6b7460e9a`. Its source bytes are retained in evidence. The current reference form exposes accessible labels for text fields, the native file picker, and “Paste image for Reference image.”

The successful clipboard fixture used a separate form with one image field, keeping its paste target visible. Copying pixels from macOS Preview and pasting into that target yielded a 401-byte, 64×64 PNG with SHA-256 `83e160394a50c08746093b421e7bbf7865d6591aac9f9b2d1f0850bea504c2d9`. Both the original and clipboard-encoded bytes are retained. The crash fixture had one text field: after “Draft saved,” only its identified page-renderer PID was terminated. The restored window visibly contained the acknowledged value before it was closed.

Do not overwrite the user's edits while they are trying a demo. Preserve acknowledged drafts. A changing desktop focus or failed clipboard setup is a driver limitation until separately diagnosed; it is not a product pass.

## Harness probes

Commands installed at the time of the run:

| Harness | Version | Test configuration |
| --- | --- | --- |
| Codex | 0.163.0 | Disposable `exec`, ephemeral JSON output, workspace-write sandbox for Fungi storage |
| Claude Code | 2.1.302 | Safe mode, no session persistence, Bash only, exact fixture command allowed, bounded turn count |
| OpenCode | 1.18.34 | `run --pure --format json`, canonical disposable working directory |
| Grok Build | 1.0.62 | Headless JSON events, bounded turns, no subagents/web search, exact fixture command |
| Kimi Code | 2.0.1 | Print/JSON attempt; blocked by missing provider credential |
| Pi | 1.3.0 | No extensions/MCP/skills/context files/session persistence; Bash only; invocation-only `openai-codex/gpt-6.1-sol` |

Pi's configured default model was rejected by its provider. The invocation-only override used its existing provider; no saved model setting changed. Raw streams remain in the disposable probe directory; the wiki contains sanitized tool calls, terminal texts and task notifications without account credentials or model reasoning payloads.

The probe driver imposes a finite external test deadline (120 or 240 seconds) and terminates only its own process groups on timeout. This is a test bound, not a Fungi human timeout.

## No-setup follow-up fixtures

The follow-up lives under `verification-output/no-setup-verification/`. `run-ui-probe.py` launches one disposable harness conversation with an explicit fixture command and a 240-second external test bound. Each form has a unique title, response marker and state directory; the connected desktop driver chooses the exact packaged Electron app path. The user-edited source demo is separate.

Kimi's successful follow-up reports 2.2.0. It uses `run_in_background: true` and `disable_timeout: true`; no authentication files are read. Claude's headless test deliberately supplies 600,000 ms and does not represent unlimited interactive waiting. Non-Codex child probes do not inherit this test coordinator's `CODEX_THREAD_ID`.

Do not name a Fungi form `opencode.json` in an OpenCode working directory: that name is reserved for its configuration and may be rewritten. The corrected fixture is `opencode-form.json`. The initial collision is recorded as a failed fixture attempt, not hidden or treated as a Fungi product failure.

For manual delivery, let the launch return and the harness end its turn, submit through the native window, then resume the exact recorded harness session with an explicit “done.” It must call the receipt's `collect` argv with `collectEnv` once. This tests manual fallback and cannot be counted as automatic completion.

## Observation and cleanup

Inspect `launch.json`, `draft.json`, `result.json`, artifact bytes/hashes and `delivery.json` separately. The OS file lock proves owner liveness without a PID guess. Preserve the result and evidence before removing disposable data.

Do not delete a user-modified demo draft. No persistent harness configuration or authentication was changed. Public signing/notarization, a clean machine with network disabled, Windows runtime and Linux desktop UI verification, long sleep/wake cycles and abrupt power-loss behavior remain outside current evidence. The later Linux core run is recorded below.

## Skill authoring and idle probes

The simpler questionnaire lives at `skills/grill-form/assets/form.html`; pass its question JSON through `--input`. Native fixtures and receipts live under `verification-output/skill-proof/`. The `example.json` is explicitly a demonstration, so synthetic test answers are safe there. Never submit test answers into an actual interview.

```sh
node --test scripts/grill-form.test.mjs
go build -o verification-output/idle-helper ./scripts/idle-helper
./fungi open skills/grill-form/assets/form.html --input skills/grill-form/example.json --delivery manual
```

For native recovery: enter a choice/note, select the known reference image, wait for “Draft saved,” close, collect once, and reopen. The initial session must stay `closed`. To test abrupt owner loss, identify the exact `_run SESSION` process in this disposable root before killing it; never terminate by generic Electron/Fungi process name. Collect the interrupted session once, reopen it, and verify the acknowledged draft and image. The submitted result must preserve the image hash.

For the Linux core recipe, copy only root Go sources/module, `examples/helper`, `scripts/protocol-fixture` and any selected `scripts/idle-helper` fixture into an isolated directory. The [recorded Docker command](../evidence/2026-10-20-skill-proof/linux-command.json) pins the official Go image, uses an unprivileged user, disables networking and mounts that source read-only. Its disposable container removes itself. The passed run is Linux arm64 under OrbStack; it is not a Linux desktop verification.

The idle helper takes `--input` JSON with `seconds` (1–86400) and `marker`. It acknowledges a draft, announces readiness, sleeps locally once, then submits timings and a JSON artifact. The five-second smoke check must pass before launching a long run. For a real callback test, run with `--delivery auto` from the intended conversation, preserve the receipt, and end the agent turn. Its one completion notice supplies the exact executable/session/environment for one collection. For other harnesses, use their already documented completion route. Do not wrap the wait in repeated tool calls.

The completed two-hour probe is session `622f59b42e05840dd408622f23c8416a`, with state beneath `verification-output/skill-proof/idle-long-sessions`. Its original executable is in `dist/fungi-darwin-arm64-2026-10-20T09-41-25-905Z`. Its [result](../evidence/2026-10-20-skill-proof/idle-long-result.json) records one collection after the exact-origin notice, matching timing artifact bytes and exit of the helper/owner. Retain the evidence. For future runs, keep the launched package and helper available until completion; do not suspend the machine without settling timing with its user. Record both the helper's elapsed measurement and the UTC timestamp difference when sleep occurs.

The user manually completed sleep/wake fixture `1897358a56ceff72064a252004170716`, using the same package and `verification-output/skill-proof/sleep-wake-sessions`. Its [baseline](../evidence/2026-10-20-skill-proof/sleep-wake-baseline.json) records the initial note/image and preparation time. After its completion notice, one collection returned the preserved note/image and user's post-wake edit. A minimal [power-log excerpt](../evidence/2026-10-20-skill-proof/sleep-wake-power.log) independently confirms actual sleep/wake within the interval. The [result record](../evidence/2026-10-20-skill-proof/sleep-wake-result.json) retains the callback, data and byte validation. For future repetitions, obtain the user's timing, prepare a fresh fixture, let them sleep/wake and submit, then collect once; inspect only relevant `Sleep`, `Wake` and `DarkWake` event rows since the baseline, excluding scheduled `Wake Requests`. Do not select an answer or submit on the user's behalf.

## Platform foundation checks

```sh
node --test scripts/*.test.mjs
go test -race -v . ./examples/helper ./scripts/protocol-fixture
go vet . ./examples/helper ./scripts/protocol-fixture
node scripts/package.mjs
node scripts/renderer-smoke.mjs /absolute/package/fungi
node scripts/package.mjs --platform win32 --arch x64
node scripts/package.mjs --platform linux --arch x64
```

The renderer driver writes a new directory beneath `verification-output/`,
opens only its disposable fixture, checks browser features, cancels, reopens
and submits. It uses one pending CLI call per interaction; it does not invoke
a model callback or poll a person. `evidence.json` records results and artifact
hashes. Its 45-second process bound is a test timeout, not a product wait limit.
It clears the child executable search path, preserves worker diagnostics on
failure, and does not turn unavailable WebGL into a passing rendering check.
Its default scope requires all 21 checks. `--allow-unavailable-webgl` permits
only a null WebGL context and records it explicitly as unavailable; an
incorrect render still fails. `--expect-startup-failure` instead requires one
saved `renderer_startup_failed` result and the indicated sandbox diagnostic
log. It fails if startup succeeds, prompting review of the platform record.

`.github/workflows/verify.yml` runs the core, Node and wiki checks before the
actual Electron fixture on hosted macOS, Windows, Ubuntu 22.04 and Ubuntu
24.04. The renderer sandbox stays enabled. Linux uses the runner's existing
Xvfb; diagnostic namespace checks do not change its security policy. See the
[hosted report](runs/2026-10-22-hosted.md) for failures and tested revisions.
After the user's [Linux decision](../planning.md#linux-distribution-decision),
Ubuntu 22.04 uses the optional-WebGL scope; Ubuntu 24.04 tests the unsupported
startup error, not successful rendering. CI job names state those scopes.

The portable core test fixture replaces shell-only helpers. Desktop execution
uses the [handoff recipe](desktop.md). The packaging command can download a
pinned official target runtime during development; a completed package does
not need a runtime download, Node, npm or Go on the user's machine.

## Installation and updates

Follow the [managed installation recipe](../features/installation.md#build-and-verification-recipe)
to build two versioned archives and run `scripts/install-smoke.mjs`. It serves
one actual archive/checksum from a local fixture server, executes the platform
bootstrap in a fresh directory with spaces in its name, holds a real Electron
session through an update, checks its saved PNG, verifies runtime file sharing
and rolls back. Installed CLI children have an empty PATH. The driver emits
full evidence and saves it beneath its new `verification-output/install-smoke-*`
directory. It sends no model callback and changes no persistent PATH.

CI runs this after its renderer check on each target. Ubuntu 24.04 uses
`--core-only`; Ubuntu 22.04 uses its existing Xvfb. Neither test changes host
sandbox policy. These fixtures prove local release archives and bootstrap
downloads, not public GitHub assets, signing or clean consumer-machine setup.

## Experimental Wails renderer

The nested `renderer/go.mod` pins Wails v2.17.0. From the repository root:

```sh
go build -o verification-output/fungi-wails-core .
cd renderer
go build -tags production -trimpath -ldflags '-s -w' -o ../verification-output/fungi-wails ./wails
cd ..
node scripts/renderer-smoke.mjs verification-output/fungi-wails-core --helper verification-output/fungi-wails
node scripts/wails-probe.mjs verification-output/fungi-wails-core verification-output/fungi-wails
```

Use `.exe` output names on Windows. Linux builds use
`-tags production,webkit2_41`, GTK3/WebKitGTK 4.1 development packages and
Xvfb for the hosted fixture. Real desktop use needs compatible installed
runtime libraries; the CI prerequisite installation is not consumer setup
proof. The existing smoke fixture can explicitly permit unavailable WebGL,
but does not waive media/module failures. The host probe records direct and
Blob media separately without redefining the original expectation.

The first driver starts the actual CLI/owner with an empty child PATH,
cancels, reopens and submits a PNG with an independently checked SHA-256.
The second starts disposable sessions, saves/restores a zero-byte attachment,
invokes the native quit/close hook, and checks whether one independent session
inherits another's run-specific localStorage value. It removes only that key.
Native API names are inspected without reading the clipboard or opening URLs.
Both drivers have bounded execution, save `evidence.json` in unique local
directories, and send no harness callbacks.

For native UI, invoke the helper with `fungi open --exec ABSOLUTE_HELPER
--delivery manual -- ABSOLUTE_HTML_OR_JSON`, then follow the existing
form/file/close/reopen recipe. A minimal `.app` bundle was used to identify the
Mac prototype through accessibility tools; it does not change the binary.
Use `collect` once after observed completion. Do not treat synthetic answers
as user decisions. Exact builds, observations and outstanding scope are in
the [Wails report](runs/2026-10-25-wails.md).

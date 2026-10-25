---
title: Verification setup
updated: 2026-10-25
---
# Verification setup

Run from `/Users/rengwu/Desktop/Projects/fungi`. Early runs predate Git initialization; current runs record the Git base, dirty source hashes and binary SHA-256.

For the user's Windows checkout, follow the [Windows agent handoff](windows.md).
It includes the correct branch, source build, actual desktop checks and result
recording without relying on this Mac's paths or a sibling jstack checkout.

## Build and launch

Development tools used here are Go 1.26.4 and Node 22.19.0. The nested
`renderer/go.mod` pins Wails v2.17.0. Build on the target OS/CPU: macOS needs
Xcode command-line tools, Linux needs GTK3/WebKitGTK 4.1 development packages,
and Windows needs no C compiler. These are build prerequisites, not end-user
requirements. No npm installation or Electron download is used.

```sh
go vet ./...
go test -race ./... # Windows CI uses go test ./...
node --test
node scripts/package.mjs
# Use the exact complete package path printed by the command:
/absolute/package/fungi doctor
/absolute/package/fungi describe
```

Packaging creates a unique directory under `dist/`. Do not use a partial
build or guess the latest directory. Read its `BUILD.json` and verify both
CLI and helper SHA-256. Use the complete package; `go build .` alone does
not supply the window helper.

Use `FUNGI_STATE_DIR="$PWD/verification-output/sessions"` for local UI fixtures. Use a separate root for each harness batch. `open` returns only after readiness, or keeps one pending pipe in `wait` mode. Never poll `collect` while awaiting a person.

## UI driver and fixtures

For current macOS UI checks, select the exact packaged or installed
`renderer/Fungi.app` path through the connected desktop driver. Match its
window title and disposable session before typing. The following original
clipboard/PNG observations used Electron; retain their build-specific scope.

The generated PNG fixture is 64×64 RGB, 180 bytes, with SHA-256 `1bd7e798fbfd3bd20de881435c7bfa7811db8466e73714a8ab7c1af6b7460e9a`. Its source bytes are retained in evidence. The current reference form exposes accessible labels for text fields, the native file picker, and “Paste image for Reference image.”

The successful clipboard fixture used a separate form with one image field, keeping its paste target visible. Copying pixels from macOS Preview and pasting into that target yielded a 401-byte, 64×64 PNG with SHA-256 `83e160394a50c08746093b421e7bbf7865d6591aac9f9b2d1f0850bea504c2d9`. Both the original and clipboard-encoded bytes are retained. The crash fixture had one text field: after “Draft saved,” only its identified page-renderer PID was terminated. The restored window visibly contained the acknowledged value before it was closed.

Do not overwrite the user's edits while they are trying a demo. Preserve acknowledged drafts. A changing desktop focus or failed clipboard setup is a driver limitation until separately diagnosed; it is not a product pass.

## Historical harness probes

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

## Historical no-setup follow-up fixtures

The follow-up lives under `verification-output/no-setup-verification/`. `run-ui-probe.py` launches one disposable harness conversation with an explicit fixture command and a 240-second external test bound. Each form has a unique title, response marker and state directory; the connected desktop driver chooses the exact packaged Electron app path. The user-edited source demo is separate.

Kimi's successful follow-up reports 2.2.0. It uses `run_in_background: true` and `disable_timeout: true`; no authentication files are read. Claude's headless test deliberately supplies 600,000 ms and does not represent unlimited interactive waiting. Non-Codex child probes do not inherit this test coordinator's `CODEX_THREAD_ID`.

Do not name a Fungi form `opencode.json` in an OpenCode working directory: that name is reserved for its configuration and may be rewritten. The corrected fixture is `opencode-form.json`. The initial collision is recorded as a failed fixture attempt, not hidden or treated as a Fungi product failure.

For manual delivery, let the launch return and the harness end its turn, submit through the native window, then resume the exact recorded harness session with an explicit “done.” It must call the receipt's `collect` argv with `collectEnv` once. This tests manual fallback and cannot be counted as automatic completion.

## Observation and cleanup

Inspect `launch.json`, `draft.json`, `result.json`, artifact bytes/hashes and `delivery.json` separately. The OS file lock proves owner liveness without a PID guess. Preserve the result and evidence before removing disposable data.

Do not delete a user-modified demo draft. No persistent harness configuration or authentication was changed. Public signing/notarization, clean consumer-machine journeys, physical Windows/Linux Wails input, long sleep/wake cycles and abrupt power-loss behavior remain outside current evidence. Historical results below retain their original renderer/build scope.

## Skill authoring and idle probes

The simpler questionnaire lives at `skills/grill-form/assets/form.html`; pass its question JSON through `--input`. Native fixtures and receipts live under `verification-output/skill-proof/`. The `example.json` is explicitly a demonstration, so synthetic test answers are safe there. Never submit test answers into an actual interview.

```sh
node --test scripts/grill-form.test.mjs
go build -o verification-output/idle-helper ./scripts/idle-helper
/absolute/package/fungi open skills/grill-form/assets/form.html --input skills/grill-form/example.json --delivery manual
```

For native recovery: enter a choice/note, select the known reference image, wait for “Draft saved,” close, collect once, and reopen. The initial session must stay `closed`. To test abrupt owner loss, identify the exact `_run SESSION` process in this disposable root before killing it; never terminate by generic Electron/Fungi process name. Collect the interrupted session once, reopen it, and verify the acknowledged draft and image. The submitted result must preserve the image hash.

For the Linux core recipe, copy only root Go sources/module, `examples/helper`, `scripts/protocol-fixture` and any selected `scripts/idle-helper` fixture into an isolated directory. The [recorded Docker command](../evidence/2026-10-20-skill-proof/linux-command.json) pins the official Go image, uses an unprivileged user, disables networking and mounts that source read-only. Its disposable container removes itself. The passed run is Linux arm64 under OrbStack; it is not a Linux desktop verification.

The idle helper takes `--input` JSON with `seconds` (1–86400) and `marker`. It acknowledges a draft, announces readiness, sleeps locally once, then submits timings and a JSON artifact. The five-second smoke check must pass before launching a long run. For a real callback test, run with `--delivery auto` from the intended conversation, preserve the receipt, and end the agent turn. Its one completion notice supplies the exact executable/session/environment for one collection. For other harnesses, use their already documented completion route. Do not wrap the wait in repeated tool calls.

The completed two-hour probe is session `622f59b42e05840dd408622f23c8416a`, with state beneath `verification-output/skill-proof/idle-long-sessions`. Its original executable is in `dist/fungi-darwin-arm64-2026-10-20T09-41-25-905Z`. Its [result](../evidence/2026-10-20-skill-proof/idle-long-result.json) records one collection after the exact-origin notice, matching timing artifact bytes and exit of the helper/owner. Retain the evidence. For future runs, keep the launched package and helper available until completion; do not suspend the machine without settling timing with its user. Record both the helper's elapsed measurement and the UTC timestamp difference when sleep occurs.

The user manually completed sleep/wake fixture `1897358a56ceff72064a252004170716`, using the same package and `verification-output/skill-proof/sleep-wake-sessions`. Its [baseline](../evidence/2026-10-20-skill-proof/sleep-wake-baseline.json) records the initial note/image and preparation time. After its completion notice, one collection returned the preserved note/image and user's post-wake edit. A minimal [power-log excerpt](../evidence/2026-10-20-skill-proof/sleep-wake-power.log) independently confirms actual sleep/wake within the interval. The [result record](../evidence/2026-10-20-skill-proof/sleep-wake-result.json) retains the callback, data and byte validation. For future repetitions, obtain the user's timing, prepare a fresh fixture, let them sleep/wake and submit, then collect once; inspect only relevant `Sleep`, `Wake` and `DarkWake` event rows since the baseline, excluding scheduled `Wake Requests`. Do not select an answer or submit on the user's behalf.

## Platform foundation checks

```sh
node --test
go test -race ./... # omit -race for the Windows fixture
go vet ./...
node scripts/package.mjs
node scripts/renderer-smoke.mjs /absolute/package/fungi
```

The renderer driver opens a disposable page through ordinary `open`, checks
21 browser features, cancels, reopens and submits. It verifies the PNG hash,
clears the child PATH, preserves diagnostics and uses finite fixture bounds.
It does not call a model or poll a human. `--allow-unavailable-webgl` permits
only a null context; an incorrect render still fails. Linux's explicit
`--allow-unavailable-local-media` records the known direct-media failure in
`excluded`, never as passed. Neither exception is used on Mac/Windows.

`.github/workflows/verify.yml` builds native packages on macOS, Windows,
Ubuntu 22.04 and 24.04. It checks the core, JS bridge, wiki, native helper,
actual packaged renderer and installation/update path. Linux uses Xvfb.
No host security policy is changed. A package built for one Linux version
must also be tested on the other before claiming a single shared archive.

## Installation and updates

Use the [managed installation recipe](../features/installation.md#build-and-verification-recipe)
to build two versioned archives and run `scripts/install-smoke.mjs`. It serves
actual archives/checksums from a local fixture server, runs the real platform
bootstrap in a fresh path with spaces, holds a Wails session through update,
checks its saved PNG and rolls back. Installed CLI children have an empty
PATH. It sends no model callback and changes no persistent PATH. Both Linux
versions exercise the actual UI under Xvfb, not a core-only substitute.

On a Linux host with Docker, `node scripts/runtime-setup-smoke.mjs ARCHIVE
ubuntu-22.04` (or `ubuntu-24.04`) starts a disposable official container with
no WebKit libraries. It checks doctor and ordinary-open failures, declined
setup, an injected package-manager failure without activation, then real
opt-in official runtime installation and reuse. Only this test container
installs prerequisites. It makes no desktop/input claim. Windows CI with an
already installed WebView2 runtime does not prove the missing-runtime
Microsoft bootstrap path.

These fixtures verify local release archives and bootstrap downloads, not
public GitHub release assets, signing or a clean consumer desktop.

## Wails diagnostics

For standalone helper work, build `./wails` from `renderer/` with tags
`production` (macOS), `production,wv2runtime.error` (Windows), or
`production,webkit2_41` (Linux). Use `.exe` names on Windows.

```sh
node scripts/renderer-smoke.mjs CORE --helper HELPER
node scripts/wails-probe.mjs CORE HELPER
```

The host probe checks native close, zero-byte Blob recovery and whether a
second session inherits a run-specific localStorage value; it removes only
that key. Native API names are inspected without reading clipboard contents
or opening URLs. Targeted Linux GStreamer logs diagnose direct/Blob media
separately. Known limits are accepted, not changed into passing results.
Both drivers have finite bounds, unique evidence directories and no callbacks.

For visible input, use ordinary `fungi open PAGE --delivery manual`, then
follow the [desktop recipe](desktop.md). Collect once after an observed
close/submit/cancel. Never submit synthetic answers into a real interview.

## Offline discovery verification

`go test -run TestOfflineHelpAndDiscovery .` runs the embedded reference without
adjacent docs/renderer or PATH tools. For actual packages, use two unique release
versions, then `node scripts/install-smoke.mjs ARCHIVE_A ARCHIVE_B`. Its discovery
checks run through the stable launcher before update, after update and after
rollback, under a path with spaces, unrelated working directory and empty child
PATH. Each help output must match its installed Markdown, and each example's
absolute launch paths must select the active version.

For native reply verification, run the installed CLI's `describe examples`,
select `reply` and execute its argument array with `--delivery manual` and a
fresh `FUNGI_STATE_DIR` outside the install. Use the actual window to type
synthetic notes and Send; collect once from its receipt. Check exact result data,
empty artifacts and absence of `draft.json`. Also exercise Cancel, window close
and explicit reopen (which should be blank without a draft). Keep result JSON
and the CLI/helper hashes as evidence. Close all fixture windows when finished.

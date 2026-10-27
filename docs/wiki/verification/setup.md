---
title: Verification setup
updated: 2026-10-27
---
# Verification setup

Run from the Fungi repository root. Early runs predate Git initialization; current runs record the Git base, dirty source hashes and binary SHA-256.

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

## Scratch output lifecycle

The 2026-10-26 cleanup retired the old `dist/` packages, root executable,
Electron dependencies and `verification-output/` scratch trees. Completed
feedback and dated evidence remain in this wiki. Historical absolute paths in
reports identify their original runs; they are not current launch commands.
For a new run, build a fresh package and use the maintained drivers above.
There is no general dependency on a pre-existing root `./fungi` executable.

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

## Native helper idle probe

`scripts/idle-helper` exercises the native executable protocol without a web
window. Build it separately:

```sh
mkdir -p verification-output
go build -o verification-output/idle-helper ./scripts/idle-helper
```

Write an input file containing `{"seconds":5,"marker":"unique-test-marker"}`,
then run the helper with the complete package:

```sh
/absolute/package/fungi open --exec /absolute/verification-output/idle-helper --input /absolute/idle-input.json --delivery wait
```

The helper acknowledges a draft, waits on one local timer, then submits timing
data and a JSON artifact. Verify the marker, elapsed time and artifact hash.
The five-second test is bounded; use an appropriate native completion route
for longer waits. For an exact-caller callback, use `--delivery auto`, preserve
the receipt and end the turn. Collect once on its completion notice.

The [October 20 skill and idle proof](runs/2026-10-20-skill-proof.md) preserves
historical questionnaire, native-helper and sleep/wake results. The in-repo
questionnaire and jstack adapter have been removed; workflow integration now
lives in `jstack-fungi-test`. Use the current page fixture below for Wails
idle and sleep/wake checks.

## Platform foundation checks

### Current Wails idle and sleep/wake fixture

`scripts/fixtures/idle.html` uses only the public four-method page API. Pass
`{"marker":"unique-test-marker","seconds":5}` for a short smoke check with
`--delivery wait`. It saves a note and PNG, idles using one timer, submits, then
closes. Check the returned marker, input preservation, elapsed time and artifact
hash before scheduling a long run. A bounded driver may hold the five-second
test call; this does not establish host support for arbitrary human waits.

Use `seconds:7200` and `--delivery auto` for a two-hour callback run. Preserve
the receipt and end the turn; collect once on its exact-session notification.
The page has no polling loop, heartbeat or model calls. A `manual` receipt
requires the documented user cue instead. Do not retry an uncertain notification.

For user-controlled sleep/wake, use `seconds:0`. Confirm the window says its
draft and PNG are saved before handoff. The user sleeps/wakes the Mac, edits the
note and submits. Record the exact before/after data and PNG hash. Inspect only
sleep/wake power-log events within the test interval; without those observations,
do not claim actual sleep. Do not change power settings or sleep the user's
machine automatically. Wall-clock and performance-clock timings are both returned;
neither alone establishes that OS sleep occurred.

Preserve the package, fixture and state directory until every pending session
is collected. Capture build/source hashes and distinguish current Wails results
from the historical native-helper and Electron evidence above.

### Callback across an intervening chat message

For a callback across an intervening user turn, use the packaged reply example
and a unique test title/input. Launch once with `--delivery auto`, preserve the
receipt, and require `delivery.mode: callback`. Ask the user to type a note,
leave the form open and send another chat message. During the resulting agent
turn, ask them to submit while a bounded verification command runs. Record that
command's start/end and the completion notice's arrival. Collect once only on
the exact-session notice, then verify the note and result identity. No “done”
message, polling or second waiter should be needed. A host-managed `wait` route
needs its own test; this callback scenario does not qualify it.

If the human submits before the bounded command starts, record that timing
instead of claiming overlap. For a separate mid-tool check, launch the existing
idle page with a five-second timer and callback delivery, then immediately run
`go test -run '^TestWaitConnectionHasNoHumanTimeout$' -count=1 -v .` in the same
tool call. Remove caller routing from the Go test's environment. On the callback,
collect once and compare the page's finish time and original result-file time
with the command interval; verify the input marker and PNG hash. Report deferred
delivery after turn end separately from immediate interruption. The
[completed local run](runs/2026-10-27-interleaved-message.md) records both cases.

### Automated core and renderer checks

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

For window dimensions, run `node scripts/window-smoke.mjs /absolute/package/fungi`
on a desktop large enough for a 960 × 780 window plus native decorations.
It opens disposable HTML sessions with default, custom, single-axis and minimum
sizes, saves their measured viewports, cancels, and reopens them. It also checks
a width override while preserving height and draft. Each command has a 30-second
test bound and uses `--delivery wait`; no model callbacks or human polling occur.
The printed evidence path records results and package identity. Preserve it
before removing that run's scratch directory. The driver is verified on macOS;
Windows/Linux execution remains untested.

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
PATH. Each help output must preserve its installed manual's title and code
examples while formatting prose/tables for the terminal; each example's
absolute launch paths must select the active version.

For native reply verification, run the installed CLI's `describe examples`,
select `reply` and execute its argument array with `--delivery manual` and a
fresh `FUNGI_STATE_DIR` outside the install. Use the actual window to type
synthetic notes and Send; collect once from its receipt. Check exact result data,
empty artifacts and absence of `draft.json`. Also exercise Cancel, window close
and explicit reopen (which should be blank without a draft). Keep result JSON
and the CLI/helper hashes as evidence. Close all fixture windows when finished.


## Public alpha downloads

`node scripts/public-release-smoke.mjs v0.1.0-alpha.1` uses anonymous HTTPS to
fetch the published installer/checksum, installs to a temporary path with spaces,
runs doctor and public update, discovers examples, and invokes the real renderer
fixture. On Linux use `xvfb-run -a` when no physical display is present. The
runtime must already be installed. The driver changes no persistent PATH and
deletes its managed test installation; preserve its JSON and renderer evidence
before removing fixture sessions. Maintainer Node is not an end-user dependency.
The manual **Publish alpha** workflow runs this check for one selected version
and platform, skipping build/publication. See the
[first release record](runs/2026-10-26-release.md) for executed recipes and the
unexplained intermittent Mac media result. These checks do not establish signing,
clean-machine security prompts or physical input qualification.

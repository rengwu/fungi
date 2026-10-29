# Fungi

<img alt="Fungi icon" src="./docs/assets/icon/fungi-mac.svg" width="34%" align="right">

A headless Go session engine for agent workflows. Fungi launches external
helpers, saves drafts, JSON responses and file artifacts, and returns completion
to the calling agent without polling. Other projects own all user interfaces,
rendering, assets and helper runtimes.

The engine is one executable with an embedded offline manual. It needs no
webview, graphical desktop, Go, Node or npm at runtime.

<br clear="right">

```sh
fungi open --exec ./helper --input ./input.json --delivery auto
fungi help helper
fungi help examples
fungi describe
fungi doctor
```

Pass helper arguments after `--`. The helper exchanges newline-delimited JSON
with Fungi over stdin/stdout. It reports readiness, optionally saves drafts,
then submits a result and waits for acknowledgment. Fungi saves artifact bytes
before acknowledging success. A helper can implement any interface in its own
project; the engine treats its input and response schemas as opaque JSON.

Read the launch receipt: automatic delivery may fall back to manual. For a
callback, end the turn and collect once on notification. For manual delivery,
ask the user to say “done,” then collect once. Use `--delivery wait` only when
the host supports a pending call or automatic background completion. Never
poll `collect`. Outcomes are `submitted`, `cancelled`, `closed`, and `failed`.

## Engine-only migration

This is the **0.2 development line**, using **protocol 2** and **package format 3**.
Previously published `v0.1.0-alpha.*` packages include the old UI. No engine
release is published by these source changes.

Bundled Wails/webviews, HTML hosting, the JavaScript page bridge, JSON forms,
UI examples and window flags have been removed. All UI implementations are
left to other projects. `fungi help migration` explains integration changes.
Existing protocol-1 sessions require their original executable and state
directory; the engine rejects them without modifying their records.

## Build and use

Developers need Go 1.26. A standalone engine build is sufficient:

```sh
go build -o fungi .
./fungi help
./fungi doctor
```

For a complete headless round trip from the source checkout:

```sh
go build -o /tmp/fungi-helper ./examples/helper
./fungi open --exec /tmp/fungi-helper --delivery wait
```

The embedded examples manual also includes a Python standard-library helper.
The selected example/helper supplies its runtime; the engine supplies none.

Node 22+ is used for packaging and verification only. Build a versioned archive:

```sh
node scripts/package.mjs --release v0.2.0-test.1
```

The printed archive and adjacent SHA-256 file identify the package. Extract it
and run the extracted engine directly, or create a separate managed installation:

```sh
./fungi install --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256 --dir "$HOME/.fungi-engine"
export PATH="$HOME/.fungi-engine/bin:$PATH"
fungi doctor
```

Add that PATH entry to shell configuration to retain it. A separate directory
keeps old UI installations available during migration. Managed engine installs
support `update`, `rollback` and `prune`. See [installation](docs/manual/installation.md).
`doctor` checks atomic session storage and locking without starting a helper.

## Verify

```sh
go vet ./...
go test -race ./...
node scripts/engine-smoke.mjs /absolute/package/fungi
node scripts/install-smoke.mjs ARCHIVE_A ARCHIVE_B
node docs/wiki/_build.mjs --check
```

Omit `-race` for the Windows fixture. The release workflow targets Mac arm64,
Windows x64 and Linux x64; the packager can cross-build using `--platform` and
`--arch`. Build success alone does not qualify execution on another platform.
The [verification record](docs/wiki/verification/runs/2026-10-28-headless.md)
records exact tested behavior and remaining limits.

Session state defaults to the OS user configuration directory, or
`FUNGI_STATE_DIR`. Keep it outside managed installs. Sessions remain until
explicitly removed. Helpers are trusted native processes; no application
sandbox or remote UI transport is supplied.

[Offline manual](docs/manual/start.md) · [Helper protocol](docs/manual/helper.md) ·
[Completion and recovery](docs/manual/delivery.md) · [Publishing](docs/RELEASING.md)

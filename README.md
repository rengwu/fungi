# Fungi

A small Go session engine with a Wails window for skill-authored human interaction.
Skills bring ordinary HTML/CSS/JavaScript or a native helper; Fungi owns the
session, saves responses and returns them to the calling agent without polling.
This is a development preview, not a published or signed release.

Each platform package contains the CLI and a native window helper. It uses
WKWebView on macOS, WebView2 on Windows and WebKitGTK 4.1 on Linux. End users
need no Go, Node, npm or bundled Electron. Keep the package directory together.
Run `fungi doctor` to check the runtime without opening a window.

```sh
./fungi help                   # short capability overview and task navigation
./fungi help examples          # complete input → interface → response flows
./fungi describe examples      # absolute installed paths and launch argv
./fungi open ./ui/index.html --input ./input.json --delivery auto
```

The [offline manual](docs/manual/start.md) is embedded in `fungi help TOPIC`
and shipped as Markdown. It covers the [four-method page API](docs/manual/html.md),
[completion and results](docs/manual/delivery.md), [CLI](docs/manual/cli.md),
[native helpers](docs/manual/helper.md), [experimental forms](docs/manual/form.md)
and [installation](docs/manual/installation.md). No separate Fungi skill is needed.

Use ordinary local HTML/CSS/JS and libraries for any interaction the webview
can render. Skills own the UI and meaning of its data. Draft saving is optional;
submitted JSON and file artifacts are durable independently. The small
[reply](examples/reply/index.html) and [canvas comparison](examples/compare/README.md)
examples are starting points, not a list of allowed workflows.

Inspect the launch receipt: automatic delivery can fall back to manual.
End the turn and collect once on notification or the user's “done.” Use
`--delivery wait` only if the host keeps one call pending or automatically
returns background completion. **Never poll `collect`.**

## Skill example

The repository’s [grill-form skill](skills/grill-form/SKILL.md) adapts jstack's simpler committed questionnaire. It uses ordinary HTML/assets and the existing bridge, with no generator or registration:

```sh
./fungi open skills/grill-form/assets/form.html --input skills/grill-form/example.json
```

It supports choices, plain notes, offline Markdown/diagrams, acknowledged drafts and one optional image. The current jstack review/editor rework is intentionally outside this example. UI libraries live with the skill; the engine has no questionnaire-specific API.

The separate [rich jstack demo](examples/jstack-grill-form/README.md) adapts the newer generated form, including visual choices, prototype controls, region comments and inline images. The simpler example skill remains available.

## Install and update

The shell/PowerShell installers download one OS/CPU package and verify its
checksum. If the system webview is missing, installation stops with instructions.
Explicitly opt into setup with `install.sh --setup-runtime` or
`install.ps1 -SetupRuntime`. Linux automatic setup currently supports Ubuntu
22.04/24.04; Windows uses Microsoft's official WebView2 bootstrapper. Normal
`open`, `doctor` and `update` commands never install system runtimes.

Managed installs support `fungi update`, `fungi rollback` and `fungi prune`.
Public installer URLs and releases are not published yet. For local archives:

```sh
./fungi install --archive PATH_TO_ARCHIVE --sha256 EXPECTED_SHA256
```

Packages use format 2 (Wails). Older Electron preview packages are no longer
accepted for installation or rollback. See [installation](docs/wiki/features/installation.md).

## Develop

Build on the target OS/CPU with Go 1.26, Node 22+ and native build tools.
Linux needs the GTK3/WebKitGTK 4.1 development packages. No npm install is needed.

```sh
go test -race ./...
node --test
node scripts/package.mjs
node scripts/package.mjs --release v0.1.0-test.1
```

Packaging prints the output directory or archive path. Use its `fungi` command,
not a standalone core build lacking its window helper. For a versioned build:

```sh
dist/release-v0.1.0-test.1/fungi-darwin-arm64/fungi doctor
node scripts/renderer-smoke.mjs dist/release-v0.1.0-test.1/fungi-darwin-arm64/fungi
```

Adjust the target path for Windows/Linux. Linux checks explicitly record the
accepted direct-media limitation with `--allow-unavailable-local-media`.
`FUNGI_STATE_DIR` selects isolated session storage; `FUNGI_RENDERER_DIR` is a
local development override for the packaged renderer directory. Sessions remain
until explicitly removed. No service remains after a completed session.

The [wiki](docs/wiki/index.md) records exact build verification and remaining
limits. Hosted execution is not clean-machine, signing or indefinite-idle proof.

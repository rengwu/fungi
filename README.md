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
./fungi describe
./fungi describe html
./fungi open examples/feedback.json
./fungi open examples/review.html --delivery wait
./fungi collect EXACT_SESSION_ID
./fungi reopen EXACT_SESSION_ID
```

`auto` targets an inherited exact Codex conversation when `codex queue` is available; otherwise it returns an explicit manual fallback. `wait` holds **one** command until completion, with no Fungi human timeout. Use it only through a harness that keeps the call pending or automatically delivers background completion. A tool session ID alone is not an automatic return mechanism. **Never poll `collect`.** In manual mode, the user says “done,” then the agent executes the receipt's `collect` argument array once with its `collectEnv` environment.

Fungi requires no harness plugin or configuration step. `describe delivery` explains the native routes and limits. Use manual mode when the installed harness cannot sustain the intended human wait; a successful short tool call does not make a bounded route AFK-safe.

Results are durable JSON envelopes with `submitted`, `cancelled`, `closed` or `failed` outcomes. Attachments are saved files referenced by key and absolute path. Collection does not delete anything. Reopening creates a new session from the last acknowledged draft; it never changes the original result.

```html
<button id="send">Send</button>
<script>
(async () => {
  const { input, draft } = await fungi.ready();
  document.querySelector('#send').onclick = () =>
    fungi.submit({ notes: draft?.notes ?? 'Looks good' });
})();
</script>
```

The bridge has four methods: `ready()`, `saveDraft(data, attachments?)`, `submit(data, attachments?)`, and `cancel()`. An attachment is `{key, name, blob}`. `ready()` returns restored attachments in that shape. Save drafts explicitly in custom HTML. Mutations run in call order, including Blob conversion; failed saves can be retried. Only acknowledged drafts are recoverable after close or interruption.

Skills can use local HTML/CSS/JS, modules, custom elements, Canvas/SVG/WebGL,
workers and WebAssembly. Pages are trusted local application code and can access
Wails native APIs. Fungi keeps drafts/results separate; browser storage isolation
is best effort. Use `saveDraft` for durable state and `sessionStorage` for
window-local state. Linux direct media URLs have a known WebKitGTK limitation;
Blob loading can help with small files, but no transparent workaround is promised.
See the [authoring guide](docs/wiki/authoring.md) and
[comparison example](examples/compare/README.md).

The JSON renderer supports flat `text`, `textarea`, `range` and `image` fields. Its schema is experimental. See the examples and `describe form` / `describe helper`.

## Skill example

The bundled [grill-form skill](skills/grill-form/SKILL.md) adapts jstack's simpler committed questionnaire. It uses ordinary HTML/assets and the existing bridge, with no generator or registration:

```sh
./fungi open skills/grill-form/assets/form.html --input skills/grill-form/example.json
```

It supports choices, plain notes, offline Markdown/diagrams, acknowledged drafts and one optional image. The current jstack review/editor rework is intentionally outside this example. UI libraries live with the skill; the engine has no questionnaire-specific API.

The separate [rich jstack demo](examples/jstack-grill-form/README.md) adapts the newer generated form, including visual choices, prototype controls, region comments and inline images. The simpler bundled skill remains available.

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

Existing Electron preview installs need the new bootstrap once to migrate to
package format 2; their older updater rejects the new format. The managed layout,
legacy rollback support and saved sessions are retained. See
[installation](docs/wiki/features/installation.md).

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

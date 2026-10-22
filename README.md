# Fungi

A small Go session engine with an Electron window for skill-authored human interaction. This is a development preview. The [compatibility record](docs/wiki/verification/runs/2026-10-20-no-setup.md) covers tested native return paths and manual fallback. The [skill proof](docs/wiki/verification/runs/2026-10-20-skill-proof.md) adds one macOS sleep/wake cycle and a two-hour helper wait with automatic Codex return. The [foundation run](docs/wiki/verification/runs/2026-10-22-platform.md) verifies broader rendering and recovery. [Hosted platform checks](docs/wiki/verification/runs/2026-10-22-hosted.md) exercise Windows and Linux and record current failures. Arbitrary-duration reliability and clean-machine distribution remain unverified.

Run a local HTML file or a flat JSON form. A native executable can implement the same lifecycle without using Electron. No registration, plugin SDK or frontend framework is required.

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

Skills can use local HTML/CSS/JS, modules, custom elements, Canvas/SVG/WebGL, workers, WebAssembly and local streaming media. Inline `srcdoc` previews work; frame document requests, outbound network, popups, navigation and OS permission requests are denied. Skills should sandbox their previews; only the main page can call the bridge. Trusted native helpers have ordinary OS access. See the compact [authoring guide](docs/wiki/authoring.md) and the [comparison example](examples/compare/README.md), which exports a custom canvas as a PNG without engine changes.

The JSON renderer supports flat `text`, `textarea`, `range` and `image` fields. Its schema is experimental. See the examples and `describe form` / `describe helper`.

## Skill example

The bundled [grill-form skill](skills/grill-form/SKILL.md) adapts jstack's simpler committed questionnaire. It uses ordinary HTML/assets and the existing bridge, with no generator or registration:

```sh
./fungi open skills/grill-form/assets/form.html --input skills/grill-form/example.json
```

It supports choices, plain notes, offline Markdown/diagrams, acknowledged drafts and one optional image. The current jstack review/editor rework is intentionally outside this example. UI libraries live with the skill; the engine has no questionnaire-specific API.

The separate [rich jstack demo](examples/jstack-grill-form/README.md) adapts the newer generated form, including visual choices, prototype controls, region comments and inline images. The simpler bundled skill remains available.

## Develop

Build dependencies are Go 1.26 and Node 22.12+; they are not end-user runtime dependencies.

```sh
npm ci --prefix renderer
node renderer/node_modules/electron/install.js
go build -o fungi .
go test -race . ./examples/helper
node --test scripts/*.test.mjs
node scripts/renderer-smoke.mjs
node scripts/package.mjs
```

The pinned Electron 44 package uses an explicit runtime download command. Packaging copies that runtime into a new `dist/fungi-PLATFORM-ARCH-…/` directory beside the executable and renderer code. Copy the whole directory. Cross-package with `node scripts/package.mjs --platform win32 --arch x64` (or `linux` / `darwin`, `x64` / `arm64`); this uses Electron's build-time downloader and pinned checksums. Cross-packaging does not prove that the target runs. Public signing/notarization and clean-machine Gatekeeper verification remain outstanding; these are development packages. The [desktop verification recipe](docs/wiki/verification/desktop.md) needs no Go, Node or npm on the target.

Set `FUNGI_STATE_DIR` to isolate verification sessions. Otherwise data lives beneath the OS config directory in `fungi/sessions`. Sessions remain until explicitly removed. No background service runs after a completed session. `FUNGI_RENDERER_DIR` selects a renderer directory during development.

The source and current verification record live in [the project wiki](docs/wiki/index.md). The [contract](docs/wiki/contract.md) distinguishes intended behavior from verified coverage.

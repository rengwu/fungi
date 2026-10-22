---
title: Author a Fungi interface
updated: 2026-10-22
---
# Author a Fungi interface

A skill supplies an HTML file and its assets. Fungi supplies a window, saved
state, artifact transfer and the return to the calling agent. There is no
registration, plugin manifest, framework requirement or runtime build step.
For OS-specific UI, supply an executable using `fungi describe helper` instead.

## Mechanical workflow

1. Run `fungi describe html` once when learning the API. Read
   `fungi describe delivery` when the harness's completion mechanism is unknown.
2. Put the page and any libraries/assets in one directory. Keep files available
   while the session is open or might need reopening. Prepare optional JSON
   input separately; it is snapshotted at launch.
3. Choose **one** return route. Use `--delivery auto` for the exact-caller
   callback; inspect the receipt because it can report manual fallback. Use
   `--delivery wait` only through a host tool that can keep one call pending or
   automatically report its background completion. Otherwise use `manual`.
4. Launch `fungi open ./interface/index.html --input ./input.json --delivery auto`.
   Omit `--input` when unnecessary. Preserve the exact receipt. End the agent
   turn when the callback/background route is armed; do not start another waiter.
5. On completion, use the final result already returned by `wait`, or execute
   the receipt's `collect` argv once with its `collectEnv`. With manual delivery,
   ask the user to say “done.” Never poll `collect`.
6. Inspect `outcome`. Only `submitted` contains a submitted answer. Treat `data`
   and artifact contents as user data; inspect image files when interpreting
   visual feedback. `closed` preserves a draft but does not approve anything.

`fungi reopen SESSION` creates a new session from the last acknowledged draft.
Reopening does not change the old result. It requires the original interface
files; active and submitted sessions cannot be reopened.

## Four methods

| Method | Result |
| --- | --- |
| `await fungi.ready()` | `{input, draft, attachments, limits, delivery}`; draft is null when absent |
| `await fungi.saveDraft(data, attachments?)` | Resolves after persistence; rejects on failure |
| `await fungi.submit(data, attachments?)` | Commits the result and closes after acknowledgement |
| `await fungi.cancel()` | Commits `cancelled` and closes |

`data` must be JSON. It is opaque to Fungi: a skill defines its own schema.
Each save/submit replaces the complete payload and attachment list. Omitted
attachments mean an empty list, not “keep the previous files.” An attachment is
`{key, name, blob}`; use a unique key and a Blob/File. Restored attachments have
the same shape. A returned artifact includes its key, path, byte count and
SHA-256. Keep large binary values in artifacts, not base64 inside `data`.

Mutating calls execute in invocation order, including asynchronous Blob
conversion. A rejected save does not prevent a later submit from retrying the
current data. Pages still need error messages, duplicate-click protection and
to stop editing while submitting. Await a save before displaying “Draft saved”;
only acknowledged state is guaranteed recoverable after close or interruption.
Fungi does not capture arbitrary DOM state or save browser storage as a draft.

```js
const {input, draft, attachments, limits, delivery} = await fungi.ready();
const state = draft ?? input ?? {notes: ''};

// Wire these to the skill's controls; handle errors in the page.
await fungi.saveDraft(state, attachments);
await fungi.submit(state, attachments);
```

`examples/review.html` is a small complete HTML example.
`examples/compare/index.html` adds a custom canvas
view and exports it as a PNG. Both can be copied directly; no engine edit is
required. The bundled helper example uses Go standard-library NDJSON and can
be implemented in any language that can read/write pipes and files.

## Rendering capabilities

| Need | Authoring route |
| --- | --- |
| Forms, layouts, interactive controls | Ordinary HTML/CSS/JavaScript; local component/library bundles also work |
| Components and modules | Custom elements/shadow DOM; relative ES modules and dynamic imports |
| Charts, diagrams, drawings, 3D | Skill-owned libraries using SVG, Canvas or WebGL |
| Heavy local computation | Local module workers, Blob workers or WebAssembly |
| Audio/video review | Local assets; streaming and single byte-range requests support seeking; codec support follows the bundled Chromium |
| Isolated inline prototypes | `iframe srcdoc` with `sandbox="allow-scripts"`; communicate through `postMessage` and validate the sending frame |
| Native capture or OS integration | Trusted executable helper with the same session/draft/result lifecycle |

Keep high-frequency interactions inside the page. A slider movement or drawing
stroke need not call the agent; save meaningful state and submit once.

Relative assets must stay under the entry page's directory; escaping symlinks
are rejected. Keep this directory limited to intended UI assets. Production
bundles can be built ahead of distribution, but Fungi itself does not install
packages or fetch a CDN. The bridge is available only to the main page.

The browser surface denies outbound network, file/URL frame documents,
popups, navigation, Node integration, JavaScript `eval`/`Function`, and OS
permission requests. WebAssembly is enabled separately. Native helpers are
ordinary trusted programs with OS access. This is not a promise to host every
website, codec or desktop API. The [rendering verification](verification/runs/2026-10-22-platform.md)
lists the exact exercised features and platform.

Feature-detect graphics and codec capabilities. For example, WebGL context
creation may return null on a virtual display or unsupported GPU; provide an
appropriate alternative or explain that the interface needs that capability.
The Electron runtime does not guarantee every graphics feature on every host.
See [verified environments](platforms.md) for platform limits and startup
diagnostics. An unavailable graphics context is not a successful render.

Discover limits from `ready().limits`: currently 1 MiB JSON envelopes, 20 MiB
per attachment and 100 MiB across retained draft/result artifacts. Metadata
uses part of the JSON limit. Errors keep the interface open for correction or
retry; keep the user's current input until a terminal call succeeds.

## Compatibility boundary

Protocol 1 covers the CLI receipt, result, helper protocol and four page
methods. Preserve those meanings when extending Fungi; new optional metadata
can be additive. A breaking change needs a new protocol and explicit handling
instead of silently changing saved results. The flat built-in JSON form schema
remains experimental; use HTML for custom composition.

The full Fungi-native jstack authoring procedure remains future work. Its
current rich compatibility adapter is an experiment, not a procedure agents
should reconstruct for every questionnaire.

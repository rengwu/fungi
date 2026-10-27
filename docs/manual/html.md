# HTML page API

Put the HTML page and local assets in one directory; input JSON may live
elsewhere. No manifest, registration or framework adapter is required.

```sh
fungi open ./ui/index.html --input ./input.json --delivery auto
```

Add `--width 1100 --height 800` to choose the starting window size in logical
pixels. Each dimension is optional (defaults: 820 × 780; minimum: 420).
`reopen` reuses the launch size unless overridden. The OS may adjust sizes to
fit the display; use responsive page layouts. Full flag rules: `fungi help cli`.

## Four methods on `window.fungi`

| Awaitable method | Behavior |
| --- | --- |
| `ready()` | Returns `{input, draft, attachments, limits, delivery}`. Absent input/draft are `null`. Restored attachments are `{key,name,blob}`. Repeated calls reuse readiness. |
| `saveDraft(data, attachments?)` | Optional. Resolves after saving; ignore its return value. |
| `submit(data, attachments?)` | Saves a `submitted` result, resolves with it, then closes the window. |
| `cancel()` | Saves a `cancelled` result, resolves with it, then closes. |

`data` is JSON-serializable; the skill defines its schema. Each save or submit
replaces all data and attachments. Omitted attachments remove previous files
from the new snapshot. Each attachment is `{key,name,blob}`: a unique nonempty
key (maximum 200 UTF-8 bytes), nonempty name (255 bytes), and Blob/File.
The Blob's type supplies `mediaType`.

```js
const {input, draft, attachments, limits, delivery} = await fungi.ready();
const state = draft ?? input ?? {notes: ''};
// Use your own controls. Draft saving is optional:
await fungi.saveDraft(state, attachments);
// At submission, send the current complete state:
await fungi.submit(state, attachments);
```

Mutations execute in invocation order, including Blob encoding; a rejected
save does not block later calls. Await acknowledgment before showing success.
On failure, preserve input, show the error and allow retry. Prevent duplicate
submissions and editing while submitting. For complete error handling, see the
`reply` page from `fungi describe examples`.

## Images and other files

```js
const blob = await new Promise(resolve => canvas.toBlob(resolve, 'image/png'));
if (!blob) throw new Error('Could not export the canvas');
await fungi.submit({selection: 'after'}, [
  {key: 'comparison', name: 'comparison.png', blob}
]);
```

Use `limits` from `await fungi.ready()`: `jsonBytes` (currently 1 MiB), `artifactBytes`
(20 MiB per file), `totalArtifactBytes` (100 MiB retained per session).
Allow for JSON envelope metadata. Use attachments for binary data; repeated
saves of new bytes count toward the retained-artifact limit.

Custom pages do not autosave; Fungi does not capture DOM state or browser storage.
For recovery, save state and attachments and restore them through `ready()`.
Only acknowledged drafts survive closing. Drafts are optional; final results
persist independently.

## Assets and rendering

Assets, including symlink targets, must stay beneath the entry page's directory.
Supported features include local modules/dynamic imports, CSS, fonts, custom
elements, Canvas/SVG/WebGL, workers and WebAssembly. Bundle libraries with the
skill and prebuild frameworks as needed. Fungi provides no CDN, remote fetch
or package installation.

Only the main page receives the bridge. Isolate untrusted previews in
`iframe srcdoc` with `sandbox="allow-scripts"`; validate the sending frame for
`postMessage`. Skill pages are trusted code with Wails native API access;
asset confinement is not an application sandbox.

Browser storage isolation is best effort. Use `sessionStorage` for window-local
state and drafts for durable state. Detect graphics/codec support; a missing
WebGL context is a rendering failure. Linux WebKitGTK cannot load direct `wails:`
media URLs. A Blob may work for small files but loads all bytes and guarantees
neither codec support nor audible playback. Runtime baselines:
`fungi help installation`.

## Compatibility

Protocol 1 preserves the documented four page methods, CLI session commands,
receipts, results and helper messages: their names, types and meanings.
Optional fields and new capabilities may be added. Ignore unknown metadata;
check CLI/helper `protocol` and result `kind`/`outcome` before using data.
Breaking these contracts requires a new protocol. Read advertised limits rather
than hardcoding them; error messages and help prose are not machine interfaces.

Skills own their input/result schemas, UI and libraries. Wails native APIs,
internal commands and storage layout are implementation details, outside this
compatibility promise. Web APIs remain subject to the system webview.
The built-in JSON form schema remains experimental.

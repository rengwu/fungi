# Wails renderer experiment

This is a small, separate Go renderer using **Wails v2.17.0** and the existing
Fungi native-helper protocol. It evaluates system webviews without replacing
the released/default Electron renderer or changing the skill page API.

From the repository root, with Go and the target's native build tools:

```sh
go build -o verification-output/fungi-wails-core .
cd renderer
go build -tags production -trimpath -ldflags '-s -w' -o ../verification-output/fungi-wails ./wails
cd ..
verification-output/fungi-wails-core open --exec "$PWD/verification-output/fungi-wails" --delivery wait -- "$PWD/examples/review.html"
```

On Windows add `.exe` to both output names. On Linux use
`-tags production,webkit2_41` and install the GTK3/WebKitGTK 4.1 development
packages on the **build machine**. The renderer uses macOS WKWebView,
Windows WebView2 and Linux WebKitGTK; it does not include a browser engine.
Windows still needs the WebView2 runtime, and Linux still needs compatible
GTK/WebKitGTK runtime libraries and a desktop session. No automatic runtime
installation or browser fallback is implemented.

Pass a local HTML entry or the existing experimental flat-form JSON file
after `--`. Assets beside the entry load normally. `--input` goes before `--`.
The existing `ready`, `saveDraft`, `submit`, `cancel`, Blob attachments,
owner persistence and delivery modes are reused. For user waits, follow
Fungi's usual automatic-delivery instructions; `wait` above is a local demo,
not a recommendation to poll a model tool.

Skill pages are trusted local application code. Wails exposes its native
runtime and Go bindings in addition to `window.fungi`; this access is accepted.
Fungi keeps each session's process, drafts, results and artifacts separate.
Windows also uses a per-session WebView2 profile. Browser storage isolation is
best effort: macOS/Linux can share `localStorage`. Use `saveDraft` for recoverable
state and `sessionStorage` for window-local state. No framework fork or storage
wrapper is required.

Recovery is explicit. After a crash, collect once to retrieve any committed
result, or reopen the session's acknowledged draft. If a blank/crashed window
remains open, close it first; active sessions cannot be reopened. Unsaved edits
are not recoverable. There is no automatic restart or model polling.

Run the existing real-owner rendering/recovery fixture:

```sh
node scripts/renderer-smoke.mjs verification-output/fungi-wails-core --helper verification-output/fungi-wails
```

For native-host diagnostics, run `node scripts/wails-probe.mjs CORE HELPER`.
This checks native close, empty-Blob recovery, media loading and shared storage
using disposable pages and one run-specific localStorage key. API availability
is observed without reading the user's clipboard or opening external URLs.

Linux CI uses Xvfb and explicitly permits an unavailable WebGL context.
An available but incorrect WebGL result still fails. The manual **Wails proof**
GitHub workflow (also triggered by prototype changes on `foundation-hardening`)
builds and exercises each platform without building Electron.
Hosted runners with prerequisites are not clean-machine installation proof.
Mac/Windows pass all 21 checks. Both Linux fixtures fail direct local media
loading while the same WAV's metadata loads through a Blob URL. An experiment
adding `wails` to `WEBKIT_GST_ALLOWED_URI_PROTOCOLS` did not fix this and was
removed. WebKit's GStreamer web source registers HTTP, HTTPS and Blob schemes,
not Wails' custom scheme. See the wiki follow-up for evidence. There is no
transparent media workaround; Blob loading preloads the file and is not a
replacement for streaming large media. Codec availability also depends on the
installed engine.

These known limitations are accepted for now. Add workarounds only when they
stay simple; they do not block the next packaging/installation integration.
Keep the failing media observation visible rather than relabel it as working.

The prototype is intentionally outside release packaging/updating. See the
project wiki's renderer research and dated verification record for measured
compatibility, sizes and remaining migration work.

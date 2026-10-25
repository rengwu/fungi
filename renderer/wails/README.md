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

**Use trusted local pages only.** Wails exposes its own native runtime and
Go bindings in addition to `window.fungi`. The content policy and root-confined
asset handler do not establish Electron-equivalent native isolation, navigation
control or private storage. Windows uses a per-session WebView2 profile;
macOS/Linux tests confirmed storage shared across independent sessions. Do not ship
this as the default renderer on the strength of rendering tests alone.

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
The current Mac/Windows fixture passes all 21 checks. Both Linux fixtures fail
direct local media loading; a separate diagnostic loads the same WAV through
a Blob URL. This is a compatibility gap, not an automatic media workaround.

The prototype is intentionally outside release packaging/updating. See the
project wiki's renderer research and dated verification record for measured
compatibility, sizes and remaining migration work.

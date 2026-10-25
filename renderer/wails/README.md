# Wails renderer

Fungi's default window helper uses **Wails v2.17.0** and the existing native
helper protocol. The Go CLI owns persistence, delivery and lifecycle. Skills
use ordinary HTML/CSS/JS and `ready`, `saveDraft`, `submit`, `cancel`.

Build the complete native package from the repository root:

```sh
node scripts/package.mjs
```

Go, Node and native build tools are development prerequisites only. macOS
builds use Xcode command-line tools; Linux uses GTK3/WebKitGTK 4.1 development
packages; Windows builds need no C compiler. Build on the target OS/CPU.
The package contains the CLI, helper and examples; no browser engine or Node.

```sh
/absolute/package/fungi doctor
/absolute/package/fungi open /absolute/page.html --delivery manual
node scripts/renderer-smoke.mjs /absolute/package/fungi
```

The helper reads the source and session launch from stdin. `--check-runtime`
reports availability without opening a window. macOS uses WKWebView, Windows
uses WebView2, Linux uses WebKitGTK 4.1. Ordinary launches never install a
runtime. The bootstrap offers explicit `--setup-runtime` / `-SetupRuntime`;
see [installation](../../docs/wiki/features/installation.md).

For standalone helper development, build `./wails` from `renderer/` with tags
`production` on macOS, `production,wv2runtime.error` on Windows, or
`production,webkit2_41` on Linux. Use `fungi open --exec ABSOLUTE_HELPER
--delivery manual -- ABSOLUTE_PAGE` or the renderer smoke driver's `--helper`
option. Keep user waits on the normal documented delivery route; never poll a
person from the model.

Skill pages are trusted local application code. Wails exposes its native
runtime and Go bindings in addition to `window.fungi`; this access is accepted.
Fungi separates session processes, drafts, results and artifacts. Windows
also uses a per-session WebView2 profile. macOS/Linux can share localStorage;
use `saveDraft` for recovery and `sessionStorage` for window-local state.

Recovery is explicit: collect once after completion, or reopen the last
acknowledged draft. Close a surviving blank/crashed window first. Unsaved
edits are not recoverable. There is no automatic restart or heartbeat.

Linux direct local media loading is an accepted limitation: WebKit's GStreamer
web source does not handle the `wails:` scheme. Blob metadata loading worked
in the diagnostic fixture, but audible playback and large-file streaming
were not established. No transparent workaround or framework fork is used.
Feature-detect WebGL and codecs. Linux CI records unavailable WebGL and failed
direct media separately; it does not call them passing rendering checks.

`node scripts/wails-probe.mjs CORE HELPER` diagnoses native close, empty-Blob
recovery, media and shared storage. It uses disposable pages and a run-specific
localStorage key. Native API names are inspected without reading the user's
clipboard or opening external URLs. See the [verification setup](../../docs/wiki/verification/setup.md)
and [platform evidence](../../docs/wiki/platforms.md).

---
title: Reconsidering the bundled renderer
updated: 2026-10-25
---
# Reconsidering the bundled renderer

The user reopened the Electron decision because the bundled runtime is a
significant cost. This is research and a proposed experiment, not an accepted
replacement or a verified compatibility claim. Existing Electron proofs remain
valid for their recorded builds. [Planning](../planning.md)

## What can stay

The current Go owner launches the renderer as a child and exchanges the same
JSON protocol used by native helpers. Go owns authoritative drafts, answers,
artifacts and delivery. Electron owns the web window, local asset serving,
page bridge, temporary attachment transfer and native lifecycle observation.
Replacing that child does not require redesigning the session protocol or
asking skill authors to adopt framework-specific APIs.

Preserve `fungi.ready`, `saveDraft`, `submit` and `cancel`. Skills continue to
provide ordinary HTML/CSS/JavaScript and assets. The experimental JSON form
can remain a page built on the same bridge. This boundary makes an isolated
renderer experiment possible; it does not make a new renderer compatible by
construction. [Current contract](../contract.md)

## Where the browser engine would come from

| Approach | Engine source | Main tradeoff |
| --- | --- | --- |
| System webview in a dedicated window | WKWebView on macOS, WebView2 on Windows, WebKitGTK on Linux | Small Fungi package and native lifecycle; runtime availability/version varies |
| Installed browser | A browser already installed on the machine | No bundled engine; browser profile, window ownership and lifecycle need separate handling |
| Another bundled engine | Fungi still distributes a web engine | Can eliminate Electron specifically, but retains engine distribution, updates and compatibility work |
| Native widget/drawing toolkit | Toolkit's own rendering | Does not preserve arbitrary web components and existing HTML skills without additional web rendering |

Tauri documents the same system-engine split. Its framework does not supply a
uniform web platform across operating systems. macOS WebKit is tied to the
system version; Linux WebKitGTK comes from the installed runtime packages.
[Tauri webview versions](https://v2.tauri.app/reference/webview-versions/)

Microsoft documents WebView2 as a separate runtime: Windows 11 includes it,
and most Windows 10 machines already have it, but applications should detect
availability and handle missing installations. Production WebView2 cannot
substitute the ordinary stable Edge browser. A bootstrap can install the
runtime, but that changes Fungi's no-extra-runtime-setup requirement and cannot
guarantee an offline first launch when the runtime is absent.
[Microsoft distribution guidance](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/distribution)

Linux is the largest distribution constraint. The webview project's production
requirements explicitly include GTK and WebKitGTK packages. Build-time headers
and compilers are distinct from these end-user runtime requirements. There is
no evidence that a system-webview host runs on every clean Linux desktop with
no setup. Fungi's current Electron renderer also has an unsupported Ubuntu
24.04 sandbox baseline; removing Electron does not by itself establish support.
[webview prerequisites](https://github.com/webview/webview#prerequisites),
[Fungi's measured platforms](../platforms.md)

## Maintained implementation candidates

The [GitHub snapshot](../attachments/research/system-renderers-2026-10-23.json)
records API observations on 2026-10-23. Stars measure interest, not reliability.
These candidates are hosts or libraries around engines, not interchangeable
engine implementations.

| Candidate | Observed interest and activity | Assessment for Fungi |
| --- | --- | --- |
| [Wails](https://github.com/wailsapp/wails) | 36,662 stars; substantive Windows focus and WebView2 fixes on Oct 20 | First Go-native candidate to probe. v2 is stable; v3 is beta. Recent v3 work does not prove v2 has those fixes |
| [Tauri/Wry](https://github.com/tauri-apps/wry) | Tauri 112,135 stars; Wry 5,024; Wry initialization/IPC fixes on Oct 20 | Credible alternative host; adds a Rust build component while retaining the same OS runtime requirements |
| [webview](https://github.com/webview/webview) | 14,304 stars; C/C++ fixes Oct 19 and 22 | Small API, but verify needed navigation, permission and lifecycle hooks rather than equating fewer APIs with less maintenance |
| [webview_go](https://github.com/webview/webview_go) | 458 stars; latest sampled commits Aug 2024 | Separate binding maintenance concern; upstream C/C++ activity does not establish binding freshness |
| [WebUI](https://github.com/webui-dev/webui) | 4,631 stars; implementation fixes Oct 22 | Useful installed-browser design reference; separate lifecycle/security compatibility proof needed |
| [Neutralinojs](https://github.com/neutralinojs/neutralinojs) | 8,671 stars; v6.10.0 changelog Oct 6 | Another system-webview host; framework substitution does not remove platform runtime constraints |

Wails supports arbitrary frontend technologies and native close hooks. It can
be evaluated as the renderer child while Go's owner continues using the
existing process protocol. Authors should see only the Fungi bridge, not
generated Wails application bindings or project scaffolding. The first probe
should explicitly select a stable version; beta-only requirements must be
identified before making it the default.
[Wails version status](https://github.com/wailsapp/wails#readme),
[window and asset options](https://wails.io/docs/reference/options/)

## Alternatives covering all three desktops, checked 2026-10-25

The user reiterates that macOS, Linux and Windows are mandatory and asks for
alternatives to Wry. The [refreshed snapshot](../attachments/research/renderer-candidates-2026-10-25.json)
records Wails 36,689 stars, webview 14,312, Neutralinojs 8,671, WebUI 4,633
and wxWidgets 7,286, with recent commits sampled separately from push dates.
Wails now has observed v2 dependency maintenance dated October 25; its v3
remains beta. These observations do not qualify a Fungi renderer.

The leading alternatives remain Wails for Go integration and webview for a
small C/C++ host. Neutralinojs is another credible option: it supplies a
prebuilt host, local web server and WebSocket bridge, and supports child-process
integration. It uses system browser libraries and documents packages for all
three desktops. Fungi would still need an adapter for its existing protocol;
authors should not acquire a second required API.
[Neutralinojs architecture](https://github.com/neutralinojs/neutralinojs),
[distribution](https://neutralino.js.org/docs/distribution/overview/)

wxWidgets adds a mature C++ option through `wxWebView`, with WebKit on macOS
and GTK/Linux and a WebView2 backend on Windows. For Fungi it is a larger native
toolkit integration than a thin webview host. Its 3.2.11 interface documentation
also records backend-specific custom-scheme and scripting differences, so
desktop support alone is insufficient for our asset/bridge contract. Select
the modern Windows backend explicitly; the legacy IE backend is unsuitable.
[wxWebView interface and backend documentation](https://github.com/wxWidgets/wxWidgets/blob/v3.2.11/interface/wx/webview.h)

WebUI remains conditional on accepting an installed-browser interface and
proving its lifecycle. Pywebview also supports all three desktops, but adds
Python and platform bindings to package or require; it is less attractive for
Fungi's existing Go core and minimal distribution goal.
[Pywebview prerequisites](https://pywebview.flowrl.com/guide/installation)

No candidate is selected. The proposed first comparison remains Wails and a
thin host (webview or Wry), using the same Fungi fixtures and actual runtimes.

## Browser fallback is a product choice

WebUI demonstrates using an installed browser with a private profile and local
communication. Fungi could similarly serve session assets locally and return
submissions directly to the owner. This would not require the agent to inspect
the browser or scrape answers. Permission to launch applications or listen on
loopback remains subject to the caller's environment.
[WebUI design](https://github.com/webui-dev/webui#how-does-it-work)

An ordinary tab cannot preserve all native lifecycle guarantees. Browsers can
freeze/discard pages, and a lost connection does not distinguish close from
sleep, reload or suspension. An explicit submitted/cancelled result can remain
reliable; treating disconnection or an idle timer as a terminal close would
violate the AFK requirement. A dedicated browser process may offer stronger
signals but needs its own platform proof. Local communication also requires
session authentication, origin/path confinement and persistence-before-ack.
[Browser lifecycle](https://developer.chrome.com/docs/web-platform/page-lifecycle-api)

The initial requirements favor dedicated native windows. A browser fallback
or automatic OS-runtime installation is not silently accepted by this research.

## Proposed proof and unresolved decision

Probe one Wails-backed renderer behind the existing protocol before removing
Electron. Reuse the actual questionnaire and the 21 renderer checks: local
modules/workers, WebAssembly, SVG/canvas PNG, media/ranges, draft Blob recovery,
ordered persistence, limits and page isolation. Add native close/crash,
clipboard, long idle/sleep and concurrent-session checks on the real platform
webviews. Exercise the existing installed update/rollback journey with a live
window, and measure archive size and startup/memory instead of using framework
marketing numbers. A generic browser test is not proof of an embedded webview.

Keep protocol reliability invariant. Define an explicit supported OS/runtime
baseline for UI features; feature-detect optional GPU/media behavior. Existing
skill libraries need representative compatibility tests. Neither Electron nor
a system webview can guarantee every arbitrary web library on every machine.

The consequential open choice is how to handle a machine missing the runtime:
require/install the system runtime, permit an explicit browser fallback with
documented lifecycle differences, or declare that native UI environment
unsupported. The previous no-runtime-install policy remains in force until
changed. That research checkpoint preceded implementation.

## Implemented Wails assessment on 2026-10-25

The user selected Wails for a practical proof. A separate Wails v2.17 helper
now reuses Fungi's protocol and page API. The [verification record](../verification/runs/2026-10-25-wails.md)
contains the measured results and distinguishes the earlier research from
runtime observations. Electron is still the default; production migration
has not been accepted or completed.

The strongest benefit is a much smaller payload without changing skill
authoring. The proof found broader native page capabilities, shared browser
storage and Linux direct-media loading failures. The user's
[subsequent clarification](../planning.md#wails-simplicity-clarification-on-2026-10-25)
accepts native access, best-effort browser storage isolation and manual crash
recovery. Missing runtimes may be set up with explicit opt-in at installation.
Those choices supersede the earlier restrictive assumptions; no framework
fork or storage virtualization is needed.

## Linux media follow-up

The failing Ubuntu 24.04 job uses WebKitGTK 2.52.6. WebKit's
[media protocol change](https://github.com/WebKit/WebKit/commit/7e871b506eaba52c7d74519aeb0047d702101331)
allows a limited set of URL schemes by default and accepts additional schemes
through `WEBKIT_GST_ALLOWED_URI_PROTOCOLS` (comma-separated). Wails v2 uses
`wails://wails/` for local assets. The restriction is a concrete candidate
cause for direct WAV failure alongside successful Blob loading. Testing that
setting on both Ubuntu targets did not fix direct loading, so the change was
removed rather than retained as an ineffective workaround.

The matching WebKitGTK 2.52.6
[GStreamer web source](https://github.com/WebKit/WebKit/blob/webkitgtk-2.52.6/Source/WebCore/platform/graphics/gstreamer/WebKitWebSourceGStreamer.cpp#L888)
registers HTTP, HTTPS and Blob schemes, not custom schemes. An allowlist cannot
by itself supply a loader for `wails://`. A second real run with GStreamer
diagnostics confirms `No URI handler implemented for "wails"` on both Ubuntu
targets. The upstream
[custom-scheme media issue](https://bugs.webkit.org/show_bug.cgi?id=146351)
also distinguishes direct source loading from fetch-to-Blob loading. Our
[follow-up verification](../verification/runs/2026-10-25-wails-followup.md)
preserves the failed experiment and separates observed behavior from source
analysis. Small-media Blob loading is a possible authored-page workaround;
preloading an entire file is not equivalent to large-media streaming.

This sits at the Wails/WebKit integration boundary, not evidence that Linux
cannot play audio. [Wry also uses WebKitGTK on Linux](https://docs.rs/wry/latest/wry/#linux),
so changing to Wry/Tauri alone does not guarantee different media behavior.
Chromium uses a different engine but reintroduces the bundle the user wants
to remove. No renderer guarantees every codec across every installed runtime.

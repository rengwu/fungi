---
title: Rendering and packaging
updated: 2026-10-19
---
# Rendering and packaging

[Research overview](index.md)

After the first interview, the default must bundle its rendering engine. The earlier Wails/thin-webview recommendation is superseded. Those hosts remain optional extension references. The user also accepted bundled Chromium/Electron as the fallback if lighter choices are inadequate. See the [focused follow-up](bundled-renderers.md).

The second interview requires an open-source default; subsequent requests add substantial adoption and healthy recent development. Portable folder/archive packaging is accepted. Following the [maintenance assessment](renderer-health.md), the user selected Electron + Go. The broad catalog below preserves the earlier research, not a list of equally eligible defaults.

## Separate three promises

A single download, a single executable, and no non-system runtime dependencies are different properties. One Go executable can embed all HTML/JS/CSS assets while still relying on an operating system's web engine. A downloadable executable can also extract helper libraries on launch, which is one distribution file but several runtime files.

For fungi's UI-free core, a self-contained Go executable is a reasonable target. Each renderer must state its actual runtime requirements. A native window is not implemented by the Go language itself; a Go host ultimately calls platform window and rendering APIs.

| Platform | Typical embedded web engine | Consequence |
| --- | --- | --- |
| macOS | System WebKit/Cocoa | No separately shipped browser engine for this route |
| Windows | WebView2 | Check availability; Microsoft documents inclusion on Windows 11 and widespread but not universal availability on Windows 10 |
| Linux | GTK and WebKitGTK | Required libraries vary by distribution and are not an unconditional OS baseline |

The upstream webview project explicitly lists Linux distribution dependencies. Microsoft's documentation distinguishes preinstalled runtime use from distributing a runtime yourself; a Fixed Version runtime adds substantial files and size. These constraints also apply when a wrapper advertises a small executable. [webview platform requirements](https://github.com/webview/webview#prerequisites), [WebView2 distribution](https://learn.microsoft.com/en-us/microsoft-edge/webview2/concepts/distribution)

## Renderer candidates

| Candidate | Useful fit | Main tradeoff | Fungi assessment |
| --- | --- | --- | --- |
| webview + Go binding | Small host for arbitrary HTML and JS/native calls | Thin API means fungi owns more lifecycle, navigation, and platform handling | Optional extension; system-webview prerequisites exclude it as the default |
| Wails | Go desktop framework with window lifecycle, bindings, assets, and platform services | Larger framework surface; system-webview requirements remain | Optional extension; superseded as the proposed default |
| Tauri / Wry | Native desktop host with web content and a capability-oriented architecture | Rust implementation; same system-webview dependency class | Good alternative renderer, not a reason to rewrite the Go core |
| WebUI | Language-independent library with browser and optional webview modes | Browser mode relies on an installed browser and a local bridge | Useful optional fallback, not a solution to the user's preference for a dedicated window |
| Purego webview wrappers | Go cross-compilation without compiling C on the author machine | Native shared library loading/extraction still exists | Interesting packaging experiment; no-cgo does not mean no native dependencies |
| Sciter / Quark | Compact embedded HTML-like UI engine and standalone packaging | Distinct styling/runtime behavior and proprietary engine | Optional extension only under the accepted open-source-default requirement |
| Ultralight | Embeddable WebKit-derived engine with native rendering APIs | Current developer-platform and static-linking restrictions need separate agreement | Exclude from the default shortlist under ordinary SDK terms |
| Servo | Embeddable Rust web engine and WebView API | Requires a dedicated compatibility and packaging assessment | Open-source bundled-engine candidate; still pre-1.0 |
| CEF | Chromium through native C/C++ APIs | Native integration and runtime resource/helper packaging | Accepted fallback family; compare maintenance with Electron |
| Electron | Established Chromium/Node desktop architecture and web ecosystem | Bundles a broader runtime and uses multiple processes | Accepted default alongside the Go core |
| Gio / Fyne | Go-oriented native GUI approaches | Their own widget/rendering models; browser component libraries do not transfer directly | Optional native renderers for workflows that benefit |
| Slint | Declarative native UI with a Rust runtime and several language integrations | Separate UI language/toolchain and licensing choices | Alternative extension, not the default skill-authoring requirement |

The webview Go binding inspected here contains cgo directives and an older Linux WebKitGTK target than the current upstream C++ project's available choices. A Go wrapper must be evaluated on its own revision, not assumed to track upstream. [Go binding](https://github.com/webview/webview_go), [actual build directives](https://github.com/webview/webview_go/blob/master/webview.go)

Wails provides a Go/JavaScript bridge and asset embedding. Its current documentation identifies GTK/WebKitGTK and WebView2 requirements; build tool requirements should not be confused with requirements of a prebuilt release. Compare its lifecycle support with a thin binding before deciding that less framework code necessarily means a simpler fungi implementation. [Wails architecture](https://v3.wails.io/contributing/architecture/), [Wails prerequisites](https://v3.wails.io/getting-started/installation/)

Tauri documents system-webview prerequisites too. WebUI now advertises both installed-browser and optional embedded-webview modes, so it should not be dismissed using an outdated browser-only description. [Tauri prerequisites](https://v2.tauri.app/start/prerequisites/), [WebUI](https://github.com/webui-dev/webui)

The moxcomic Go wrapper explicitly describes loading native libraries and optionally extracting prebuilt copies, while still requiring GTK/WebKitGTK on Linux. This is an example of simplifying builds without eliminating runtime dependencies. Its reliability has not been tested here. [go-webview wrapper](https://github.com/moxcomic/go-webview)

Sciter's scapp supports appending UI resources into an executable. Quark supplies its own extended CSS and Reactor model; that is not evidence of drop-in compatibility with the full React/browser ecosystem. The current SDK is on GitLab; the older GitHub SDK is archived. Its engine EULA is separate from the repository's BSD license. [scapp packaging](https://docs.sciter.com/docs/scapp/), [Quark capabilities](https://quark.sciter.com/), [current engine EULA](https://gitlab.com/sciter-engine/sciter-js-sdk/-/blob/main/SCITER-ENGINE-EULA.md)

Ultralight provides official C/C++ APIs and describes a customized WebKit engine. Follow-up examination of its 2026-08 terms found an explicit developer-platform restriction in section 4.9 and a static-linking restriction in section 4.2. Fungi's intended third-party skill/runtime role makes ordinary SDK licensing an unsuitable assumption; a separate agreement would need to cover that use. [Ultralight](https://ultralig.ht/), [current terms](https://ultralig.ht/terms/2026-08)

Servo exposes a WebView API and supports multiple platforms. That establishes an embedding direction; it does not establish compatibility with fungi's eventual library suite. Electron's process model provides a useful reference for separating privileged native code from page content. [Servo](https://servo.org/), [Electron process model](https://www.electronjs.org/docs/latest/tutorial/process-model)

Gio uses Go and an immediate-mode model with platform window/input/GPU dependencies. Fyne and Slint are additional UI ecosystems; adopting any of them would add a different authoring model rather than make arbitrary web components native. [Gio](https://gioui.org/), [Fyne documentation](https://docs.fyne.io/started/), [Slint](https://slint.dev/)

## Small native tools are instructive

swiftDialog demonstrates that a command-driven native application can present Markdown, controls, and forms without becoming a full application-development environment. Its macOS scope makes it a workflow reference rather than fungi's cross-platform foundation. [swiftDialog](https://github.com/swiftDialog/swiftDialog)

ncruces/zenity is a useful example of a compact Go-facing dialog API. Its dependencies are platform-specific: Windows native support, macOS osascript, and external dialog tools on other Unixes. That is appropriate for an optional adapter, but not a universal no-install guarantee. [zenity](https://github.com/ncruces/zenity)

Terminal forms such as Huh/Bubble Tea offer another renderer, including an accessible prompt mode. A terminal owned by an agent tool is not necessarily an interactive terminal the human can use, so this cannot be an automatic GUI fallback. [Huh](https://github.com/charmbracelet/huh)

## Native capture remains its own capability

Drawing into a canvas and capturing other desktop windows have different requirements. Browser screen capture requires an explicit user gesture and renewed permission; an arbitrary webview may not implement the full API. Linux's screenshot portal has its own request/response and capability model. Reusing screencap's interaction concept is sensible, but the platform capture implementation should live behind a capability rather than be implied by loading HTML. [getDisplayMedia](https://developer.mozilla.org/en-US/docs/Web/API/MediaDevices/getDisplayMedia), [XDG Screenshot portal](https://flatpak.github.io/xdg-desktop-portal/docs/doc-org.freedesktop.portal.Screenshot.html)

## Recommended next comparison

When implementation is authorized, exercise a bundled-engine candidate against the [focused comparison](bundled-renderers.md), using the same HTML and JSON forms, slider preview, close/reopen behavior, and blob submission. Compare actual packaging on macOS, Windows, and a named Linux baseline; test keyboard input, IME, accessibility, focus, and clipboard behavior. Do not choose on unverified download-size claims alone.

The accepted fallback is a bundled Chromium/Electron renderer, distinct from opening the user's installed browser. The latter can remain an optional extension. Remote/container transport is deferred by the first interview.

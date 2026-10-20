---
title: Bundled renderer follow-up
updated: 2026-10-19
---
# Bundled renderer follow-up

[Research overview](index.md) · [Accepted decisions](../planning.md)

The user chose bundled rendering to avoid WebView2/WebKitGTK setup, then accepted Chromium/Electron as a worst-case fallback. This page records the investigation that followed. The user has since selected Electron + Go; no renderer has been run as part of fungi.

Later decisions require an open-source default, substantial adoption and recent healthy development, and portable folder/archive packaging. These supersede the proprietary-engine and extraction questions discussed below. See the [current maintenance assessment and recommendation](renderer-health.md).

## Revised assessment

| Candidate | Evidence and fit | Recommendation |
| --- | --- | --- |
| Sciter.JS / scapp | Small standalone engine; JavaScript and substantial desktop UI facilities; intentionally different DOM/CSS behavior | Excluded from the default by the later open-source requirement; retain the measurements as research |
| Servo | Open-source embeddable web engine, current library API and release process | Evaluate if open source and broader web compatibility matter enough to accept a young embedding integration |
| CEF | Bundled Chromium with native APIs, multiprocess isolation, mature browser compatibility | Strong fallback for a Go/native host; compare integration effort instead of assuming it is simpler than Electron |
| Electron | Bundled Chromium plus Node and desktop shell APIs | Strong fallback when reducing platform integration and author friction is more valuable than minimum footprint |
| Ultralight | WebKit-derived engine with official native APIs | Exclude under ordinary SDK terms: developer-platform use requires a separate agreement |
| Blitz, litehtml, RmlUi | Smaller rendering/UI building blocks with narrower web semantics | Useful specialized extensions; avoid making fungi implement the missing browser platform |
| Blitsen | Niche Blitz-based HTML/JS runtime; now embeds Bun | Interesting watchlist entry; pre-alpha, partial web APIs, changing architecture and larger current baseline |

## Sciter: inspect the current distribution, not old size claims

The archived GitHub SDKs point to the current GitLab project. Revision `18d47968a560693b5286fa54995149991a480f6f` identifies release 6.0.5.3. Downloaded its Linux x64 `scapp` and `libsciter.so` into a temporary directory and inspected ELF section/dynamic tables without executing them. [Pinned SDK](https://gitlab.com/sciter-engine/sciter-js-sdk/-/tree/18d47968a560693b5286fa54995149991a480f6f/bin/linux/x64)

| Inspected file | Bytes | Meaning |
| --- | --- | --- |
| `scapp` | 18,888,792 | Standalone runner, about 18.9 MB decimal |
| `libsciter.so` | 25,052,768 | Separate embeddable shared library, about 25.1 MB decimal |

These are alternative artifacts, not two sizes to add together for every app. Neither is a measured fungi distribution. [Inspection record](../attachments/research/sciter-linux-inspection.json)

Both list the same direct dependencies: `libuuid.so.1`, `libfontconfig.so.1`, `libfreetype.so.6`, `libGLESv2.so.2`, `libEGL.so.1`, `libstdc++.so.6`, `libm.so.6`, `libgcc_s.so.1`, `libc.so.6`, and `ld-linux-x86-64.so.2`. There is no direct GTK/WebKitGTK entry. This does not prove there are no dynamically loaded or transitive dependencies, nor establish Wayland/X11 compatibility or a clean-machine baseline. Embedded version strings include `GLIBC_2.35`; strings alone do not determine the minimum supported glibc version.

The official documentation explicitly describes browser differences, simplified DOM APIs, and its own layout alternatives. For example, ordinary inline HTML event attributes are not supported by default. A JavaScript engine is not the same thing as browser API compatibility. [Design differences](https://docs.sciter.com/docs/intro/), [HTML behavior](https://docs.sciter.com/docs/HTML/), [DOM](https://docs.sciter.com/docs/DOM/)

Scapp can carry appended UI resources. That does not establish how to combine a Go orchestrator and this renderer into one no-extraction executable. The engine has its own EULA; the SDK repository's BSD license must not be assumed to license the engine binary itself. [Packaging](https://docs.sciter.com/docs/scapp/), [engine EULA](https://gitlab.com/sciter-engine/sciter-js-sdk/-/blob/main/SCITER-ENGINE-EULA.md)

## Servo: a real embedding option, still requiring evidence

Servo announced its first crates.io library release and an LTS policy in April 2026. Its own announcement explicitly says it is not a 1.0 release. [Library release](https://servo.org/blog/2026/04/13/servo-0.1.0-release/)

GitHub release metadata for v0.7.0 lists a 49,312,528-byte Linux x64 archive, a 103,179,362-byte Windows x64 ZIP, and a 114,699,280-byte macOS arm64 DMG. These are compressed demo/application distributions, not equivalent to Sciter's individual uncompressed binaries or a minimal embedded renderer. They establish that small-source/API claims do not mean a tiny download. No Servo binary was downloaded or executed. [Release assets](https://github.com/servo/servo/releases/tag/v0.7.0)

A Go core could invoke a Rust renderer executable through the same process contract as any other native capability. Do not make every skill author learn Rust or the Servo embedding API.

## Chromium fallback: choose CEF or Electron on total complexity

CEF exposes C/C++ APIs and supports JavaScript/native messaging. It brings runtime libraries, data/resources and multiple processes; macOS has a prescribed application-bundle layout. Its documentation's single-executable process-entry example does not mean all runtime files disappear. [CEF architecture and packaging](https://chromiumembedded.github.io/cef/general_usage)

Electron provides an application shell around Chromium and Node. Its packaging also contains runtime files; an application archive is not the entire engine. A Go orchestrator can launch a bundled Electron renderer while the page sees only a narrow preload bridge. [Electron packaging](https://www.electronjs.org/docs/latest/tutorial/application-distribution), [context isolation](https://www.electronjs.org/docs/latest/tutorial/context-isolation)

Energy is an existing Go/CEF route, but its current design also includes LCL and Delphi-based bridges plus a native runtime. That may save platform code, but it is not automatically the smallest maintenance surface. [Energy architecture](https://github.com/energye/energy)

The initial proposal was automatic extraction into a private cache. The user's subsequent simplicity instruction instead selects a portable folder/archive with renderer files already present. It must work offline and need no runtime download or package-manager installation. Signing, helper layout and Linux system-library compatibility still need actual packaging tests. Do not use disabled Chromium sandboxing as a packaging shortcut.

## Why several smaller alternatives do not win automatically

Ultralight's 2026-08 license section 4.9 requires a separate agreement for developer platforms/frameworks/runtimes offered to third parties. Section 4.2 restricts static-linking approaches. This is a direct selection constraint for fungi's stated purpose, not just a question of the listed price. No license was purchased or accepted. [Versioned terms](https://ultralig.ht/terms/2026-08)

Blitz's current README calls it beta, says its HTML wrapper lacks interactivity, and says it has no JavaScript bindings; the interactive route uses Dioxus. Litehtml provides layout with host drawing callbacks and warns against using it as a full HTML engine. RmlUi expects documents authored for its own HTML/CSS subset and host-driven behavior/rendering. None is a drop-in host for arbitrary existing browser widgets. [Blitz](https://github.com/DioxusLabs/blitz), [litehtml](https://github.com/litehtml/litehtml), [RmlUi](https://github.com/mikke89/RmlUi)

Blitsen tries to bridge that gap, but its current Bun migration supersedes earlier QuickJS size claims. Its own Linux x64 reference measurement is now 144.3 MB uncompressed / 54.7 MB gzip, with a pre-alpha warning and partial web compatibility. Those are upstream measurements, not reproduced here. This is a useful example of how reducing a dependency can move complexity into a homegrown browser API implementation. [Current runtime decision](https://blitsen.dev/docs/bun-migration/)

## Small next proof

Prefer one chosen renderer in the first implementation. Supporting every surveyed engine would defeat rule-lean.

Use the same tiny HTML and JSON forms, local slider animation, canvas-to-PNG result, and one representative existing library across the candidate and a Chromium baseline. Check clean launch, keyboard/IME/focus/accessibility, artifacts, close/crash handling and offline execution. Record total shipped files/bytes, cold start and idle memory on named platforms. If a smaller engine requires substantial adapters or a second primary UI language, count that as an authoring cost rather than hiding it behind binary-size savings.

This proof remains planned. Proprietary-default acceptability and packaging are now settled; current selection depends on the [maintenance evidence](renderer-health.md) and the candidate's actual integration and compatibility results.

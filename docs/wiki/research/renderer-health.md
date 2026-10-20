---
title: Renderer adoption and maintenance
updated: 2026-10-19
---
# Renderer adoption and maintenance

[Research overview](index.md) · [Accepted decisions](../planning.md) · [Bundled renderer investigation](bundled-renderers.md)

**Decision following this research: the user accepted Electron + Go on 2026-10-19.** Electron best combines the requested substantial following, maintained releases, browser compatibility, and low integration complexity. CEF remains the principal researched alternative. No Fungi integration or performance benchmark has been run. [Accepted direction](../planning.md#accepted-electron-go-direction)

The second interview requires an open-source default. The subsequent chat instruction chooses simple portable folder/archive packaging. Together these remove the main reasons to favor a proprietary compact engine or engineer a self-extracting executable. Download size still matters, but the user has accepted Chromium/Electron if lighter options undermine reliability or ease of authoring.

## What was measured

Public GitHub API observations on **2026-10-19** cover repository stars/forks, recent default-branch commits, linked author accounts, and published releases. The commit query starts on 2026-07-21 and stops after 100 results. Separate follow-ups inspect the latest 12 commits in Energy's main and component repositories. The [saved evidence](../attachments/research/renderer-health-2026-10-19.json) includes the query method, source URLs, commit samples, release dates, and limitations.

Stars measure accumulated interest, not deployed installations. Author samples are not team sizes: they omit unlinked authors, coauthors, reviewers, other repositories, and downstream work. A capped sample can span much less than 90 days. Do not compare these counts as productivity scores. Recent pushes, default-branch commits, published releases, and security backports are different evidence.

| Project | GitHub stars | Recent evidence | Implication for Fungi |
| --- | ---: | --- | --- |
| [Electron](https://github.com/electron/electron) | 123,403 | Latest sampled commit Oct 18; 100-commit cap reached; 16 linked non-bot authors; v44.9.0 published Oct 17 UTC | Strongest overall candidate for an established desktop web host |
| [CEF](https://github.com/chromiumembedded/cef) | 4,837 | Latest sampled default-branch commit Oct 10; cap reached; 12 linked authors; stable build artifacts dated Oct 17 | Established embedding engine; native integration work remains ours |
| [Servo](https://github.com/servo/servo) | 38,186 | Latest sampled commit Oct 19; cap reached; 31 linked authors; v0.7.0 published Oct 5 | Substantial interest and active development; evolving web and embedding behavior |
| [NW.js](https://github.com/nwjs/nw.js) | 41,140 | Latest sampled commit Oct 17; 65 commits in the query; one linked author account; October release announcement | Credible Chromium alternative; more concentrated activity in the inspected repository |
| [Qt WebEngine](https://github.com/qt/qtwebengine) | 442 | Latest sampled commit Oct 16; cap reached; 13 linked authors | Active project; GitHub mirror stars underrepresent the Qt ecosystem |
| [Blitz](https://github.com/DioxusLabs/blitz) | 4,277 | Latest sampled commit Oct 18; cap reached; five linked authors | Active specialized renderer, not a complete browser JavaScript environment |
| [Energy](https://github.com/energye/energy) | 618 | No commits in the 90-day query; latest default-branch commit June 5; latest GitHub release June 17, 2025 | Its upstream CEF engine's health cannot substitute for binding maintenance |
| [cef2go](https://github.com/cztomczak/cef2go) | 912 | No commits in the query; repository push timestamp in 2020; latest GitHub release in 2014 | Unsuitable foundation for the requested actively maintained default |

All inspected repositories were unarchived; that alone says little about support. CEF, NW.js and Qt WebEngine do not use GitHub Releases as their complete release record, so an empty API response was followed up through their actual distribution channels.

## Why Electron leads

Electron's v44.9.0 release includes fixes to inherited settings for newly opened windows and concurrent certificate-verification behavior. These are concrete desktop/runtime fixes, beyond repository activity alone. Its GitHub publication timestamp is Oct 17 UTC; the release website displays Oct 16. [Release notes](https://releases.electronjs.org/release/v44.9.0), [GitHub release](https://github.com/electron/electron/releases/tag/v44.9.0)

The project documents support for its three latest stable major versions, with different backport scopes, and an eight-week major-release cadence. Its public governance separates releases, security, Chromium upgrades and other responsibilities into working groups. This gives Fungi an identifiable upstream maintenance process, although Fungi would still need to ship runtime updates. [Support timelines](https://www.electronjs.org/docs/latest/tutorial/electron-timelines), [governance](https://github.com/electron/governance)

The integration advantage is practical: **a Go core can launch a small Electron companion without adopting a third-party Go/Electron binding**. Electron owns the native window and web renderer; Go owns sessions, results and agent delivery. This accepted direction revises the original Go-owned-window idea. Ordinary HTML/CSS/JavaScript and the JSON component catalog can share the same narrow page bridge. Skill authors would not need Electron APIs or a local Node installation. Fungi's build would package the runtime. [Process model](https://www.electronjs.org/docs/latest/tutorial/process-model), [distribution](https://www.electronjs.org/docs/latest/tutorial/application-distribution), [context isolation](https://www.electronjs.org/docs/latest/tutorial/context-isolation)

The cost is a larger runtime and Chromium's multiprocess model. No footprint, startup or idle-memory comparison has been measured here. Electron popularity does not prove Fungi will run on every clean Linux desktop or that any particular component library will work correctly.

## CEF: healthy engine, separate integration decision

CEF's own documentation directs users to its binary builds. The distribution feed contains stable Linux and Windows artifacts dated Oct 17. Its lack of GitHub Releases is therefore not evidence of inactivity. OBS's browser component is a concrete downstream CEF integration. [CEF build documentation](https://chromiumembedded.github.io/cef/branches_and_building), [binary distribution](https://cef-builds.spotifycdn.com/index.html), [OBS browser component](https://github.com/obsproject/obs-browser)

CEF supplies C/C++ embedding APIs and a prescribed runtime/helper layout. Selecting it means maintaining a native host or depending on bindings, in addition to the Go session engine. It does not automatically eliminate Chromium's resources or helper processes. [General usage and application layout](https://chromiumembedded.github.io/cef/general_usage)

Energy deserves a fairer assessment than its main release date alone: its v3 work is split into [CEF](https://github.com/energye/cef) and [LCL](https://github.com/energye/lcl) repositories. Those have 16 and 35 stars respectively; their latest sampled commits are June 28 and July 1. The latest 12 commits in each of these and the main repository are attributed to the same linked account. This does not establish abandonment or the full contributor population, but it falls short of the desired evidence for a default dependency. Building our own CEF host avoids that wrapper dependency while increasing Fungi's maintenance scope.

CEF remains the alternative if an implementation experiment demonstrates a meaningful packaging or resource advantage that justifies owning more native integration. There is no measured advantage yet.

## Other established and emerging options

**Servo passes the popularity and activity tests.** Its April library announcement explicitly describes a pre-1.0 API. The October progress report documents CSS Grid becoming enabled by default, improvements to interactive text selection and accessibility, and breaking embedding API changes. These are signs of active progress and evidence that relevant behavior is still evolving. That makes it a valuable future renderer candidate, but a riskier default for arbitrary skill-authored web libraries today. [Library release](https://servo.org/blog/2026/04/13/servo-0.1.0-release/), [recent progress](https://servo.org/blog/2026/10/09/august-in-servo/), [v0.7.0](https://github.com/servo/servo/releases/tag/v0.7.0)

**NW.js also deserves consideration.** Its official blog announces v0.118 on Oct 10, following September releases. The repository author concentration is a maintenance signal, not proof of a one-person project: Chromium and other repositories supply substantial work. It remains a viable alternative, but this investigation found no requirement that makes it preferable to Electron's broader desktop integration ecosystem. [Release announcements](https://nwjs.io/blog/)

**Qt WebEngine is actively maintained despite modest mirror stars.** Qt announced a separate WebEngine release cadence in September to make browser updates less dependent on whole-Qt releases. Its Chromium baseline and backported security level must be distinguished. The tradeoff is adopting Qt libraries, plugins and helpers around a Go core. GitHub's automatically detected repository license is not the license assessment for the full module/distribution. [Release direction](https://www.qt.io/blog/qt-webengine-will-release-separate), [Chromium versions](https://wiki.qt.io/QtWebEngine/ChromiumVersions), [deployment](https://doc.qt.io/qt-6/qtwebengine-deploying.html), [module licensing](https://doc.qt.io/qt-6.8/qtwebengine-licensing.html)

**Blitz has recent development but a different scope.** Its project documentation distinguishes its HTML renderer from interactive Dioxus integration and says JavaScript bindings are absent. Making it the default would transfer missing web-platform work onto Fungi or constrain authoring. That conflicts with the requested simplicity and ordinary web-library reuse. [Project capabilities](https://github.com/DioxusLabs/blitz)

Sciter and Ultralight are excluded from the default by the accepted open-source requirement. Wails, Tauri, webview and similar system-webview hosts do not by themselves satisfy the accepted requirement to bundle a renderer without separate WebView2/WebKitGTK setup. These remain optional extension references; their popularity does not resolve that mismatch. [Earlier rendering investigation](rendering.md)

## Keep the resulting design small

Use one default renderer, a tiny page-to-session bridge, and a primitive reference component catalog. Let skills add local web assets or invoke native helper processes. Do not build a renderer marketplace, a universal UI language, multiple default engines, or an extraction/cache system to accommodate every researched option.

Renderer choice does not solve automatic agent resumption. The [delivery contract](delivery.md) still needs a supported harness route, durable results and the accepted manual “done” fallback. Human waiting must not produce repeated model tool calls regardless of renderer.

With Electron + Go selected, the useful next evidence comes from the [first end-to-end proof](../proposal.md#first-proof) and packaging checks. Broader compatibility experiments can then cover a drawing surface and representative external libraries, with measured archive size, startup and idle resources. Revisit CEF only if Electron fails a concrete requirement or an evidenced improvement justifies the extra native host work. These experiments are proposed, not performed or scheduled.

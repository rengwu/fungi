---
title: Research findings and recommendation
updated: 2026-10-19
---
# Research findings and recommendation

Research date: 2026-10-19. This is a design assessment, not a claim that fungi or any proposed integration has been implemented or verified.

## Main finding

The strongest design is a small session engine with an unusually easy authoring surface. Most interface capabilities already exist in reusable libraries. Fungi's distinctive work is maintaining the relationship between an interaction, its saved result, and the correct calling agent, including when the human takes hours to respond.

Keep ordinary HTML/CSS/JavaScript available directly, and provide the equally supported JSON authoring path accepted in the first interview. Both should use one renderer and result bridge. Native applications can use a small process contract. Avoid making authors implement delivery, recovery, or harness detection themselves.

## Research map

- [Engine and extension mechanisms](engine.md): process plugins, Go plugins, WebAssembly, storage, and failure boundaries.
- [Rendering and packaging](rendering.md): system webviews, Wails, niche engines, native toolkits, and the real meaning of one binary.
- [Bundled renderer follow-up](bundled-renderers.md): revised shortlist, actual Sciter binary inspection, licensing constraints, and the accepted Chromium/Electron fallback.
- [Renderer adoption and maintenance](renderer-health.md): dated GitHub/API evidence, releases, contributor samples, wrapper health, and the current recommendation.
- [Reconsidering the bundled renderer](system-renderers.md): October 23–25 system-webview research, followed by the user-authorized Wails prototype; production migration remains open.
- [Composition and component libraries](components.md): standard controls, forms, drawing, charts, editors, and agent-oriented UI formats.
- [Pele follow-up](pele.md): user-suggested optional diagram asset, maintenance signals and a bounded local probe.
- [Agent delivery and waiting](delivery.md): current harness mechanisms, their limits, and the no-model-polling contract.
- [Proposed small authoring interface](../proposal.md): illustrative commands and a minimal skill-authored UI.

## Recommendations after the first interview

The user chose bundled rendering instead of system webviews, both HTML and JSON authoring in the first release, and a manual “done” fallback. These decisions supersede the initial recommendations to compare Wails/thin webview hosts, defer the JSON path, and consider refusing launch without automatic delivery. [Accepted decisions](../planning.md)

The second round and chat clarifications additionally require an open-source default, substantial adoption and healthy recent development, and the simplest portable folder/archive packaging. The user subsequently accepted the resulting Electron + Go recommendation. CEF remains a researched alternative rather than a parallel implementation target. The selection is not a compatibility or performance benchmark. [Accepted direction](../planning.md#accepted-electron-go-direction)

| Concern | Recommendation | Why |
| --- | --- | --- |
| Core responsibility | Own sessions, results, artifacts, lifecycle, and delivery | These remain useful regardless of renderer or UI library |
| Common UI authoring | HTML/assets plus a tiny bridge, and a small JSON catalog | Both paths are accepted first-release scope; share implementation underneath |
| Primitive test UI | Accepted: text, form input, slider, submit/cancel and pasted/selected image | Tests real interaction without committing to a large component framework |
| Complex UI | Locally bundled libraries, selected by the skill | Reuses existing capabilities without inflating every fungi download |
| Native extensions | Ordinary processes with a small JSON contract; include one tiny helper test | Independent languages and releases; no shared-library ABI |
| Native host | Accepted: Electron companion beside the Go core | Strong adoption and maintained desktop integration reduce integration work; measure the actual footprint |
| Packaging | Portable folder/archive with one CLI entry point | Accepted simplicity choice; no installer, runtime download or extraction-cache framework |
| JSON scope | Accepted: flat built-in fields first; richer composition deferred | The user has larger JSON-authoring ideas to discuss separately |
| WASM runtime | Defer | Useful for portable computation, but adds little to the first UI workflow |
| Delivery | Best-effort supported route, save first, fall back to “done” | Open the interface even when automatic delivery is unavailable; no agent polling |
| Reliability | Exercise failure cases before expanding components | A beautiful UI with lost or misrouted answers fails the main requirement |

Electron + Go and the first milestone's scope are selected. Exact protocols and APIs remain proposals. Accepted decisions are recorded in [planning](../planning.md); the [first proof](../proposal.md#first-proof) targets macOS and verified automatic return in Codex, Claude Code, OpenCode, Grok, Kimi and Pi.

## Evidence boundary

Read the local jstack skills, collector code, native screencap helper, completion launcher, and Codex adapter. The source checkout was at `32324e4` with existing local changes to grill-form and its notes editor. No jstack files were modified. Its historical compatibility report describes tests performed by jstack; those were not repeated during this research.

The installed Codex CLI reported `0.163.0`; `codex queue --help` exposed `--thread` and `--message`. This establishes local command availability, not guaranteed future delivery or cross-harness compatibility.

During the first interview, jstack's completion command returned `armed`, a notice with the exact session path subsequently arrived in this conversation, and collection/finish recovered and saved all six answers. This is evidence of one jstack questionnaire handoff in this environment, not a fungi implementation or an AFK-duration/platform/harness compatibility test.

Reviewed primary project documentation, specifications, and selected source code. In particular, inspected the actual cgo directives in `webview_go`, the HashiCorp process handshake, and Go's file-rename contract. Library size and performance claims have not been independently benchmarked. Other operating systems and harnesses have not been exercised here. No packages, application frameworks, or agent integrations were installed.

An initial CEF usage URL and an older Servo embedding page could not be retrieved. Follow-up research found the current canonical CEF documentation and Servo's current library release documentation. Sciter's current GitLab SDK replaced archived GitHub sources for packaging/license assessment. Two Linux Sciter binaries were downloaded into a temporary research directory and inspected as data, without execution or installation; details are in the bundled-renderer follow-up. Source links throughout the research identify the evidence behind each comparison.
